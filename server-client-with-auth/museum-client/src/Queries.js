import { useQuery, gql, useMutation, useLazyQuery } from "@apollo/client";
import { useHistory } from "react-router-dom";

function useMutationWithAuthErrHandling(query) {
  const history = useHistory();
  return useMutation(query, {
    onError: error => {
      if (
        error.message == "Token is expired" ||
        error.message == "empty token"
      ) {
        history.push("/login");
      }
    }
  });
}

function useQueryWithAuthErrHandling(query) {
  const history = useHistory();
  return useQuery(query, {
    onError: error => {
      if (
        error.message == "Token is expired" ||
        error.message == "empty token"
      ) {
        history.push("/login");
      }
    }
  });
}

export { useMutationWithAuthErrHandling, useQueryWithAuthErrHandling };
