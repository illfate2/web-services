import React, { useState } from "react";
import { useQuery, gql, useMutation } from "@apollo/client";
import TableContainer from "./TableContainer";
import "./popup.css";
import { useHistory } from "react-router-dom";
import { useForm } from "react-hook-form";

const CREATE_ITEM_QUERY = gql`
  mutation CreateMuseumItem($input: MuseumItemInput!) {
    createMuseumItem(input: $input) {
      id
      inventoryNumber
      name
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

function CreateItemForm({ onSubmit }) {
  const { register, handleSubmit } = useForm();
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
  const [museumItemsData, setMuseumItemsData] = useState([]);
  const { loading } = useQuery(GET_ITEMS_QUERY, {
    onCompleted: data => {
      setMuseumItemsData(data.museumItems);
    }
  });

  const [addItem] = useMutation(CREATE_ITEM_QUERY, {
    onCompleted: data => {
      let dataCopy = [...museumItemsData];
      dataCopy.push({
        id: data.createMuseumItem.id,
        name: data.createMuseumItem.name,
        inventoryNumber: data.createMuseumItem.inventoryNumber
      });
      setMuseumItemsData(dataCopy);
    }
  });
  const onAddItemSubmit = data => {
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
  };

  const [deleteItem] = useMutation(DELETE_ITEM_QUERY);

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
        Cell: ({ row }) => {
          return (
            <div>
              <button
                onClick={() => {
                  deleteItem({ variables: { id: row.original.id } });
                  const dataCopy = [...museumItemsData];
                  dataCopy.splice(row.index, 1);
                  setMuseumItemsData(dataCopy);
                }}
              >
                delete
              </button>
            </div>
          );
        }
      }
    ],
    [museumItemsData]
  );

  if (loading) return "Loading...";

  return (
    <div>
      <CreateItemForm onSubmit={onAddItemSubmit} />
      <TableContainer columns={columns} data={museumItemsData} />
    </div>
  );
};

export default MuseumItems;
