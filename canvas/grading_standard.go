package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type GradingStandard struct {
	ID            int              `json:"ID"`
	Title         string           `json:"title"`
	ContextType   string           `json:"context_type"`
	ContextID     int              `json:"context_id"`
	GradingScheme []*GradingScheme `json:"grading_scheme"`
}

type GradingScheme struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

func (c *APIClient) GetGradingStandardsByAccountContext(accountID int) ([]*GradingStandard, error) {
	gradingStandards := []*GradingStandard{}
	requestURL := fmt.Sprintf("%s/accounts/%d/grading_standards?page=1&per_page=%d", c.BaseURL, accountID, c.PageSize)

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

		_gradingStandards := []*GradingStandard{}
		if err := json.Unmarshal(body, &_gradingStandards); err != nil {
			return nil, terror.Error(err, "cannot unmarshal response body")
		}
		gradingStandards = append(gradingStandards, _gradingStandards...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return gradingStandards, nil
}

func (c *APIClient) GetGradingStandardsByCourseContext(courseID int) ([]*GradingStandard, error) {
	gradingStandards := []*GradingStandard{}
	requestURL := fmt.Sprintf("%s/courses/%d/grading_standards?page=1&per_page=%d", c.BaseURL, courseID, c.PageSize)

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

		_gradingStandards := []*GradingStandard{}
		if err := json.Unmarshal(body, &_gradingStandards); err != nil {
			return nil, terror.Error(err, "cannot unmarshal response body")
		}
		gradingStandards = append(gradingStandards, _gradingStandards...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return gradingStandards, nil
}
