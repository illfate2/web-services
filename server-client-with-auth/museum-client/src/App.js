import logo from "./logo.svg";
import "./App.css";
import MuseumSets from "./MuseumSets";
import Cookies from "js-cookie";

function App() {
  const authCookie = Cookies.get("auth");
  if (authCookie !== undefined) {
    localStorage.setItem("accessToken", authCookie);
  }
  return (
    <div className="App">
      <header className="App-header">
        <h1>Museum App</h1>
        <a href="/museumItems">Museum items </a>
        <a href="/museumSets">Museum sets </a>
        <a href="/museumFunds">Museum funds </a>
        <a href="/museumItemMovements">Museum movements </a>
      </header>
    </div>
  );
}

export default App;
