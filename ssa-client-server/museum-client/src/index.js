import {ApolloClient, ApolloProvider, InMemoryCache} from "@apollo/client";

import React from "react";

import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";

import EditMuseumSet from "./set/EditMuseumSet";
import MuseumSets from "./set/MuseumSets";
import EditMuseumFund from "./fund/EditMuseumFund";
import MuseumFunds from "./fund/MuseumFunds";
import ViewMuseumItemMovement from "./movement/ViewMuseumItemMovement";
import EditMuseumItemMovement from "./movement/EditMuseumItemMovement";
import MuseumItemMovements from "./movement/MuseumItemMovements";
import {EditMuseumItem} from "./item/EditMuseumItem";
import {ViewMuseumItem} from "./item/ViewMuseumItem";
import MuseumItems from "./item/MuseumItems";
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import Store from "./store/Store";

const defaultOptions = {
    watchQuery: {
        fetchPolicy: 'no-cache',
        errorPolicy: 'ignore',
    },
    query: {
        fetchPolicy: 'no-cache',
        errorPolicy: 'all',
    },
}

const client = new ApolloClient({
    cache: new InMemoryCache(),
    uri: process.env.REACT_APP_API_URL + "/query",
    // uri: "http://localhost:8082/query",
    defaultOptions: defaultOptions
});

ReactDOM.render(
    <Store>
        <Router>
            <ApolloProvider client={client}>
                <Switch>
                    <Route path="/museumSet/edit/:id">
                        <EditMuseumSet/>
                    </Route>
                    <Route path="/museumSets">
                        <MuseumSets/>
                    </Route>
                    <Route path="/museumFund/edit/:id">
                        <EditMuseumFund/>
                    </Route>
                    <Route path="/museumFunds">
                        <MuseumFunds/>
                    </Route>
                    <Route path="/museumItem/view/:id">
                        <ViewMuseumItem/>
                    </Route>
                    <Route path="/museumItem/edit/:id">
                        <EditMuseumItem/>
                    </Route>
                    <Route path="/museumItems">
                        <MuseumItems/>
                    </Route>
                    <Route path="/museumItemMovement/view/:id">
                        <ViewMuseumItemMovement/>
                    </Route>
                    <Route path="/museumItemMovement/edit/:id">
                        <EditMuseumItemMovement/>
                    </Route>
                    <Route path="/museumItemMovements">
                        <MuseumItemMovements/>
                    </Route>
                    <Route path="/">
                        <App/>
                    </Route>
                </Switch>
            </ApolloProvider>
        </Router>
    </Store>,
    document.getElementById("root")
);
