import { Alert, Button, Progress, Select, Text } from "@mantine/core";
import { IconCheck, IconInfoCircle } from "@tabler/icons-react";
import { useState } from "react";
import {
  GetAccountByID,
  GetAssignmentsByCourse,
  GetCoursesByAccount,
} from "../../wailsjs/go/canvas/APIClient";
import { ExportAssignmentsStatus } from "../../wailsjs/go/main/App";
import { canvas } from "../../wailsjs/go/models";
import { Qualifications, QualificationsWithAccountID } from "../constant";
import { colors } from "../theme";

interface UngradedAssignmentsProps {
  inProgress: boolean;
  changeInProgress: (val: boolean) => void;
}

export default function UngradedAssignments({
  inProgress,
  changeInProgress,
}: UngradedAssignmentsProps) {
  const [accountID, setAccountID] = useState<number | null>(null);
  const [errorMsg, setErrorMsg] = useState("");
  const [successMsg, setSuccessMsg] = useState("");
  const [progress, setProgress] = useState(0);
  const iconCheck = <IconCheck />;
  const iconInfoCircle = <IconInfoCircle />;

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    changeInProgress(true);
    setSuccessMsg("");
    setErrorMsg("");

    if (accountID === null) {
      return;
    }

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
      setSuccessMsg("Created 2 csv files in currrent folder.");
    } catch (err: any) {
      setErrorMsg(err);
    } finally {
      changeInProgress(false);
      setProgress(0);
    }
  };

  return (
    <div style={{ maxWidth: "24em" }}>
      <Text fw={500} c={colors.blue}>
        Export ungraded assignments report
      </Text>
      <form onSubmit={handleSubmit}>
        <Select
          label="Select a qualification."
          placeholder="Pick a qualification"
          data={Qualifications}
          onChange={(val) => {
            const accountID = val
              ? QualificationsWithAccountID.get(val)
              : undefined;
            accountID ? setAccountID(accountID) : null;
          }}
          mb={"md"}
          disabled={inProgress}
          searchable
          clearable
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
        <div>
          <Progress
            color="teal"
            radius="md"
            value={progress}
            striped
            animated
            mt={"md"}
          />
          <Text c="blue" fw={700} ta="center" size="md">
            {Math.floor(progress)}%
          </Text>
        </div>
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
