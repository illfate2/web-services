import { ApolloClient, createHttpLink, InMemoryCache } from "@apollo/client";

import React from "react";

import { ApolloProvider } from "@apollo/client";

import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import { setContext } from "@apollo/client/link/context";
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
import { RegistrationForm } from "./RegistrationForm";
import { LoginForm } from "./LoginForm";
import { onError } from "@apollo/client/link/error";
import { Cookies, CookiesProvider } from "react-cookie";

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors)
    graphQLErrors.forEach(({ message, locations, path }) =>
      console.log(
        `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
      )
    );

  if (networkError) console.log(`[Network error]: ${networkError}`);
});

const httpLink = createHttpLink({
  uri: "http://localhost:8082/query",
  errorLink: errorLink
});

const authLink = setContext((_, { headers }) => {
  // get the authentication token from local storage if it exists
  const token = localStorage.getItem("accessToken");
  // return the headers tso the context so httpLink can read them
  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : ""
    }
  };
});

const client = new ApolloClient({
  link: authLink.concat(httpLink),
  cache: new InMemoryCache()
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
        <Route path="/register">
          <RegistrationForm />
        </Route>
        <Route path="/login">
          <LoginForm />
        </Route>
        <Route path="/">
          <CookiesProvider>
            <App />
          </CookiesProvider>
        </Route>
      </Switch>
    </ApolloProvider>
  </Router>,
  document.getElementById("root")
);
