import { MantineProvider } from "@mantine/core";
import "@mantine/core/styles.css";
import HomePage from "./pages/home";
import { theme } from "./theme";

export default function App() {
  return (
    <MantineProvider theme={theme}>
      <HomePage />
    </MantineProvider>
  );
}
