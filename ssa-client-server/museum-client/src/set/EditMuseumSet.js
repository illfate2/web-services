import React, {useContext, useState} from "react";
import {useHistory, useParams} from "react-router-dom";
import {gql, useMutation, useQuery} from "@apollo/client";
import {Context} from "../store/Store";


const GET_SET_QUERY = gql`
    query MuseumSet($id: ID!) {
        museumSet(id: $id) {
            id
            name
        }
    }
`;

const UPDATE_SET_QUERY = gql`
    mutation UpdateMuseumSet($id: ID!, $input: MuseumSetInput!) {
        updateMuseumSet(id: $id, input: $input) {
            id
            name
        }
    }
`;

const EditMuseumSet = () => {
    const [state, dispatch] = useContext(Context);
    const id = useParams().id;
    const [name, setName] = useState();
    const history = useHistory();
    const {data} = useQuery(GET_SET_QUERY, {
        variables: {id: id},
        onCompleted: () => {
            setName(data.museumSet.name);
        }
    });

    const [updateSet] = useMutation(UPDATE_SET_QUERY, {
        onCompleted: data => {
            dispatch({ type: "UPDATE_SET", payload: data.updateMuseumSet });
        }
    });

    return (
        <div>
            <h1>Edit museum set with ID: {id}</h1>
            <form
                onSubmit={event => {
                    updateSet({variables: {id: id, input: {name: name}}});
                    event.preventDefault();
                    history.push("/museumSets");
                }}
            >
                <p>
                    <label>Name:</label>
                    <br/>
                    <input
                        type="text"
                        value={name}
                        onChange={e => {
                            setName(e.target.value);
                        }}
                    />
                </p>
                <input type="submit" value="send"/>
            </form>
        </div>
    );
};

export default EditMuseumSet;
