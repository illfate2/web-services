import React from "react";
import { gql, useMutation } from "@apollo/client";
import { useForm } from "react-hook-form";
import { useHistory } from "react-router-dom";
const LOGIN_MUTATION = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      accessToken
      refreshToken
    }
  }
`;

export const LoginForm = () => {
  const history = useHistory();

  const [login] = useMutation(LOGIN_MUTATION, {
    onCompleted: data => {
      localStorage.setItem("accessToken", data.login.accessToken);
      localStorage.setItem("refreshToken", data.login.refreshToken);
      history.push("/");
    },
    onError: error => {
      console.log(error);
    }
  });

  const { register, handleSubmit } = useForm();
  const submit = data => {
    login({
      variables: {
        email: data.email,
        password: data.password
      }
    });
  };
  return (
    <form onSubmit={handleSubmit(submit)}>
      <h1>Login:</h1>
      <label>
        Email:
        <input type="text" ref={register} name="email" />
      </label>
      <label>
        Password:
        <input type="password" ref={register} name="password" />
      </label>
      <input type="submit" value="Login" />
      <a href="http://localhost:8082/login/google">Login with google</a>
    </form>
  );
};
