package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type Section struct {
	ID            int    `json:"id"`
	SISSectionID  string `json:"sis_section_id"`
	Name          string `json:"name"`
	StartAt       string `json:"start_at"`
	EndAt         string `json:"end_at"`
	CourseID      int    `json:"course_id"`
	TotalStudents int    `json:"total_students"`
}

type SectionWithEnrollments struct {
	ID           int
	SISSectionID string
	Name         string
	Teachers     []string
}

func (c *APIClient) GetSectionsByCourseID(courseID int) ([]*Section, error) {
	sections := []*Section{}
	requestURL := fmt.Sprintf("%s/courses/%d/sections?page=1&per_page=%d&include[]=total_students", c.BaseURL, courseID, c.PageSize)

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

		_enrollments := []*Section{}
		if err := json.Unmarshal(body, &_enrollments); err != nil {
			return nil, terror.Error(err, "cannot unmarshal response body")
		}
		sections = append(sections, _enrollments...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return sections, nil
}

func (c *APIClient) GetSectionByID(sectionID int) (*Section, error) {
	section := &Section{}
	requestURL := fmt.Sprintf("%s/sections/%d", c.BaseURL, sectionID)

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

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, terror.Error(err, "cannot read response body")
	}

	if res.Status != "200 OK" {
		return nil, terror.Error(fmt.Errorf("status code: %d", res.StatusCode), "something went wrong and did not receive 200 OK status")
	}

	if err := json.Unmarshal(body, section); err != nil {
		return nil, terror.Error(err, "cannot unmarshal response body")
	}

	return section, nil
}
