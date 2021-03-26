import {ApolloClient, ApolloProvider, createHttpLink, InMemoryCache} from "@apollo/client";

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
import {onError} from "@apollo/client/link/error";
import {setContext} from "@apollo/client/link/context";
import {LoginForm} from "./LoginForm";
import {RegistrationForm} from "./RegistrationForm";

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

const errorLink = onError(({graphQLErrors, networkError}) => {
    if (graphQLErrors)
        graphQLErrors.forEach(({message, locations, path}) =>
            console.log(
                `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
            )
        );

    if (networkError) console.log(`[Network error]: ${networkError}`);
});

const authLink = setContext((_, {headers}) => {
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


const httpLink = createHttpLink({
    uri: process.env.REACT_APP_API_URL + "/query",
    errorLink: errorLink
});

const client = new ApolloClient({
    cache: new InMemoryCache(),
    link: authLink.concat(httpLink),
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
                    <Route path="/login">
                        <LoginForm/>
                    </Route>
                    <Route path="/signup">
                        <RegistrationForm/>
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
