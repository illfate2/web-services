import { gql, useMutation, useQuery } from "@apollo/client";
import { useForm } from "react-hook-form";
import { useParams, useHistory } from "react-router-dom";
import React, { useState } from "react";

const UPDATE_MOVEMENT_QUERY = gql`
  mutation UpdateMuseumItemMovement($id: ID!, $input: UpdateMuseumMovementInput!) {
    updateMuseumItemMovement(id: $id, input: $input) {
      id
    }
  }
`;

const EditMuseumItemMovement = () => {
  const id = useParams().id;
  const [updateItem] = useMutation(UPDATE_MOVEMENT_QUERY);
  const [item, setItem] = useState();
  const { register, handleSubmit } = useForm();
  const history = useHistory();
  const onSubmit = data => {
    updateItem({
      variables: {
        id: id,
        input: {
          acceptDate: data.accept_date + ":00Z",
          exhibitTransferDate: data.exhibit_transfer_date + ":00Z",
          exhibitReturnDate: data.exhibit_return_date + ":00Z"
        }
      }
    });
    history.push("/museumItemMovements");
  };

  return (
    <div>
      <form onSubmit={handleSubmit(onSubmit)}>
        <label>Accept date</label>
        <input
          name="accept_date"
          required
          type="datetime-local"
          ref={register}
          value={item && item.acceptDate}
          placeholder="Accept date"
          onChange={e => {
            setItem({
              ...item,
              acceptDate: e.target.value
            });
          }}
        />
        <br></br>
        <label>Exhibit transfer date</label>
        <input
          name="exhibit_transfer_date"
          required
          type="datetime-local"
          ref={register}
          value={item && item.transferDate}
          placeholder="Exhibit transfer date"
          onChange={e => {
            setItem({
              ...item,
              transferDate: e.target.value
            });
          }}
        />
        <br></br>
        <label>Exhibit return date</label>
        <input
          name="exhibit_return_date"
          required
          type="datetime-local"
          ref={register}
          value={item && item.creationDate}
          placeholder="Exhibit return date"
          onChange={e => {
            setItem({
              ...item,
              creationDate: e.target.value
            });
          }}
        />
        <br></br>
        <button type="submit">Update</button>
      </form>
    </div>
  );
};

export default EditMuseumItemMovement;
