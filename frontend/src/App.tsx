import { useState } from "react";
import "./App.css";
import UngradedSubmissions from "./components/ungradedSubmissions";
import { Export } from "./types";

function App() {
  const [exportItem, setExportItem] = useState<Export>(
    Export.UngradedAssignments
  );
  const [inProgress, setInProgress] = useState(false);

  const changeInProgres = (val: boolean) => {
    setInProgress(val);
  };

  return (
    <div id="App">
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "center",
          gap: "1em",
          paddingTop: "1em",
          paddingBottom: "1em",
          borderBottom: "1px solid white",
          marginBottom: "1em",
        }}
      >
        <label>Please select an action:</label>
        <select
          onChange={(e) => {
            setExportItem(e.target.value as Export);
          }}
          disabled={inProgress}
        >
          <option value={Export.UngradedAssignments} style={{ padding: 10 }}>
            Export ungraded assignments
          </option>
          <option value={Export.StudentAssessments} style={{ padding: 10 }}>
            Export student assessments
          </option>
        </select>
      </div>
      {exportItem === Export.UngradedAssignments && (
        <UngradedSubmissions
          inProgress={inProgress}
          changeInProgress={changeInProgres}
        />
      )}
    </div>
  );
}

export default App;
