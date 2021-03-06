import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { useQuery, useMutation, gql } from "@apollo/client";
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
  const id = useParams().id;
  const [name, setName] = useState();
  const {  data } = useQuery(GET_SET_QUERY, {
    variables: { id: id },
    onCompleted: () => {
      setName(data.museumSet.name);
    }
  });

  const [updateSet] = useMutation(UPDATE_SET_QUERY);

  return (
    <div>
      <h1>Edit museum set with ID: {id}</h1>
      <form
        onSubmit={event => {
          updateSet({ variables: { id: id, input: { name: name } } });
          event.preventDefault();
        }}
      >
        <p>
          <label>Name:</label>
          <br />
          <input
            type="text"
            value={name}
            onChange={e => {
              setName(e.target.value);
            }}
          />
        </p>
        <input type="submit" value="send" />
      </form>
    </div>
  );
};

export default EditMuseumSet;
