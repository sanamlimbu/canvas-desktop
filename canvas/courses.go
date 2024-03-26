package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type CourseEnrollmentType string

const (
	TeacherCourseEnrollment  CourseEnrollmentType = "teacher"
	StudenCourseEnrollment   CourseEnrollmentType = "student"
	TaCourseEnrollment       CourseEnrollmentType = "ta"
	ObserverCourseEnrollment CourseEnrollmentType = "observer"
	DesignerCourseEnrollment CourseEnrollmentType = "designer"
)

// For Wails EnumBind
var AllCourseEnrollmentType = []struct {
	Value  CourseEnrollmentType
	TSName string
}{
	{TeacherCourseEnrollment, "TEACHER"},
	{StudenCourseEnrollment, "STUDENT"},
	{TaCourseEnrollment, "TA"},
	{ObserverCourseEnrollment, "OBSERVER"},
	{DesignerCourseEnrollment, "DESIGNER"},
}

type Course struct {
	ID                int        `json:"id" csv:"-"`
	AccountName       string     `json:"account_name" csv:"Account"`
	CourseCode        string     `json:"course_code" csv:"Course Code"`
	Name              string     `json:"name" csv:"Course Name"`
	GradingStandardID int        `json:"grading_standard_id" csv:"Grading Standard ID"`
	GradingStandard   string     `json:"grading_standard" csv:"Grading Standard"`
	AccountID         int        `json:"account_id" csv:"-"`
	RootAccountID     int        `json:"root_account_id" csv:"-"`
	FriendlyName      string     `json:"friendly_name" csv:"-"`
	WorkflowState     string     `json:"workflow_state" csv:"Course State"`
	StartAt           string     `json:"start_at" csv:"Start At"`
	EndAt             string     `json:"end_at" csv:"End At"`
	IsPublic          bool       `json:"is_public" csv:"-"`
	EnrollmentTermID  int        `json:"enrollment_term_id" csv:"-"`
	Account           *Account   `json:"account" csv:"-"`
	Sections          []*Section `json:"sections" csv:"-"`
}

func (c *APIClient) GetCourseByID(id int) (*Course, error) {
	course := &Course{}

	requestURL := fmt.Sprintf("%s/courses/%d?include[]=account", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, terror.Error(err, "cannot create a get request")
	}
	bearer := "Bearer " + c.AccessToken
	req.Header.Add("Authorization", bearer)

	res, err := c.do(req)
	if err != nil {
		return nil, terror.Error(err, "error on get request call")
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		return nil, terror.Error(fmt.Errorf("status code: %d", res.StatusCode), "something went wrong and did not receive 200 OK status")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, terror.Error(err, "cannot read response body")
	}

	if err := json.Unmarshal(body, course); err != nil {
		return nil, terror.Error(err, "cannot unmarshal response body")
	}
	return course, nil
}

// enrollmentType allowed values: teacher, student, ta, observer, designer
func (c *APIClient) GetCoursesByAccount(account *Account, enrollmentType CourseEnrollmentType) ([]*Course, error) {
	courses := []*Course{}
	requestURL := fmt.Sprintf("%s/accounts/%d/courses?page=1&per_page=%d&enrollment_type[]=%s&include[]=account", c.BaseURL, account.ID, c.PageSize, enrollmentType)

	for {
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, terror.Error(err, "cannot create http request")
		}
		bearer := "Bearer " + c.AccessToken
		req.Header.Add("Authorization", bearer)

		res, err := c.do(req)
		if err != nil {
			return nil, terror.Error(err, "cannot make http call")
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, terror.Error(err, "cannot read response body")
		}

		_courses := []*Course{}
		if err := json.Unmarshal(body, &_courses); err != nil {
			return nil, terror.Error(err, "cannot unmarshall response body")
		}
		courses = append(courses, _courses...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return courses, nil
}

func (c *APIClient) GetCoursesByAccountID(accountID int) ([]*Course, error) {
	courses := []*Course{}
	requestURL := fmt.Sprintf("%s/accounts/%d/courses?page=1&per_page=%d&include[]=account", c.BaseURL, accountID, c.PageSize)

	for {
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, terror.Error(err, "cannot create http request")
		}
		bearer := "Bearer " + c.AccessToken
		req.Header.Add("Authorization", bearer)

		res, err := c.do(req)
		if err != nil {
			return nil, terror.Error(err, "cannot make http call")
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, terror.Error(err, "cannot read response body")
		}

		_courses := []*Course{}
		if err := json.Unmarshal(body, &_courses); err != nil {
			return nil, terror.Error(err, "cannot unmarshall response body")
		}
		courses = append(courses, _courses...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return courses, nil
}

func (c *APIClient) GetCoursesByUser(user *User) ([]*Course, error) {
	courses := []*Course{}
	requestURL := fmt.Sprintf("%s/users/%d/courses?page=1&per_page=%d&include[]=account&include[]=sections", c.BaseURL, user.ID, c.PageSize)

	for {
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, terror.Error(err, "cannot create http request")
		}
		bearer := "Bearer " + c.AccessToken
		req.Header.Add("Authorization", bearer)

		res, err := c.do(req)
		if err != nil {
			return nil, terror.Error(err, "cannot make http call")
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, terror.Error(err, "cannot read response body")
		}

		_courses := []*Course{}
		if err := json.Unmarshal(body, &_courses); err != nil {
			return nil, terror.Error(err, "cannot unmarshall response body")
		}
		courses = append(courses, _courses...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return courses, nil
}
