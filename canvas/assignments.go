package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ninja-software/terror/v2"
)

type AssignmentBucket string

const (
	PastBucket        AssignmentBucket = "past"
	OverdueBucket     AssignmentBucket = "overdue"
	UndatedBucket     AssignmentBucket = "undated"
	UngradedBucket    AssignmentBucket = "ungraded"
	UnsubmittedBucket AssignmentBucket = "unsubmitted"
	UpcomingBucket    AssignmentBucket = "upcoming"
	FutureBucket      AssignmentBucket = "future"
)

// For Wails EnumBind
var AllAssignmentBucket = []struct {
	Value  AssignmentBucket
	TSName string
}{
	{PastBucket, "PAST"},
	{OverdueBucket, "OVERDUE"},
	{UndatedBucket, "UPDATED"},
	{UngradedBucket, "UNGRADED"},
	{UnsubmittedBucket, "UNSUBMITTED"},
	{UpcomingBucket, "UPCOMING"},
	{FutureBucket, "FUTURE"},
}

type Assignment struct {
	ID                         int                    `json:"id" csv:"-"`
	CourseID                   int                    `json:"course_id" csv:"-"`
	Account                    string                 `json:"account" csv:"Account"`
	CourseName                 string                 `json:"course_name" csv:"Course Name"`
	Name                       string                 `json:"name" csv:"Assignment"`
	DueAt                      string                 `json:"due_at" csv:"Due"`
	UnlockAt                   string                 `json:"unlock_at" csv:"Available From"`
	LockAt                     string                 `json:"lock_at" csv:"Until"`
	NeedsGradingCount          int                    `json:"needs_grading_count" csv:"-"`
	Section                    string                 `json:"section" csv:"Section"`
	NeedingGradingSection      int                    `json:"needs_grading_section" csv:"Needs Grading"`
	Teachers                   string                 `json:"teachers" csv:"Teachers"`
	Bucket                     string                 `json:"bucket" csv:"Bucket"`
	Published                  bool                   `json:"published" csv:"Published"`
	GradebookURL               string                 `json:"gradebook_url" csv:"Gradebook URL"`
	NeedsGradingCountBySection []*SectionNeedsGrading `json:"needs_grading_count_by_section" csv:"-"`
	AllDates                   []*AssignmentDate      `json:"all_dates" csv:"-"`
	GradingStandardID          int                    `json:"grading_standard_id" csv:"-"`
	GradingType                string                 `json:"grading_type" csv:"-"`
	OmitFromFinalGrade         bool                   `json:"omit_from_final_grade" csv:"-"`
	WorkflowState              string                 `json:"workflow_state" csv:"-"`
}

type SectionNeedsGrading struct {
	SectionID         int `json:"section_id" csv:"-"`
	NeedsGradingCount int `json:"needs_grading_count" csv:"-"`
}

type AssignmentResult struct {
	AssignmentID   int     `json:"assignment_id" csv:"-"`
	UserSisID      string  `json:"user_sis_id" csv:"SIS ID"`
	Name           string  `json:"name" csv:"Name"`
	Acccount       string  `json:"account" csv:"Account"`
	CourseName     string  `json:"course_name" csv:"Course Name"`
	Title          string  `json:"title" csv:"Assignment"`
	MaxScore       float32 `json:"max_score" csv:"Max Score"`
	MinScore       float32 `json:"min_score" csv:"Min Score"`
	PointsPossible float32 `json:"points_possible" csv:"Points Possible"`
	Submission     struct {
		Score       float32 `json:"score" csv:"Score"`
		SubmittedAt string  `json:"submitted_at" csv:"Submitted At"`
	} `json:"submission"`
	Status      string `json:"status" csv:"Submission Status"`
	DueAt       string `json:"due_at" csv:"Due At"`
	CourseState string `json:"course_state" csv:"Course State"`
}

type AssignmentDate struct {
	ID       int    `json:"id"`
	DueAt    string `json:"due_at"`
	UnlockAt string `json:"unlock_at"`
	LockAt   string `json:"lock_at"`
	Title    string `json:"title"`
	SetType  string `json:"set_type"`
	SetID    int    `json:"set_id"`
}

type AssignmentGradingStandard struct {
	ID                 int    `json:"id" csv:"-"`
	CourseID           int    `json:"course_id" csv:"-"`
	Account            string `json:"Account" csv:"Account"`
	CourseName         string `json:"course_name" csv:"Course Name"`
	Name               string `json:"name" csv:"Assignment"`
	CourseState        string `json:"course_state" csv:"Course State"`
	GradingStandardID  int    `json:"grading_standard_id" csv:"Grading Standard ID"`
	GradingStandard    string `json:"grading_standard" csv:"Grading Standard"`
	GradingType        string `json:"grading_type" csv:"Grading Type"`
	OmitFromFinalGrade bool   `json:"omit_from_final_grade" csv:"Omit Final Grade"`
	WorkflowState      string `json:"workflow_state" csv:"Workflow State"`
	DueAt              string `json:"due_at" csv:"Due"`
	UnlockAt           string `json:"unlock_at" csv:"Available From"`
	LockAt             string `json:"lock_at" csv:"Until"`
}

func (c *APIClient) GetAssignmentsResultsByUser(user *User) ([]*AssignmentResult, error) {
	results := []*AssignmentResult{}
	courses := make(map[int]*Course)

	_courses, err := c.GetCoursesByUser(user)
	if err != nil {
		return nil, terror.Error(err, fmt.Sprintf("cannot get courses of user SIS ID:%s", user.SISUserID))
	}

	for _, course := range _courses {
		courses[course.ID] = course
	}

	enrollments, err := c.GetEnrollmentsByUser(user)
	if err != nil {
		return nil, terror.Error(err, fmt.Sprintf("cannot get enrollments of user SIS ID: %s", user.SISUserID))
	}
	for _, enrollment := range enrollments {
		_results, err := c.GetAssignmentsResultByEnrollment(user, enrollment)
		if err != nil {
			return nil, terror.Error(err, fmt.Sprintf("cannot get assignments result of user SIS ID: %s", user.SISUserID))
		}

		for _, result := range _results {
			result.UserSisID = user.SISUserID
			result.Name = user.Name
			if courses[enrollment.CourseID] != nil {
				result.Acccount = courses[enrollment.CourseID].Account.Name
				result.CourseName = courses[enrollment.CourseID].Name
				result.CourseState = courses[enrollment.CourseID].WorkflowState

				sections := ""
				total := len(courses[enrollment.CourseID].Sections)
				for i, section := range courses[enrollment.CourseID].Sections {
					if i == total-1 {
						sections = sections + section.Name
					} else {
						sections = sections + section.Name + ";"
					}
				}
			}
			results = append(results, result)
		}
	}

	return results, nil
}

// bucket allowed values: past, overdue, undated, ungraded, unsubmitted, upcoming, future
func (c *APIClient) GetAssignmentsByCourseAndBucket(course *Course, bucket AssignmentBucket) ([]*Assignment, error) {
	assignments := []*Assignment{}
	requestURL := fmt.Sprintf("%s/courses/%d/assignments?page=1&per_page=%d&bucket=%s&needs_grading_count_by_section=true&include[]=all_dates", c.BaseURL, course.ID, c.PageSize, bucket)
	sections := make(map[int]*SectionWithEnrollments)
	trimmedBaseURL := strings.TrimSuffix(c.BaseURL, "/api/v1")

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

		_assignments := []*Assignment{}
		if err := json.Unmarshal(body, &_assignments); err != nil {
			return nil, terror.Error(err, "cannot unmarshall response body")
		}

		for _, _assignment := range _assignments {
			dates := make(map[int]*AssignmentDate)
			for _, date := range _assignment.AllDates {
				if date.SetID != 0 {
					dates[date.SetID] = date
				}
			}

			for _, section := range _assignment.NeedsGradingCountBySection {
				if sections[section.SectionID] == nil {
					enrollments, err := c.GetEnrollmentsBySectionID(section.SectionID, TeacherEnrollment)
					if err != nil {
						return nil, terror.Error(err, "error retreiving enrollments")
					}
					// No teachers in the section
					if len(enrollments) == 0 {
						_section, err := c.GetSectionByID(section.SectionID)
						if err != nil {
							return nil, terror.Error(err, "error retreiving section")
						}

						if _section.SISSectionID == "" {
							sections[section.SectionID] = &SectionWithEnrollments{
								ID:           _section.ID,
								SISSectionID: _section.Name,
								Teachers:     []string{},
							}

						} else {
							sections[section.SectionID] = &SectionWithEnrollments{
								ID:           _section.ID,
								SISSectionID: _section.SISSectionID,
								Teachers:     []string{},
							}
						}

						continue
					}

					teachers := []string{}
					for _, enrollment := range enrollments {
						teachers = append(teachers, enrollment.User.Name)
					}

					if enrollments[0].SISSectionID == "" {
						_section, err := c.GetSectionByID(section.SectionID)
						if err != nil {
							return nil, terror.Error(err, "error retreiving section")
						}

						sections[section.SectionID] = &SectionWithEnrollments{
							ID:           section.SectionID,
							SISSectionID: _section.Name,
							Teachers:     teachers,
						}
					} else {
						sections[section.SectionID] = &SectionWithEnrollments{
							ID:           section.SectionID,
							SISSectionID: enrollments[0].SISSectionID,
							Teachers:     teachers,
						}
					}
				}

				assignment := &Assignment{
					ID:                         _assignment.ID,
					CourseID:                   _assignment.CourseID,
					Name:                       _assignment.Name,
					NeedsGradingCount:          _assignment.NeedsGradingCount,
					Section:                    sections[section.SectionID].SISSectionID,
					NeedingGradingSection:      section.NeedsGradingCount,
					Teachers:                   strings.Join(sections[section.SectionID].Teachers, ";"),
					Published:                  _assignment.Published,
					NeedsGradingCountBySection: _assignment.NeedsGradingCountBySection,
					Account:                    course.Account.Name,
					CourseName:                 course.Name,
					Bucket:                     string(bucket),
					GradebookURL:               fmt.Sprintf(`%s/courses/%d/gradebook`, trimmedBaseURL, course.ID),
				}

				// Date set type is ADHOC
				if dates[section.SectionID] == nil {
					assignment.DueAt = ""
					assignment.LockAt = ""
					assignment.UnlockAt = ""
				} else {
					assignment.DueAt = dates[section.SectionID].DueAt
					assignment.LockAt = dates[section.SectionID].LockAt
					assignment.UnlockAt = dates[section.SectionID].UnlockAt
				}

				assignments = append(assignments, assignment)
			}
		}

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return assignments, nil
}

func (c *APIClient) GetAssignmentsByCourse(course *Course) ([]*Assignment, error) {
	assignments := []*Assignment{}
	requestURL := fmt.Sprintf("%s/courses/%d/assignments?page=1&per_page=%d", c.BaseURL, course.ID, c.PageSize)

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

		_assignments := []*Assignment{}
		if err := json.Unmarshal(body, &_assignments); err != nil {
			return nil, terror.Error(err, "cannot unmarshall response body")
		}

		assignments = append(assignments, _assignments...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return assignments, nil
}

// bucket allowed values: past, overdue, undated, ungraded, unsubmitted, upcoming, future
func (c *APIClient) GetAssignmentsByAccount(account *Account, bucket AssignmentBucket) ([]*Assignment, error) {
	assignments := []*Assignment{}

	courses, err := c.GetCoursesByAccount(account, StudenCourseEnrollment)
	if err != nil {
		return nil, terror.Error(err, "error retrieving courses")
	}

	for _, course := range courses {
		_assignments, err := c.GetAssignmentsByCourseAndBucket(course, UngradedBucket)
		if err != nil {
			return nil, terror.Error(err, "error retrieving ungraded assignments")
		}

		assignments = append(assignments, _assignments...)
	}

	return assignments, nil
}

func (c *APIClient) GetAssignmentsResultByEnrollment(user *User, enrollment *Enrollment) ([]*AssignmentResult, error) {
	results := []*AssignmentResult{}
	requestURL := fmt.Sprintf("%s/courses/%d/analytics/users/%d/assignments?page=1&per_page=%d", c.BaseURL, enrollment.CourseID, user.ID, c.PageSize)

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

		// Unpublished course
		if res.Status != "200 OK" {
			fmt.Println("unpublished course: ", enrollment.CourseID)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, terror.Error(err, "cannot read response body")
		}

		_results := []*AssignmentResult{}
		if err := json.Unmarshal(body, &_results); err != nil {
			return nil, terror.Error(err, "cannot unmarshal response body")
		}
		results = append(results, _results...)

		nextURL := getNextURL(res.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		requestURL = nextURL
	}

	return results, nil
}

func (c *APIClient) GetAssignmentGradingStandardsByCourse(course *Course) ([]*AssignmentGradingStandard, error) {
	assignments := []*AssignmentGradingStandard{}

	_assignments, err := c.GetAssignmentsByCourse(course)
	if err != nil {
		return nil, terror.Error(err, "error retrieving assignments")
	}

	for _, _assignment := range _assignments {
		assignment := &AssignmentGradingStandard{
			ID:                 _assignment.ID,
			CourseID:           _assignment.CourseID,
			Name:               _assignment.Name,
			GradingStandardID:  _assignment.GradingStandardID,
			GradingType:        _assignment.GradingType,
			OmitFromFinalGrade: _assignment.OmitFromFinalGrade,
			DueAt:              _assignment.DueAt,
			WorkflowState:      _assignment.WorkflowState,
			UnlockAt:           _assignment.UnlockAt,
			LockAt:             _assignment.LockAt,
			Account:            course.Account.Name,
			CourseName:         course.Name,
			GradingStandard:    course.GradingStandard,
			CourseState:        course.WorkflowState,
		}
		assignments = append(assignments, assignment)
	}

	return assignments, nil
}
