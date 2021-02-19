import { ApolloClient, InMemoryCache } from "@apollo/client";

import React from "react";

import { ApolloProvider } from "@apollo/client";

import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

import EditMuseumSet from "./EditMuseumSet";
import MuseumSets from "./MuseumSets";
import EditMuseumFund from "./EditMuseumFund";
import MuseumFunds from "./MuseumFunds";
import ViewMuseumItemMovement from "./ViewMuseumItemMovement";
import EditMuseumItemMovement from "./EditMuseumItemMovement";
import MuseumItemMovements from "./MuseumItemMovements";
import { EditMuseumItem } from "./EditMuseumItem";
import { ViewMuseumItem } from "./ViewMuseumItem";
import MuseumItems from "./MuseumItems";
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
        <Route path="/museumFund/edit/:id">
          <EditMuseumFund />
        </Route>
        <Route path="/museumFunds">
          <MuseumFunds />
        </Route>
        <Route path="/museumItem/view/:id">
          <ViewMuseumItem />
        </Route>
        <Route path="/museumItem/edit/:id">
          <EditMuseumItem />
        </Route>
        <Route path="/museumItems">
          <MuseumItems />
        </Route>
        <Route path="/museumItemMovement/view/:id">
          <ViewMuseumItemMovement />
        </Route>
        <Route path="/museumItemMovement/edit/:id">
          <EditMuseumItemMovement />
        </Route>
        <Route path="/museumItemMovements">
          <MuseumItemMovements />
        </Route>
        <Route path="/">
          <App />
        </Route>
      </Switch>
    </ApolloProvider>
  </Router>,
  document.getElementById("root")
);
