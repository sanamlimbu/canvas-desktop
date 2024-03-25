import { Alert, Button, Loader, Text, TextInput } from "@mantine/core";
import { IconCheck, IconInfoCircle } from "@tabler/icons-react";
import { useRef, useState } from "react";
import {
  GetEnrollmentResultsByUser,
  GetUserBySisID,
} from "../../wailsjs/go/canvas/APIClient";
import { ExportEnrollmentsResults } from "../../wailsjs/go/main/App";

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
  const iconCheck = <IconCheck />;
  const iconInfoCircle = <IconInfoCircle />;

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
    <div style={{ maxWidth: "24em" }}>
      <Text fw={500} c="blue">
        Export ungraded assignments report
      </Text>
      <form onSubmit={handleSubmit}>
        <TextInput
          label="Please enter SIS ID."
          placeholder="SIS ID is case sensitive."
          mb={"md"}
        />
        <Button
          type="submit"
          disabled={inProgress}
          variant="outline"
          color="cyan"
        >
          Start
        </Button>
      </form>
      {inProgress && <Loader mt={"md"} />}
      {errorMsg && (
        <Alert
          variant="light"
          color="red"
          title="Error"
          icon={iconInfoCircle}
          mt={"md"}
        >
          {errorMsg}
        </Alert>
      )}
      {successMsg && (
        <Alert
          variant="light"
          color="teal"
          title="Successful"
          icon={iconCheck}
          mt={"md"}
        >
          Created a csv file in currrent folder.
        </Alert>
      )}
    </div>
  );
}
