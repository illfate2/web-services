import React from "react";
import { useQuery, gql, useMutation, useLazyQuery } from "@apollo/client";
import TableContainer from "./TableContainer";
import "./popup.css";
import { useHistory } from "react-router-dom";
import { useForm } from "react-hook-form";

const CREATE_ITEM_QUERY = gql`
  mutation CreateMuseumItem($input: MuseumItemInput!) {
    createMuseumItem(input: $input) {
      id
    }
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

const GET_FUNDS_QUERY = gql`
  {
    museumFunds {
      id
      name
    }
  }
`;

const DELETE_ITEM_QUERY = gql`
  mutation DeleteMuseumItem($id: ID!) {
    deleteMuseumItem(id: $id)
  }
`;

const GET_ITEMS_QUERY = gql`
  {
    museumItems {
      id
      inventoryNumber
      name
    }
  }
`;

function CreateItemForm() {
  const [addItem] = useMutation(CREATE_ITEM_QUERY);
  const { register, handleSubmit } = useForm();
  const onSubmit = data => {
    addItem({
      variables: {
        input: {
          name: data.name,
          inventoryNumber: data.inventory_number,
          annotation: data.annotation,
          creationDate: data.creation_date + ":00Z",
          setID: data.sets,
          fundID: data.funds,
          personInput: {
            firstName: data.first_name,
            lastName: data.last_name,
            middleName: data.middle_name
          }
        }
      }
    });
    window.location.reload(false);
  };
  const { loading: funds_loading, data: funds_data } = useQuery(
    GET_FUNDS_QUERY
  );
  const { loading: sets_loading, data: sets_data } = useQuery(GET_SETS_QUERY);

  if (funds_loading) return "loadding";
  if (sets_loading) return "loadding";
  return (
    <div>
      <form onSubmit={handleSubmit(onSubmit)}>
        <input name="name" ref={register} placeholder="Name" />
        <input
          name="inventory_number"
          ref={register}
          placeholder="Inventory number"
        />
        <input
          name="creation_date"
          required
          type="datetime-local"
          ref={register}
          placeholder="Creation date"
        />
        <input name="annotation" ref={register} placeholder="Annotation" />
        <br></br>
        <input name="first_name" ref={register} placeholder="First name" />
        <input name="last_name" ref={register} placeholder="Last name" />
        <input name="middle_name" ref={register} placeholder="Middle name" />
        <select name="funds" ref={register}>
          {funds_data.museumFunds.map(f => (
            <option key={f.id} value={f.id}>
              {f.name}
            </option>
          ))}
        </select>
        <select name="sets" ref={register}>
          {sets_data.museumSets.map(s => (
            <option key={s.id} value={s.id}>
              {s.name}
            </option>
          ))}
        </select>
        <button type="submit">Create</button>
      </form>
    </div>
  );
}

const MuseumItems = () => {
  let { loading, error, data } = useQuery(GET_ITEMS_QUERY);

  const [deleteItem, { deleteData }] = useMutation(DELETE_ITEM_QUERY);

  const history = useHistory();

  const handleClick = (id, path) => {
    history.push("/museumItem/" + path + "/" + id);
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
        Header: "Inventory Number",
        accessor: "inventoryNumber"
      },
      {
        Header: "view",
        accessor: "view",
        Cell: ({ row }) => (
          <button
            onClick={() => {
                handleClick(row.original.id, "view");
            }}
            value={"view"}
          >
            {"view"}
          </button>
        )
      },
      {
        Header: "edit",
        accessor: "edit",
        Cell: ({ row }) => (
          <button
            onClick={() => {
              handleClick(row.original.id, "edit");
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
              deleteItem({ variables: { id: row.original.id } });
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
      <CreateItemForm />
      <TableContainer columns={columns} data={data.museumItems} />
    </div>
  );
};

export default MuseumItems;
