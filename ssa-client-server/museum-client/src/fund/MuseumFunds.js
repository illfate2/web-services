import React, { useState, useContext } from "react";
import { useQuery, gql, useMutation } from "@apollo/client";
import TableContainer from "../table/TableContainer";
import { useHistory } from "react-router-dom";
import { useForm } from "react-hook-form";
import { Context } from "../store/Store";

const CREATE_FUND_QUERY = gql`
  mutation CreateMuseumFund($input: MuseumFundInput!) {
    createMuseumFund(input: $input) {
      id
      name
    }
  }
`;

const DELETE_FUND_QUERY = gql`
  mutation DeleteMuseumFund($id: ID!) {
    deleteMuseumFund(id: $id)
  }
`;

const GET_FUNDS_QUERY = gql`
  {
    museumFunds {
      id
      name
    }
  }
`;

function CreateFundForm({ addFundToTable }) {
  const { register, handleSubmit } = useForm();
  const [addFund, { data }] = useMutation(CREATE_FUND_QUERY, {
    onCompleted: data => {
      addFundToTable(data);
    }
  });
  const onSumbit = data => {
    addFund({ variables: { input: { name: data.fund } } });
  };
  return (
    <div>
      <form onSubmit={handleSubmit(onSumbit)}>
        <label>Fund name:</label>
        <input name="fund" ref={register} />
        <button type="submit">Create</button>
      </form>
    </div>
  );
}

const MuseumFunds = () => {
  const [state, dispatch] = useContext(Context);

  let { loading } = useQuery(GET_FUNDS_QUERY, {
    onCompleted: data => {
      dispatch({ type: "SET_FUNDS", payload: data.museumFunds });
    }
  });

  const [deleteFund] = useMutation(DELETE_FUND_QUERY);

  const history = useHistory();

  const handleClick = id => {
    history.push("/museumFund/edit/" + id);
  };

  const onCreateFund = data => {
    dispatch({ type: "ADD_FUND", payload: data.createMuseumFund });
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
              deleteFund({ variables: { id: row.original.id } });
              dispatch({ type: "REMOVE_FUND", payload: row.original.id });
            }}
            value={"remove"}
          >
            {"remove"}
          </button>
        )
      }
    ],
    [state.funds]
  );

  if (loading) return "Loading...";

  return (
    <div>
      <CreateFundForm addFundToTable={onCreateFund} />
      <TableContainer columns={columns} data={state.funds} />
    </div>
  );
};

export default MuseumFunds;
