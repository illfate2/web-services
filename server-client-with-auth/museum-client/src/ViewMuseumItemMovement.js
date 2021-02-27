import React from "react";
import { useQuery, gql, useMutation, useLazyQuery } from "@apollo/client";
import { useParams } from "react-router-dom";

const GET_MOVEMENTS_QUERY = gql`
  query MuseumMovement($id: ID!) {
    museumMovement(id: $id) {
      id
      person {
        firstName
        lastName
        middleName
      }
    }
  }
`;

const ViewMuseumItemMovement = () => {
    const id = useParams().id;
    const { loading, data } = useQuery(GET_MOVEMENTS_QUERY, {
      variables: {
        id: id
      }
    });
    if (loading) return "loading...";
    return (
      <div>
        <h1>Movement ID: {data.museumMovement.id}</h1>
        {/* <h2>Name: {data.museumItem.name}</h2> */}
        {/* <h2>Inventory number: {data.museumItem.inventoryNumber}</h2> */}
        {/* <h2>Creation Date: {data.museumItem.creationDate}</h2> */}
        {/* <h2>Annotation: {data.museumItem.annotation}</h2> */}
        <h3>First Name: {data.museumMovement.person.firstName}</h3>
        <h3>Last Name: {data.museumMovement.person.lastName}</h3>
        <h3>Middle Name: {data.museumMovement.person.middleName}</h3>
      </div>
    );
};

export default ViewMuseumItemMovement;