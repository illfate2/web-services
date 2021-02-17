import React from "react";
import { useQuery, gql, useMutation, useLazyQuery } from "@apollo/client";
import TableContainer from "./TableContainer";
import { useHistory } from "react-router-dom";

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

function CreateFundForm() {
  let input;
  const [addFund, { data }] = useMutation(CREATE_FUND_QUERY);
  return (
    <div>
      <form
        onSubmit={e => {
          addFund({ variables: { input: { name: input.value } } });
          input.value = "";
        }}
      >
        <label>Fund name:</label>
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

const MuseumFunds = () => {
  let { loading, error, data } = useQuery(GET_FUNDS_QUERY);

  const [deleteFund, { deleteData }] = useMutation(DELETE_FUND_QUERY);

  const history = useHistory();

  const handleClick = id => {
    history.push("/museumFund/edit/" + id);
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
      <CreateFundForm />
      <TableContainer columns={columns} data={data.museumFunds} />
    </div>
  );
};

export default MuseumFunds;