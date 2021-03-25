import React from "react";
import { useQuery, gql, useMutation, useLazyQuery } from "@apollo/client";
import { useParams } from "react-router-dom";

const GET_ITEMS_QUERY = gql`
  query MuseumItem($id: ID!) {
    museumItem(id: $id) {
      id
      inventoryNumber
      name
      annotation
      creationDate
      person {
        firstName
        lastName
        middleName
      }
    }
  }
`;

export const ViewMuseumItem = () => {
  const id = useParams().id;
  const { loading, data } = useQuery(GET_ITEMS_QUERY, {
    variables: {
      id: id
    }
  });
  if (loading) return "loading...";
  return (
    <div>
      <h1>Item ID: {data.museumItem.id}</h1>
      <h2>Name: {data.museumItem.name}</h2>
      <h2>Inventory number: {data.museumItem.inventoryNumber}</h2>
      <h2>Creation Date: {data.museumItem.creationDate}</h2>
      <h2>Annotation: {data.museumItem.annotation}</h2>
      <h3>First Name: {data.museumItem.person.firstName}</h3>
      <h3>Last Name: {data.museumItem.person.lastName}</h3>
      <h3>Middle Name: {data.museumItem.person.middleName}</h3>
    </div>
  );
};
