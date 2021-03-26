import {useMutation, useQuery} from "@apollo/client";
import {useHistory} from "react-router-dom";

function useMutationWithAuthErrHandling(query) {
    const history = useHistory();
    return useMutation(query, {
        onError: error => {
            if (
                error.message === "Token is expired" ||
                error.message === "empty token" || error.message === "signature is invalid"
            ) {
                history.push("/login");
            }
        }
    });
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
