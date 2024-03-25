import { Divider, Flex, Image, Select, Title } from "@mantine/core";
import { useState } from "react";
import SAILogo from "../assets/images/sai-logo.png";
import EnrollmentsResult from "../components/enrollmentsResult";
import ToggleTheme from "../components/toggleTheme";
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
          marginRight: "1em",
        }}
      >
        <Flex justify="space-between" mb="md">
          <Image src={SAILogo} height={100} fit="contain" width="auto" />
          <ToggleTheme />
        </Flex>

        <Title
          style={{
            background: "linear-gradient(to right, #007FFF, #0059B2)",
            WebkitBackgroundClip: "text",
            WebkitTextFillColor: "transparent",
            fontWeight: "bold",
          }}
          mb={"md"}
        >
          Canvas reports
        </Title>

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
        {exportItem && <Divider mt={"lg"} />}
      </div>
      {exportItem === Export.UngradedAssignments && (
        <UngradedSubmissions
          inProgress={inProgress}
          changeInProgress={changeInProgres}
        />
      )}
      {exportItem === Export.EnrollmentResults && (
        <EnrollmentsResult
          inProgress={inProgress}
          changeInProgress={changeInProgres}
        />
      )}
    </div>
  );
}

export default HomePage;
