import React from "react";
import {gql,} from "@apollo/client";
import {useParams} from "react-router-dom";
import {useQueryWithAuthErrHandling} from "../Queries";

const GET_MOVEMENTS_QUERY = gql`
    query MuseumMovement($id: ID!) {
        museumMovement(id: $id) {
            id
            acceptDate
            exhibitTransferDate
            exhibitReturnDate
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
    const {loading, data} = useQueryWithAuthErrHandling(GET_MOVEMENTS_QUERY, {
        variables: {
            id: id
        }
    });
    if (loading) return "loading...";
    return (
        <div>
            <h1>Movement ID: {data.museumMovement.id}</h1>
            <h2>Accept Date: {data.museumMovement.acceptDate}</h2>
            <h2>Exhibit Transfer Date: {data.museumMovement.exhibitTransferDate}</h2>
            <h2>Exhibit Return Date: {data.museumMovement.exhibitReturnDate}</h2>
            <h3>First Name: {data.museumMovement.person.firstName}</h3>
            <h3>Last Name: {data.museumMovement.person.lastName}</h3>
            <h3>Middle Name: {data.museumMovement.person.middleName}</h3>
        </div>
    );
};

export default ViewMuseumItemMovement;
