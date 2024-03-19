package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ninja-software/terror/v2"
)

type Submission struct {
	ID         int    `json:"id" csv:"-"`
	Account    string `json:"-" csv:"Qualification"`
	CourseName string `json:"-" csv:"Course"`
	User       struct {
		SISUserID string `json:"sis_user_id" csv:"ID"`
		Name      string `json:"name" csv:"Name"`
	} `json:"user" csv:"User"`
	UserID          int    `json:"user_id" csv:"-"`
	AssignmentID    int    `json:"assignment_id" csv:"-"`
	AssignmentName  string `json:"assignment_name" csv:"Assignment Name"`
	AssignmentDueAt string `json:"-" csv:"Due At"`
	CourseID        int    `json:"course_id" csv:"-"`
	Grade           string `json:"grade" csv:"Grade"`
	SubmittedAt     string `json:"submitted_at" csv:"Submitted At"`
	GradedAt        string `json:"graded_at" csv:"Graded At"`
	Attempt         int    `json:"attempt" csv:"Attempt"`
	GraderID        int    `json:"grader_id" csv:"-"`
	Late            bool   `json:"late" csv:"Late"`
	Excused         bool   `json:"excused" csv:"Excused"`
	PreviewURL      string `json:"preview_url" csv:"Preview URL"`
}

func (c *APIClient) GetSubmissions(courseID int, assignmentID int) ([]Submission, error) {
	submissions := []Submission{}
	page := 1

	for {
		requestURL := fmt.Sprintf("%s/courses/%d/assignments/%d/submissions?per_page=%d&page=%d", c.BaseURL, courseID, assignmentID, c.PageSize, page)
		req, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, err
		}
		bearer := "Bearer " + c.AccessToken
		req.Header.Add("Authorization", bearer)

		res, err := c.do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var pageSubmissions []Submission
		if err := json.Unmarshal(body, &pageSubmissions); err != nil {
			return nil, err
		}

		if len(pageSubmissions) == 0 {
			break
		}

		page++
		submissions = append(submissions, pageSubmissions...)
	}

	return submissions, nil
}

func (c *APIClient) GetUngradedSubmissionsByAccount(account *Account) ([]*Submission, error) {
	submissions := []*Submission{}
	courses, err := c.GetCoursesByAccount(account, StudenCourseEnrollment)
	if err != nil {
		return nil, terror.Error(err, "error retreiving courses")
	}

	for _, course := range courses[:2] {
		assignments, err := c.GetAssignmentsByCourse(course, "ungraded")
		if err != nil {
			return nil, terror.Error(err, "error retreiving assignments")
		}
		for _, assignment := range assignments {
			requestURL := fmt.Sprintf("%s/courses/%d/assignments/%d/submissions?page=1&per_page=%d&include[]=user", c.BaseURL, course.ID, assignment.ID, c.PageSize)
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

				_submissions := []*Submission{}
				if err := json.Unmarshal(body, &_submissions); err != nil {
					return nil, terror.Error(err, "cannot unmarshall response body")
				}

				for _, submission := range _submissions {
					if submission.Grade != "" {
						continue
					}

					submission.Account = account.Name
					submission.CourseName = course.Name
					submission.AssignmentName = assignment.Name
					submission.AssignmentDueAt = assignment.DueAt

					submissions = append(submissions, submission)
				}

				nextURL := getNextURL(res.Header.Get("Link"))
				if nextURL == "" {
					break
				}

				requestURL = nextURL
			}
			fmt.Println("Completed - Course: ", course.Name, ", Assignment: ", assignment.Name)
		}
	}

	return submissions, nil
}
