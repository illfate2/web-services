import React, { useState } from "react";
import { useQuery, gql, useMutation, useLazyQuery } from "@apollo/client";
import TableContainer from "./TableContainer";
import { useHistory } from "react-router-dom";
import { useForm } from "react-hook-form";

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
  const [museumFundsData, setMuseumFundsData] = useState([]);
  let { loading } = useQuery(GET_FUNDS_QUERY, {
    onCompleted: data => {
      setMuseumFundsData(data.museumFunds);
    }
  });

  const [deleteFund, { deleteData }] = useMutation(DELETE_FUND_QUERY);

  const history = useHistory();

  const handleClick = id => {
    history.push("/museumFund/edit/" + id);
  };

  const onCreateFund = data => {
    let dataCopy = [...museumFundsData];
    dataCopy.push({
      id: data.createMuseumFund.id,
      name: data.createMuseumFund.name
    });
    setMuseumFundsData(dataCopy);
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
              const dataCopy = [...museumFundsData];
              dataCopy.splice(row.index, 1);
              setMuseumFundsData(dataCopy);
            }}
            value={"remove"}
          >
            {"remove"}
          </button>
        )
      }
    ],
    [museumFundsData]
  );

  if (loading) return "Loading...";

  return (
    <div>
      <CreateFundForm addFundToTable={onCreateFund} />
      <TableContainer columns={columns} data={museumFundsData} />
    </div>
  );
};

export default MuseumFunds;
