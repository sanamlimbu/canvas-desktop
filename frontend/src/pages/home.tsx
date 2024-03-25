import { Select } from "@mantine/core";
import { useState } from "react";
import StudentEnrollmentsResult from "../components/studentEnrollmentsResult";
import UngradedSubmissions from "../components/ungradedAssignments";
import { Export } from "../types";

function HomePage() {
  const [exportItem, setExportItem] = useState<string | null>(null);
  const [inProgress, setInProgress] = useState(false);

  const changeInProgres = (val: boolean) => {
    setInProgress(val);
  };

  return (
    <div
      id="App"
      style={{
        display: "flex",
        justifyContent: "center",
        flexDirection: "column",
        maxWidth: "40em",
        marginLeft: "auto",
        marginRight: "auto",
      }}
    >
      <div
        style={{
          marginTop: "1em",
          marginBottom: "1em",
          maxWidth: "24em",
        }}
      >
        <Select
          label="Please select an action."
          placeholder="Pick an action"
          value={exportItem}
          data={Object.values(Export)}
          onChange={setExportItem}
          clearable
          searchable
          disabled={inProgress}
        />
      </div>
      {exportItem === Export.UngradedAssignments && (
        <UngradedSubmissions
          inProgress={inProgress}
          changeInProgress={changeInProgres}
        />
      )}
      {exportItem === Export.StudentEnrollmentResults && (
        <StudentEnrollmentsResult
          inProgress={inProgress}
          changeInProgress={changeInProgres}
        />
      )}
    </div>
  );
}

export default HomePage;
