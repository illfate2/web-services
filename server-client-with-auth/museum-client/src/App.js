import "./App.css";
import {Link, useHistory} from "react-router-dom";
import Cookies from "js-cookie";

function App() {
    const authCookie = Cookies.get("auth");
    const history = useHistory();
    if (authCookie !== undefined) {
        localStorage.setItem("accessToken", authCookie);
    } else {
        history.push("/login");
    }
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
