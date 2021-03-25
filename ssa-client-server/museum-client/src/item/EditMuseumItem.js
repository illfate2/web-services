import { gql, useMutation, useQuery } from "@apollo/client";
import { useForm } from "react-hook-form";
import { useParams, useHistory } from "react-router-dom";
import React, { useState, useContext } from "react";
import { Context } from "../store/Store";

const UPDATE_ITEM_QUERY = gql`
  mutation UpdateMuseumItem($id: ID!, $input: UpdateMuseumItemInput!) {
    updateMuseumItem(id: $id, input: $input) {
      id
    }
  }
`;

const GET_ITEM_QUERY = gql`
  query MuseumItem($id: ID!) {
    museumItem(id: $id) {
      id
      inventoryNumber
      name
      annotation
      creationDate
      set {
        id
        name
      }
      fund {
        id
        name
      }
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

export const EditMuseumItem = () => {
  const [state, dispatch] = useContext(Context);
  const id = useParams().id;
  const [updateItem] = useMutation(UPDATE_ITEM_QUERY, {
    onCompleted: data => {
      dispatch({ type: "UPDATE_ITEM", payload: data.updateMuseumItem });
    }
  });
  const [item, setItem] = useState();
  const { register, handleSubmit } = useForm();
  const history = useHistory();
  const onSubmit = data => {
    updateItem({
      variables: {
        id: id,
        input: {
          name: data.name,
          inventoryNumber: data.inventory_number,
          annotation: data.annotation,
          creationDate: data.creation_date + ":00Z",
          setID: data.sets,
          fundID: data.funds
        }
      }
    });
    history.push("/museumItems");
  };
  const { loading: funds_loading, data: funds_data } = useQuery(
    GET_FUNDS_QUERY
  );
  const { loading: sets_loading, data: sets_data } = useQuery(GET_SETS_QUERY);

  const { loading: item_loading, data: item_data } = useQuery(GET_ITEM_QUERY, {
    variables: {
      id: id
    },
    onCompleted: () => {
      setItem(item_data.museumItem);
    }
  });

  if (funds_loading) return "loadding";
  if (sets_loading) return "loadding";
  if (item_loading) return "loading";
  return (
    <div>
      <form onSubmit={handleSubmit(onSubmit)}>
        <input
          name="name"
          value={item && item.name}
          ref={register}
          onChange={e => {
            setItem({
              ...item,
              name: e.target.value
            });
          }}
        />
        <input
          name="inventory_number"
          ref={register}
          value={item && item.inventoryNumber}
          placeholder="Inventory number"
          onChange={e => {
            setItem({
              ...item,
              inventoryNumber: e.target.value
            });
          }}
        />
        <input
          name="creation_date"
          required
          type="datetime-local"
          ref={register}
          value={item && item.creationDate}
          placeholder="Creation date"
          onChange={e => {
            setItem({
              ...item,
              creationDate: e.target.value
            });
          }}
        />
        <input
          name="annotation"
          ref={register}
          placeholder="Annotation"
          value={item && item.annotation}
          onChange={e => {
            setItem({
              ...item,
              annotation: e.target.value
            });
          }}
        />
        <br></br>
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
        <button type="submit">Update</button>
      </form>
    </div>
  );
};
