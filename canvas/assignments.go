package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ninja-software/terror/v2"
)

type AssignmentSubmissionStatus string

const (
	Past        AssignmentSubmissionStatus = "past"
	Overdue     AssignmentSubmissionStatus = "overdue"
	Undated     AssignmentSubmissionStatus = "undated"
	Ungraded    AssignmentSubmissionStatus = "ungraded"
	Unsubmitted AssignmentSubmissionStatus = "unsubmitted"
	Upcoming    AssignmentSubmissionStatus = "Upcoming"
	Future      AssignmentSubmissionStatus = "Future"
)

type Assignment struct {
	ID                         int    `json:"id" csv:"-"`
	CourseID                   int    `json:"course_id" csv:"-"`
	Account                    string `json:"-" csv:"Qualification"`
	CourseName                 string `json:"-" csv:"Course Name"`
	Name                       string `json:"name" csv:"Assignment"`
	DueAt                      string `json:"-" csv:"Due"`
	UnlockAt                   string `json:"-" csv:"Available From"`
	LockAt                     string `json:"-" csv:"Until"`
	NeedsGradingCount          int    `json:"needs_grading_count" csv:"Total Needs Grading"`
	Section                    string `json:"-" csv:"Section"`
	NeedingGradingSection      int    `json:"-" csv:"Needs Grading (Section)"`
	Teachers                   string `json:"-" csv:"Teachers"`
	Status                     string `json:"-" csv:"Status"`
	Published                  bool   `json:"published" csv:"Published"`
	GradebookURL               string `json:"-" csv:"Gradebook URL"`
	NeedsGradingCountBySection []*struct {
		SectionID         int `json:"section_id" csv:"-"`
		NeedsGradingCount int `json:"needs_grading_count" csv:"-"`
	} `json:"needs_grading_count_by_section" csv:"-"`
	AllDates []*AssignmentDate `json:"all_dates" csv:"-"`
}

type AssignmentResult struct {
	AssignmentID  int     `json:"assignment_id" csv:"-"`
	UserSisID     string  `csv:"Student ID"`
	StudentName   string  `csv:"Student Name"`
	Qualification string  `csv:"Qualification"`
	CourseName    string  `csv:"Course Name"`
	Title         string  `json:"title" csv:"Assignment"`
	MaxScore      float32 `json:"max_score" csv:"Max Score"`
	MinScore      float32 `json:"min_score" csv:"Min Score"`
	Submission    struct {
		Score       float32 `json:"score" csv:"Score"`
		SubmittedAt string  `json:"submitted_at" csv:"Submitted At"`
	} `json:"submission"`
	DueAt  string `json:"due-at" csv:"Due At"`
	Status string `json:"status" csv:"Submission Status"`
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

func (c *APIClient) GetAssignmentsResultsByUser(user *User) ([]*AssignmentResult, error) {
	results := []*AssignmentResult{}
	enrollments, err := c.GetEnrollmentsByUserID(user.ID)
	if err != nil {
		return nil, terror.Error(err, fmt.Sprintf("cannot get enrollments of user ID: %s", user.SISUserID))
	}

	for _, enrollment := range enrollments {
		ars := []*AssignmentResult{}

		requestURL := fmt.Sprintf("%s/courses/%d/analytics/users/%d/assignments?per_page=%d", c.BaseURL, enrollment.CourseID, user.ID, c.PageSize)
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
			res.Body.Close()
			continue
		}

		course, err := c.GetCourseByID(enrollment.CourseID)
		if err != nil {
			terror.Error(err, "error fetching course")
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, terror.Error(err, "cannot read response body")
		}

		if err := json.Unmarshal(body, &ars); err != nil {
			return nil, terror.Error(err, "cannot unmarshal response body")
		}

		for _, result := range ars {
			result.CourseName = course.Name
			result.UserSisID = user.SISUserID
			result.StudentName = user.Name
		}

		results = append(results, ars...)
	}

	return results, nil
}

// bucket allowed values: past, overdue, undated, ungraded, unsubmitted, upcoming, future
func (c *APIClient) GetAssignmentsByCourseID(courseID int, bucket string) ([]*Assignment, error) {
	assignments := []*Assignment{}
	requestURL := fmt.Sprintf("%s/courses/%d/assignments?page=1&per_page=%d&bucket=%s&needs_grading_count_by_section=true&include[]=all_dates", c.BaseURL, courseID, c.PageSize, bucket)
	sections := make(map[int]*SectionWithEnrollments)

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
					enrollments, err := c.GetEnrollmentsBySectionID(section.SectionID, "TeacherEnrollment")
					if err != nil {
						return nil, terror.Error(err, "error retreiving enrollments")
					}
					// No teachers in the section
					if len(enrollments) == 0 {
						_section, err := c.GetSectionsByID(section.SectionID)
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
						_section, err := c.GetSectionsByID(section.SectionID)
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
					Account:                    _assignment.Account,
					CourseName:                 _assignment.CourseName,
					Name:                       _assignment.Name,
					NeedsGradingCount:          _assignment.NeedsGradingCount,
					Section:                    sections[section.SectionID].SISSectionID,
					NeedingGradingSection:      section.NeedsGradingCount,
					Teachers:                   strings.Join(sections[section.SectionID].Teachers, ";"),
					Status:                     _assignment.Status,
					Published:                  _assignment.Published,
					GradebookURL:               _assignment.GradebookURL,
					NeedsGradingCountBySection: _assignment.NeedsGradingCountBySection,
				}
				// Date set type is ADHOC
				if dates[section.SectionID] == nil {
					assignment.DueAt = ""
					assignment.LockAt = ""
					assignment.UnlockAt = ""
				} else {
					dueAt, err := UTCToPerthTime(dates[section.SectionID].DueAt)
					if err != nil {
						return nil, terror.Error(err, "error converting UTC to Perth time")
					}

					lockAt, err := UTCToPerthTime(dates[section.SectionID].LockAt)
					if err != nil {
						return nil, terror.Error(err, "error converting UTC to Perth time")
					}

					unlockAt, err := UTCToPerthTime(dates[section.SectionID].UnlockAt)
					if err != nil {
						return nil, terror.Error(err, "error converting UTC to Perth time")
					}
					assignment.DueAt = dueAt
					assignment.LockAt = lockAt
					assignment.UnlockAt = unlockAt
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

// bucket allowed values: past, overdue, undated, ungraded, unsubmitted, upcoming, future
func (c *APIClient) GetAssignmentsByAccount(account *Account, bucket string) ([]*Assignment, error) {
	assignments := []*Assignment{}
	trimmedBaseURL := strings.TrimSuffix(c.BaseURL, "/api/v1")

	courses, err := c.GetCoursesByAccount(account, "student")
	if err != nil {
		return nil, terror.Error(err, "error retrieving courses")
	}

	for _, course := range courses {
		_assignments, err := c.GetAssignmentsByCourseID(course.ID, "ungraded")
		if err != nil {
			return nil, terror.Error(err, "error retrieving ungraded assignments")
		}

		for _, assignment := range _assignments {
			assignment.Account = account.Name
			assignment.CourseName = course.Name
			assignment.Status = bucket
			assignment.GradebookURL = fmt.Sprintf(`%s/courses/%d/gradebook`, trimmedBaseURL, course.ID)
		}

		assignments = append(assignments, _assignments...)
		fmt.Println("Completed: ", account.Name, " ", course.Name)
	}

	return assignments, nil
}

// // bucket allowed values: past, overdue, undated, ungraded, unsubmitted, upcoming, future
// func (c *APIClient) GetAssignmentsByAccount(account *Account, bucket string) ([]*Assignment, error) {
// 	assignments := []*Assignment{}
// 	trimmedBaseURL := strings.TrimSuffix(c.BaseURL, "/api/v1")

// 	courses, err := c.GetCoursesByAccount(account, "student")
// 	if err != nil {
// 		return nil, terror.Error(err, "error retrieving courses")
// 	}

// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	var firstErr error

// 	for _, course := range courses {
// 		wg.Add(1)
// 		go func(course *Course) {
// 			defer wg.Done()
// 			_assignments, err := c.GetAssignmentsByCourseID(course.ID, "ungraded")
// 			if err != nil {
// 				mu.Lock()
// 				if firstErr == nil {
// 					firstErr = terror.Error(err, "error retrieving ungraded assignments")
// 				}
// 				mu.Unlock()
// 				return
// 			}
// 			for _, assignment := range _assignments {
// 				assignment.Account = account.Name
// 				assignment.CourseName = course.Name
// 				assignment.Status = bucket
// 				assignment.GradebookURL = fmt.Sprintf(`%s/courses/%d/gradebook`, trimmedBaseURL, course.ID)
// 			}

// 			mu.Lock()
// 			assignments = append(assignments, _assignments...)
// 			mu.Unlock()

// 		}(course)
// 	}
// 	wg.Wait()

// 	if firstErr != nil {
// 		return nil, firstErr
// 	}

// 	return assignments, nil
// }
