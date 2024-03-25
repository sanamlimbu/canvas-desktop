import { useState } from "react";
import {
  GetAccountByID,
  GetAssignmentsByCourse,
  GetCoursesByAccount,
} from "../../wailsjs/go/canvas/APIClient";
import { ExportAssignmentsStatus } from "../../wailsjs/go/main/App";
import { canvas } from "../../wailsjs/go/models";
import { Qualifications } from "../constant";

interface UngradedAssignmentsProps {
  inProgress: boolean;
  changeInProgress: (val: boolean) => void;
}

export default function UngradedAssignments({
  inProgress,
  changeInProgress,
}: UngradedAssignmentsProps) {
  const qualifications = Qualifications;
  const [accountID, setAccountID] = useState<number>(
    qualifications[0].AccountID
  );
  const [errorMsg, setErrorMsg] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  const [progress, setProgress] = useState(0);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    changeInProgress(true);
    setSuccessMsg("");
    setErrorMsg("");

    try {
      const assignments: canvas.Assignment[] = [];
      const account = await GetAccountByID(accountID);
      const courses = await GetCoursesByAccount(
        account,
        canvas.CourseEnrollmentType.STUDENT
      );
      let completedCourses = 0;
      const totalProgress = courses.length + 1; // One for CSV export operation

      for (let i = 0; i < courses.length; i++) {
        const _assignments = await GetAssignmentsByCourse(
          courses[i],
          canvas.AssignmentBucket.UNGRADED
        );

        assignments.push(..._assignments);
        completedCourses++;
        setProgress((completedCourses / totalProgress) * 100);
      }

      await ExportAssignmentsStatus(assignments, account);
      setProgress(100);
      setSuccessMsg("Successfully created 2 csv files in currrent folder.");
    } catch (err: any) {
      setErrorMsg(err);
    } finally {
      changeInProgress(false);
      setProgress(0);
    }
  };

  return (
    <div>
      <div style={{ marginBottom: "0.5em" }}>
        <label>Export ungraded assignments report</label>
      </div>
      <div style={{ maxWidth: "40rem", margin: "0 auto", padding: "0 1rem" }}>
        <form
          id="app-cover"
          onSubmit={handleSubmit}
          className="flex-column"
          style={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            alignItems: "center",
            maxWidth: "40rem",
            gap: "1em",
          }}
        >
          <div>
            <label>Select a qualification: </label>
            <select
              onChange={(e) => setAccountID(Number(e.target.value))}
              disabled={inProgress}
            >
              {qualifications.map((qualification) => (
                <option
                  key={qualification.AccountID}
                  value={qualification.AccountID}
                >
                  {qualification.Name}
                </option>
              ))}
            </select>
          </div>
          <button type="submit" disabled={inProgress}>
            Start
          </button>
        </form>
      </div>
      {inProgress && (
        <div style={{ marginTop: "0.5em" }}>
          <span>Completed </span>
          <progress value={progress} max={100} />
          <span> {Math.floor(progress)}%</span>
        </div>
      )}
      <div style={{ marginTop: "0.5em" }}>
        {errorMsg && <span style={{ color: "#ef5350" }}> {errorMsg}</span>}
        {successMsg && <span>{successMsg}</span>}
      </div>
    </div>
  );
}
