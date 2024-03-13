package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type Course struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	CourseCode       string `json:"course_code" `
	AccountID        int    `json:"account_id"`
	RootAccountID    int    `json:"root_account_id"`
	FriendlyName     string `json:"friendly_name"`
	WorkflowState    string `json:"workflow_state"`
	StartAt          string `json:"start_at"`
	EndAt            string `json:"end_at"`
	IsPublic         bool   `json:"is_public"`
	EnrollmentTermID int    `json:"enrollment_term_id"`
	Account          struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		WorkflowState string `json:"workflow_state"`
	} `json:"account"`
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
func (c *APIClient) GetCoursesByAccount(account *Account, enrollmentType string) ([]*Course, error) {
	courses := []*Course{}
	requestURL := fmt.Sprintf("%s/accounts/%d/courses?page=1&per_page=%d&enrollment_type[]=%s", c.BaseURL, account.ID, c.PageSize, enrollmentType)

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
