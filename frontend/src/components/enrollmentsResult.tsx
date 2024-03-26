import { Alert, Button, Flex, Loader, Text, TextInput } from "@mantine/core";
import { IconCheck, IconInfoCircle } from "@tabler/icons-react";
import { useState } from "react";
import {
  GetEnrollmentResultsByUser,
  GetUserBySisID,
} from "../../wailsjs/go/canvas/APIClient";
import { ExportEnrollmentsResults } from "../../wailsjs/go/main/App";
import { colors } from "../theme";

interface EnrollmentsResultProps {
  inProgress: boolean;
  changeInProgress: (val: boolean) => void;
}

export default function EnrollmentsResult({
  inProgress,
  changeInProgress,
}: EnrollmentsResultProps) {
  const [errorMsg, setErrorMsg] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  const iconCheck = <IconCheck />;
  const iconInfoCircle = <IconInfoCircle />;
  const [sisID, setSisID] = useState("");

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (sisID === "") {
      setErrorMsg("Please enter a valid SIS ID.");
      return;
    }

    changeInProgress(true);
    setSuccessMsg("");
    setErrorMsg("");

    try {
      const student = await GetUserBySisID(sisID); // not found will throw err "404"
      const results = await GetEnrollmentResultsByUser(student);
      await ExportEnrollmentsResults(results, student.sis_user_id);
      setSuccessMsg("Created a csv file in currrent folder.");
    } catch (err: any) {
      setErrorMsg(err);
    } finally {
      changeInProgress(false);
    }
  };

  return (
    <div>
      <Text fw={500} c={colors.blue}>
        Export enrollments result
      </Text>
      <form onSubmit={handleSubmit}>
        <TextInput
          label="Please enter SIS ID."
          placeholder="SIS ID is case sensitive."
          mb={"md"}
          onChange={(event) => setSisID(event.currentTarget.value)}
          disabled={inProgress}
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
      {inProgress && (
        <Flex align={"center"} justify={"center"}>
          <Loader mt={"md"} size={"sm"} />
        </Flex>
      )}
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
          title="Success"
          icon={iconCheck}
          mt={"md"}
        >
          {successMsg}
        </Alert>
      )}
    </div>
  );
}
