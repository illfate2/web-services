import { ApolloClient, InMemoryCache } from "@apollo/client";

import React from "react";

import { ApolloProvider } from "@apollo/client";

import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

import EditMuseumSet from "./EditMuseumSet";
import MuseumSets from "./MuseumSets"
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";

const client = new ApolloClient({
  cache: new InMemoryCache(),
  uri: "http://localhost:8082/query"
});

ReactDOM.render(
  <Router>
    <ApolloProvider client={client}>
      <Switch>
        <Route path="/museumSet/edit/:id">
          <EditMuseumSet />
        </Route>
        <Route path="/museumSets">
          <MuseumSets />
        </Route>
        <Route path="/">
          <App />
        </Route>
      </Switch>
    </ApolloProvider>
  </Router>,
  document.getElementById("root")
);
