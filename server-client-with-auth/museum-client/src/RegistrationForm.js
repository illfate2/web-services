import React from "react";
import { gql, useMutation } from "@apollo/client";
import { useForm } from "react-hook-form";
const SIGNUP_MUTATION = gql`
  mutation SignupUser($input: Signup!) {
    signupUser(input: $input) {
      accessToken
      refreshToken
    }
  }
`;

export const RegistrationForm = () => {
  const [signup] = useMutation(SIGNUP_MUTATION, {
    onCompleted: data => {
      localStorage.setItem("accessToken", data.signupUser.accessToken);
      localStorage.setItem("refreshToken", data.signupUser.refreshToken);
    },
    onError: error => {
      console.log(error);
    }
  });

  const { register, handleSubmit } = useForm();
  const submit = data => {
    signup({
      variables: {
        input: {
          email: data.email,
          password: data.password
        }
      }
    });
  };
  return (
    <form onSubmit={handleSubmit(submit)}>
      <label>
        Email:
        <input type="text" ref={register} name="email" />
      </label>
      <label>
        Password:
        <input type="password" ref={register} name="password" />
      </label>
      <input type="submit" value="Отправить" />
    </form>
  );
};
