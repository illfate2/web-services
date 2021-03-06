import React, { useState } from "react";
import { useQuery, gql, useMutation } from "@apollo/client";
import TableContainer from "./TableContainer";
import "./popup.css";
import { useHistory } from "react-router-dom";
import { useForm } from "react-hook-form";

const CREATE_SET_QUERY = gql`
  mutation CreateMuseumSet($input: MuseumSetInput!) {
    createMuseumSet(input: $input) {
      id
      name
    }
  }
`;

const DELETE_SET_QUERY = gql`
  mutation DeleteMuseumSet($id: ID!) {
    deleteMuseumSet(id: $id)
  }
`;

const GET_SETS_QUERY = gql`
  {
    museumSets {
      id
      name
    }
  }
`;

function CreateSetForm({ onSubmit }) {
  const { register, handleSubmit } = useForm();
  let input;
  return (
    <div>
      <form onSubmit={handleSubmit(onSubmit)}>
        <label>Set name:</label>
        <input name="set" ref={register} />
        <button type="submit">Create</button>
      </form>
    </div>
  );
}

const MuseumSets = () => {
  const [museumSetsData, setMuseumSetsData] = useState([]);

  const { loading, data } = useQuery(GET_SETS_QUERY, {
    onCompleted: data => {
      setMuseumSetsData(data.museumSets);
    }
  });

  const [addSet] = useMutation(CREATE_SET_QUERY, {
    onCompleted: data => {
      let dataCopy = [...museumSetsData];
      dataCopy.push({
        id: data.createMuseumSet.id,
        name: data.createMuseumSet.name
      });
      setMuseumSetsData(dataCopy);
    }
  });
  const onCreateSetSumbit = input => {
    addSet({ variables: { input: { name: input.set } } });
  };

  const [deleteSet] = useMutation(DELETE_SET_QUERY);

  const history = useHistory();

  const handleClick = id => {
    history.push("/museumSet/edit/" + id);
  };

  const columns = React.useMemo(
    () => [
      {
        Header: "id",
        accessor: "id",
        show: false
      },
      {
        Header: "name",
        accessor: "name"
      },
      {
        Header: "edit",
        accessor: "edit",
        Cell: ({ row }) => (
          <button
            onClick={() => {
              handleClick(row.original.id);
            }}
            value={"edit"}
          >
            {"edit"}
          </button>
        )
      },
      {
        Header: "delete",
        accessor: "delete",
        Cell: ({ row }) => (
          <button
            onClick={() => {
              deleteSet({ variables: { id: row.original.id } });
              const dataCopy = [...museumSetsData];
              dataCopy.splice(row.index, 1);
              setMuseumSetsData(dataCopy);
            }}
            value={"remove"}
          >
            {"remove"}
          </button>
        )
      }
    ],
    [museumSetsData]
  );

  if (loading) return "Loading...";

  return (
    <div>
      <CreateSetForm onSubmit={onCreateSetSumbit} />
      <TableContainer columns={columns} data={museumSetsData} />
    </div>
  );
};

export default MuseumSets;
