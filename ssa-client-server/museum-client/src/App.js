import logo from "./logo.svg";
import "./App.css";
import MuseumSets from "./MuseumSets";

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Museum App</h1>
        <Link to="/museumItems">Museum items </Link>
        <Link to="/museumSets">Museum sets </Link>
        <Link to="/museumFunds">Museum funds </Link>
        <Link to="/museumItemMovements">Museum movements </Link>
      </header>
    </div>
  );
}

export default App;
