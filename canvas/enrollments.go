package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type Enrollment struct {
	ID              int    `json:"id"`
	UserID          int    `json:"user_id"`
	CourseID        int    `json:"course_id"`
	CourseSectionID int    `json:"course_section_id"`
	SISSectionID    string `json:"sis_section_id"`
	Grades          struct {
		HtmlUrl      string  `json:"html_url"`
		CurrentScore float32 `json:"current_score"`
		CurrentGrade string  `json:"current_grade"`
		FinalScore   float32 `json:"final_score"`
		FinalGrade   string  `json:"final_grade"`
	} `json:"grades"`
	User struct {
		Name      string `json:"name"`
		SISUserID string `json:"sis_user_id"`
	} `json:"user"`
	EnrollmentState string `json:"enrollment_state"`
	Role            string `json:"role"`
}

type EnrollmentResult struct {
	SISID           string  `json:"sis_id" csv:"SIS ID"`
	Name            string  `json:"name" csv:"Name"`
	Account         string  `json:"account" csv:"Account"`
	CourseName      string  `json:"course_name" csv:"Course Name"`
	Section         string  `json:"section" csv:"Section"`
	EnrollmentState string  `json:"enrollment_state" csv:"Enrollment State"`
	CourseState     string  `json:"course_state" csv:"Course State"`
	CurrentGrade    string  `json:"current_grade" csv:"Current Grade"`
	CurrentScore    float32 `json:"current_score" csv:"Current Score"`
	EnrollmentRole  string  `json:"enrollment_role" csv:"Enrollment Role"`
	GradesURL       string  `json:"gardes_url" csv:"Grades URL"`
}

type EnrollmentType string

const (
	TeacherEnrollment  EnrollmentType = "TeacherEnrollment"
	StudentEnrollment  EnrollmentType = "StudentEnrollment"
	TaEnrollment       EnrollmentType = "TaEnrollment"
	DesignerEnrollment EnrollmentType = "DesignerEnrollment"
	ObserverEnrollment EnrollmentType = "ObserverEnrollment"
)

// For Wails EnumBind
var AllEnrollmentType = []struct {
	Value  EnrollmentType
	TSName string
}{
	{TeacherEnrollment, "TEACHER"},
	{StudentEnrollment, "STUDENT"},
	{TaEnrollment, "TA"},
	{DesignerEnrollment, "DESIGNER"},
	{ObserverEnrollment, "OBSERVER"},
}

// "deleted" removed
var enrollmentStates []string = []string{"active", "invited", "creation_pending", "rejected", "completed", "inactive", "current_and_invited", "current_and_future", "current_and_concluded"}

func (c *APIClient) GetEnrollmentsByUser(user *User) ([]*Enrollment, error) {
	enrollments := []*Enrollment{}
	requestURL := fmt.Sprintf("%s/users/%d/enrollments?page=1&per_page=%d", c.BaseURL, user.ID, c.PageSize)
	for _, enrollmentState := range enrollmentStates {
		requestURL += fmt.Sprintf(`&state[]=%s`, enrollmentState)
	}

	for {
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

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, terror.Error(err, "cannot read response body")
		}

		if res.Status != "200 OK" {
			return nil, terror.Error(fmt.Errorf("status code: %d", res.StatusCode), "something went wrong and did not receive 200 OK status")
		}

		_enrollments := []*Enrollment{}
		if err := json.Unmarshal(body, &_enrollments); err != nil {
			return nil, terror.Error(err, "cannot unmarshal response body")
		}
		enrollments = append(enrollments, _enrollments...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return enrollments, nil
}

func (c *APIClient) GetEnrollmentResultsByUser(user *User) ([]*EnrollmentResult, error) {
	results := []*EnrollmentResult{}
	courses := make(map[int]*Course)

	enrollments, err := c.GetEnrollmentsByUser(user)
	if err != nil {
		return nil, terror.Error(err, fmt.Sprintf("cannot get enrollments of user SIS ID:%s", user.SISUserID))
	}

	_courses, err := c.GetCoursesByUser(user)
	if err != nil {
		return nil, terror.Error(err, fmt.Sprintf("cannot get courses of user SIS ID:%s", user.SISUserID))
	}

	for _, course := range _courses {
		courses[course.ID] = course
	}

	for _, enrollment := range enrollments {
		result := &EnrollmentResult{
			SISID:           enrollment.User.SISUserID,
			Name:            enrollment.User.Name,
			Section:         enrollment.SISSectionID,
			CurrentGrade:    enrollment.Grades.CurrentGrade,
			CurrentScore:    enrollment.Grades.CurrentScore,
			GradesURL:       enrollment.Grades.HtmlUrl,
			EnrollmentState: enrollment.EnrollmentState,
			EnrollmentRole:  enrollment.Role,
		}

		if courses[enrollment.CourseID] == nil {
			result.CourseName = ""
			result.CourseState = ""
			result.Account = ""
		} else {
			result.CourseName = courses[enrollment.CourseID].Name
			result.CourseState = courses[enrollment.CourseID].WorkflowState
			result.Account = courses[enrollment.CourseID].Account.Name
		}

		results = append(results, result)
	}

	return results, nil
}

// enrollmentType accepted values: StudentEnrollment, TeacherEnrollment, TaEnrollment, DesignerEnrollment, and ObserverEnrollment
func (c *APIClient) GetEnrollmentsBySectionID(sectionID int, enrollmentTypes ...EnrollmentType) ([]*Enrollment, error) {
	enrollments := []*Enrollment{}
	requestURL := fmt.Sprintf("%s/sections/%d/enrollments?page=1&per_page=%d", c.BaseURL, sectionID, c.PageSize)
	for _, enrollmentType := range enrollmentTypes {
		requestURL += fmt.Sprintf(`&type[]=%s`, enrollmentType)
	}

	for {
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

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, terror.Error(err, "cannot read response body")
		}

		if res.Status != "200 OK" {
			return nil, terror.Error(fmt.Errorf("status code: %d", res.StatusCode), "something went wrong and did not receive 200 OK status")
		}

		_enrollments := []*Enrollment{}
		if err := json.Unmarshal(body, &_enrollments); err != nil {
			return nil, terror.Error(err, "cannot unmarshal response body")
		}
		enrollments = append(enrollments, _enrollments...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return enrollments, nil
}
