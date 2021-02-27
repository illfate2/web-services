import React, { useState } from "react";
import { useParams } from "react-router-dom";
import { useQuery, useMutation, gql } from "@apollo/client";


const GET_FUND_QUERY = gql`
  query MuseumFund($id: ID!) {
    museumFund(id: $id) {
      id
      name
    }
  }
`;

const UPDATE_FUND_QUERY = gql`
  mutation UpdateMuseumFund($id: ID!, $input: UpdateMuseumFundInput!) {
    updateMuseumFund(id: $id, input: $input) {
      id
      name
    }
  }
`;

const EditMuseumFund = () => {
  const id = useParams().id;
  const [name, setName] = useState();
  const { loading, error, data } = useQuery(GET_FUND_QUERY, {
    variables: { id: id },
    onCompleted: () => {
      setName(data.museumFund.name);
    }
  });

  const [updateFund, { updateData }] = useMutation(UPDATE_FUND_QUERY);

  return (
    <div>
      <h1>Edit museum fund with ID: {id}</h1>
      <form
        onSubmit={event => {
          updateFund({ variables: { id: id, input: { name: name } } });
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

export default EditMuseumFund;
