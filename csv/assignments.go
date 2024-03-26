package csv

import (
	"fmt"
	"os"
	"strings"
	"time"

	"canvas-desktop/canvas"

	"github.com/gocarina/gocsv"
	"github.com/ninja-software/terror/v2"
)

func ExportAssignmentsResults(results []*canvas.AssignmentResult, userSisID string) error {
	time := time.Now().Format("2006-01-02-15-04-05")
	file, err := os.Create(fmt.Sprintf("%s-%s-assignment_results.csv", userSisID, time))
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

func ExportAssignmentsStatus(assignments []*canvas.Assignment, account *canvas.Account) error {
	time := time.Now().Format("2006-01-02-15-04-05")
	name := canvas.ReplaceSpaceInStr(account.Name, "_")

	filePerth, err := os.Create(fmt.Sprintf("PERTH-%s-%s-assignments_status.csv", name, time))
	if err != nil {
		return err
	}
	defer filePerth.Close()

	fileAdelaide, err := os.Create(fmt.Sprintf("ADL-%s-%s-assignments_status.csv", name, time))
	if err != nil {
		return err
	}
	defer fileAdelaide.Close()

	assignmentsPerth := []*canvas.Assignment{}
	assignmentsAdelaide := []*canvas.Assignment{}

	for _, assignment := range assignments {
		if strings.Contains(assignment.Section, "ADL") ||
			strings.Contains(assignment.Section, "ADELAIDE") ||
			strings.Contains(assignment.Section, "Adelaide") {
			assignmentsAdelaide = append(assignmentsAdelaide, assignment)
		} else {
			assignmentsPerth = append(assignmentsPerth, assignment)
		}
	}

	err = gocsv.MarshalFile(&assignmentsAdelaide, fileAdelaide)
	if err != nil {
		return terror.Error(err, "cannot write rows to csv file")
	}

	err = gocsv.MarshalFile(&assignmentsPerth, filePerth)
	if err != nil {
		return terror.Error(err, "cannot write rows to csv file")
	}

	return nil
}

func ExportAssignmentGradingStandards(assignments []*canvas.AssignmentGradingStandard) error {
	time := time.Now().Format("2006-01-02-15-04-05")
	file, err := os.Create(fmt.Sprintf("EB_assignments-%s.csv", time))
	if err != nil {
		return err
	}
	defer file.Close()

	err = gocsv.MarshalFile(&assignments, file)
	if err != nil {
		return terror.Error(err, "cannot write rows to csv file")
	}

	return nil
}
