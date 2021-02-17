import React from "react";
import MuseumSet from "./MuseumSet";
import { useQuery, gql, useMutation, useLazyQuery } from "@apollo/client";
import TableContainer from "./TableContainer";
import "./popup.css";
import { useHistory } from "react-router-dom";

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
  let input;
  const [addSet, { data }] = useMutation(CREATE_SET_QUERY);
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

function UpdateMuseumSet() {}

const MuseumSets = () => {
  let { loading, error, data } = useQuery(GET_SETS_QUERY);

  const [deleteSet, { deleteData }] = useMutation(DELETE_SET_QUERY);

  const [getSet, { setData }] = useLazyQuery(GET_SETS_QUERY);

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

  return (
    <div>
      <CreateSetForm />
      <TableContainer columns={columns} data={data.museumSets} />
    </div>
  );
};

export default MuseumSets;
