package csv

import (
	"canvas-desktop/canvas"
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/ninja-software/terror/v2"
)

func ExportUngradedSubmissions(submissions []*canvas.Submission, account *canvas.Account) error {
	time := time.Now().Format("2006-01-02-15-04-05")
	file, err := os.Create(fmt.Sprintf("%d-%s-ungraded_submissions.csv", account.ID, time))
	if err != nil {
		return err
	}
	defer file.Close()

	err = gocsv.MarshalFile(&submissions, file)
	if err != nil {
		return terror.Error(err, "cannot write rows to csv file")
	}

	return nil
}
