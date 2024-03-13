import { useNavigate } from "react-router-dom";

export default function HomeScreen() {
  const navigate = useNavigate();
  return (
    <div id="App">
      <div id="result" className="result"></div>
      <div id="input" className="input-box">
        <input
          id="name"
          className="input"
          autoComplete="off"
          name="input"
          type="text"
        />
        <button className="btn">Greet</button>
        <div
          style={{
            display: "flex",
            flexDirection: "column",
          }}
        >
          <span
            onClick={() => {
              navigate("/assignment");
            }}
          >
            Assignments
          </span>
          <span>Assignments</span>
        </div>
      </div>
    </div>
  );
}
