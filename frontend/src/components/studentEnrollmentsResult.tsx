import { useRef, useState } from "react";
import {
  GetEnrollmentResultsByUser,
  GetUserBySisID,
} from "../../wailsjs/go/canvas/APIClient";
import { ExportEnrollmentsResults } from "../../wailsjs/go/main/App";
import "../App.css";
import Loader from "./loader";

interface StudentEnrollmentsResultProps {
  inProgress: boolean;
  changeInProgress: (val: boolean) => void;
}

export default function StudentEnrollmentsResult({
  inProgress,
  changeInProgress,
}: StudentEnrollmentsResultProps) {
  const [errorMsg, setErrorMsg] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  const sisIdInput = useRef<HTMLInputElement>(null);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!sisIdInput.current?.value) {
      setErrorMsg("Please enter a valid SIS ID.");
      return;
    }

    changeInProgress(true);
    setSuccessMsg("");
    setErrorMsg("");

    try {
      const student = await GetUserBySisID(sisIdInput.current?.value); // not found will throw err "404"
      const results = await GetEnrollmentResultsByUser(student);
      console.log("total ", results.length);
      await ExportEnrollmentsResults(results, student.sis_user_id);
      setSuccessMsg("Successfully created a csv file in currrent folder.");
    } catch (err: any) {
      setErrorMsg(err);
    } finally {
      changeInProgress(false);
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
            <label>Please enter student's SIS ID: </label>
            <input type="text" ref={sisIdInput} />
          </div>
          <button type="submit" disabled={inProgress}>
            Start
          </button>
        </form>
      </div>

      {inProgress && (
        <div style={{ marginTop: "0.5em" }}>
          <Loader />
        </div>
      )}

      <div style={{ marginTop: "0.5em" }}>
        {errorMsg && <span style={{ color: "#ef5350" }}> {errorMsg}</span>}
        {successMsg && <span>{successMsg}</span>}
      </div>
    </div>
  );
}
