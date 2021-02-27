import React from "react";
import { gql } from "@apollo/client";
import TableContainer from "./TableContainer";
import "./popup.css";
import { useHistory } from "react-router-dom";
import {
  useMutationWithAuthErrHandling,
  useQueryWithAuthErrHandling
} from "./Queries";

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

function CreateSetForm() {
  const history = useHistory();

  let input;
  const [addSet] = useMutationWithAuthErrHandling(CREATE_SET_QUERY);
  return (
    <div>
      <form
        onSubmit={e => {
          addSet({ variables: { input: { name: input.value } } });
          input.value = "";
        }}
      >
        <label>Set name:</label>
        <input
          ref={node => {
            input = node;
          }}
        />
        <button type="submit">Create</button>
      </form>
    </div>
  );
}

const MuseumSets = () => {
  const history = useHistory();

  let { loading, error, data } = useQueryWithAuthErrHandling(GET_SETS_QUERY);

  const [deleteSet] = useMutationWithAuthErrHandling(DELETE_SET_QUERY);

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
              window.location.reload(false);
            }}
            value={"remove"}
          >
            {"remove"}
          </button>
        )
      }
    ],
    []
  );

  if (loading) return "Loading...";
  if (error) return "error";

  return (
    <div>
      <CreateSetForm />
      <TableContainer columns={columns} data={data.museumSets} />
    </div>
  );
};

export default MuseumSets;
