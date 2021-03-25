import React, { useContext } from "react";
import { useQuery, gql, useMutation } from "@apollo/client";
import TableContainer from "../table/TableContainer";
import { useHistory } from "react-router-dom";
import { useForm } from "react-hook-form";
import { Context } from "../store/Store";

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
  const [state, dispatch] = useContext(Context);
  const { loading } = useQuery(GET_SETS_QUERY, {
    onCompleted: data => {
      dispatch({ type: "SET_SETS", payload: data.museumSets });
    }
  });

  const [addSet, { error }] = useMutation(CREATE_SET_QUERY, {
    onCompleted: data => {
      if (!error) {
        dispatch({ type: "ADD_SET", payload: data.createMuseumSet });
      }
    }
  });
  const onCreateSetSubmit = input => {
    addSet({ variables: { input: { name: input.set } } });
  };

  const [deleteSet] = useMutation(DELETE_SET_QUERY, {
    onCompleted: data => {}
  });

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
              dispatch({ type: "REMOVE_SET", payload: row.original.id });
            }}
            value={"remove"}
          >
            {"remove"}
          </button>
        )
      }
    ],
    [state.sets]
  );

  if (loading) return "Loading...";

  return (
    <div>
      <CreateSetForm onSubmit={onCreateSetSubmit} />
      <TableContainer columns={columns} data={state.sets} />
    </div>
  );
};

export default MuseumSets;
