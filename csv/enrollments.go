package csv

import (
	"canvas-desktop/canvas"
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/ninja-software/terror/v2"
)

func ExportEnrollmentsResults(results []*canvas.EnrollmentResult, userSisID string) error {
	time := time.Now().Format("2006-01-02-15-04-05")
	file, err := os.Create(fmt.Sprintf("%s-%s-enrollments_results.csv", userSisID, time))
	if err != nil {
		return err
	}
	defer file.Close()

	err = gocsv.MarshalFile(&results, file)
	if err != nil {
		return terror.Error(err, "cannot write rows to csv file")
	}

	return nil
}
