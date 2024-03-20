package main

import (
	"canvas-desktop/canvas"
	"canvas-desktop/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	baseURL := getenv("CANVAS_BASE_URL", "https://skillsaustralia.instructure.com/api/v1")
	accessToken := getenv("CANVAS_ACCESS_TOKEN", "18033~qQ5yZwSntWRrwRv4fOH2Y7ZnfmCN3jfpsnXLlzABxPeZoAkejTxQFNyrSN4DSxnq")
	pageSizeStr := getenv("CANVAS_PAGE_SIZE", "100")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		panic(err.Error)
	}

	if accessToken == "" {
		panic("missing access token")
	}

	//rl := rate.NewLimiter(rate.Every(10*time.Second), 1000) // 1000 requests every 10 seconds
	client := canvas.NewAPIClient(baseURL, accessToken, pageSize, http.DefaultClient)

	accountID := 111

	// courses, err := client.GetCoursesByAccountID(133, StudentCourseEnrollment)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for i, course := range courses {
	// 	fmt.Println(i+1, " ", course.Name, " ", course.Account.Name)
	// }

	// for _, qualification := range canvas.Qualifications {

	// }

	// submissions, err := client.GetUngradedSubmissionsByAccount(account)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	account, err := client.GetAccountByID(accountID)
	if err != nil {
		log.Fatal(err)
	}

	assignments, err := client.GetAssignmentsByAccount(account, "ungraded")
	if err != nil {
		log.Fatal(err)
	}

	err = csv.ExportAssignmentsStatus(assignments, account)
	if err != nil {
		fmt.Println("Failed exporting assignments status")
		return
	}

	// err = csv.ExportUngradedSubmissions(submissions, account)
	// if err != nil {
	// 	fmt.Println("Failed exporting ungraded submissions")
	// }

	fmt.Print("Successfully exported assignments status")
}

func getenv(key string, other string) string {
	value := os.Getenv(key)
	if value == "" {
		return other
	}

	return value
}
