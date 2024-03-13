import { useState } from "react";
import "./App.css";
import UngradedSubmissions from "./components/ungradedSubmissions";
import { Export } from "./types";

function App() {
  const [exportItem, setExportItem] = useState<Export>(
    Export.UngradedAssignments
  );
  const [value, setValue] = useState("");
  const updateExportItem = (item: Export) => {
    setExportItem(item);
  };

  function greet() {}

  return (
    <div id="App">
      <div
        style={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "center",
          gap: "1em",
          marginTop: "1em",
        }}
      >
        <label>Please choose an action item</label>
        <select
          onChange={(e) => {
            setExportItem(e.target.value as Export);
          }}
        >
          <option value={Export.UngradedAssignments}>
            Export ungraded assignments
          </option>
          <option value={Export.StudentAssessments}>
            Export student assessments
          </option>
        </select>
      </div>
      {exportItem === Export.UngradedAssignments && <UngradedSubmissions />}
    </div>
  );
}

export default App;
