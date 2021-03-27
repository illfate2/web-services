import React from "react";
import {gql, useMutation} from "@apollo/client";
import {useForm} from "react-hook-form";
import {Link, useHistory} from "react-router-dom";
import './Login.css';

const LOGIN_MUTATION = gql`
    mutation Login($email: String!, $password: String!) {
        login(email: $email, password: $password) {
            accessToken
            refreshToken
        }
    }
`;

export const LoginForm = (props) => {
    const history = useHistory();

    const [login] = useMutation(LOGIN_MUTATION, {
        onCompleted: data => {
            localStorage.setItem("accessToken", data.login.accessToken);
            localStorage.setItem("refreshToken", data.login.refreshToken);
            history.push("/");
        },
        onError: error => {
            alert(error);
        }
    });

    const {register, handleSubmit} = useForm();
    const submit = data => {
        login({
            variables: {
                email: data.email,
                password: data.password
            }
        });
    };
    return (
        <form onSubmit={handleSubmit(submit)} className="Login">
            <h1>Login:</h1>
            <label className="Login">
                Email:
                <input type="email" ref={register} name="email" className="Login"/>
            </label>
            <label className="Login">
                Password:
                <input type="password" ref={register} name="password" className="Login"/>
            </label>
            <input type="submit" className="LoginButton"/>
            <a href={props.base_uri + "/login/google"} className="Login">Login with google</a>
            <a href={props.base_uri + "/login/discord"} className="Login">Login with discord</a>
            <a href={props.base_uri + "/login/github"} className="Login">Login with github</a>
            <Link to="/signup">Register</Link>
        </form>
    );
};
