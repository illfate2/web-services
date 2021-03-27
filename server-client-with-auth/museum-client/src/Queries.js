import {fromPromise, gql, useMutation, useQuery} from "@apollo/client";
import {useHistory} from "react-router-dom";

function useMutationWithAuthErrHandling(query, options) {
    const history = useHistory();
    if (options !== undefined) {
        options.onError = error => {
            if (
                error.message === "Token is expired" ||
                error.message === "empty token" || error.message === "signature is invalid"
            ) {
                history.push("/login");
            }
        }
    }
    return useMutation(query, options);
}

function useQueryWithAuthErrHandling(query, options) {
    const history = useHistory();
    options.onError = error => {
        if (
            error.message === "Token is expired" ||
            error.message === "empty token" || error.message === "signature is invalid"
        ) {
            history.push("/login");
        }
    }
    return useQuery(query, options);
}

export {useMutationWithAuthErrHandling, useQueryWithAuthErrHandling};
