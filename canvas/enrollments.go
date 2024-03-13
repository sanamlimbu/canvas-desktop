package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

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
}

type EnrollmentResult struct {
	StudentID     string  `csv:"Student ID"`
	StudentName   string  `csv:"Student Name"`
	Qualification string  `csv:"Qualification"`
	CourseName    string  `csv:"Course Name"`
	CourseStatus  string  `csv:"Course Status"`
	CurrentGrade  string  `csv:"Current Grade"`
	CurrentScore  float32 `csv:"Current Score"`
	GradesURL     string  `csv:"Grades URL"`
}

func (c *APIClient) GetEnrollmentsByUserID(userID int) ([]*Enrollment, error) {
	enrollments := []*Enrollment{}
	requestURL := fmt.Sprintf("%s/users/%d/enrollments?page=1&per_page=%d", c.BaseURL, userID, c.PageSize)

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

func (c *APIClient) GetAllEnrollmentsResultsByUserID(userID int) ([]*EnrollmentResult, error) {
	results := []*EnrollmentResult{}
	enrollments, err := c.GetEnrollmentsByUserID(userID)
	if err != nil {
		return nil, terror.Error(err, fmt.Sprintf("cannot get enrollments of user ID:%d", userID))
	}

	for _, enrollment := range enrollments {
		course, err := c.GetCourseByID(enrollment.CourseID)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		result := &EnrollmentResult{
			StudentID:     enrollment.User.SISUserID,
			StudentName:   enrollment.User.Name,
			Qualification: course.Account.Name,
			CourseName:    course.Name,
			CourseStatus:  course.WorkflowState,
			CurrentGrade:  enrollment.Grades.CurrentGrade,
			CurrentScore:  enrollment.Grades.CurrentScore,
			GradesURL:     enrollment.Grades.HtmlUrl,
		}

		results = append(results, result)
	}

	return results, nil
}

// enrollmentType accepted values: All, StudentEnrollment, TeacherEnrollment, TaEnrollment, DesignerEnrollment, and ObserverEnrollment
func (c *APIClient) GetEnrollmentsBySectionID(sectionID int, enrollmentTypes ...string) ([]*Enrollment, error) {
	enrollments := []*Enrollment{}
	requestURL := fmt.Sprintf("%s/sections/%d/enrollments?page=1&per_page=%d", c.BaseURL, sectionID, c.PageSize)
	if !slices.Contains(enrollmentTypes, "All") {
		for _, enrollmentType := range enrollmentTypes {
			requestURL += fmt.Sprintf(`&type[]=%s`, enrollmentType)
		}
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
