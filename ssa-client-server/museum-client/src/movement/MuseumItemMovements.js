import React, {useContext, useState} from "react";
import {useQuery, gql, useMutation} from "@apollo/client";
import TableContainer from "../table/TableContainer";
import {useHistory} from "react-router-dom";
import {useForm} from "react-hook-form";
import {Context} from "../store/Store";


const CREATE_MOVEMENT_QUERY = gql`
    mutation CreateMuseumItemMovement($input: MuseumMovementInput!) {
        createMuseumItemMovement(input: $input) {
            id
            acceptDate
            exhibitTransferDate
            exhibitReturnDate
        }
    }
`;

const DELETE_MOVEMENT_MUTATION = gql`
    mutation DeleteMuseumItemMovement($id: ID!) {
        deleteMuseumItemMovement(id: $id)
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

const GET_MOVEMENTS_QUERY = gql`
    {
        museumMovements {
            id
            acceptDate
            exhibitTransferDate
            exhibitReturnDate
        }
    }
`;

function CreateMovementForm({addMovememtToTable}) {
    const [addMovement] = useMutation(CREATE_MOVEMENT_QUERY, {
        onCompleted: data => {
            addMovememtToTable(data);
        }
    });

    const {loading, data} = useQuery(GET_ITEMS_QUERY);
    const {register, handleSubmit} = useForm();
    const onSubmit = data => {
        addMovement({
            variables: {
                input: {
                    acceptDate: data.creation_date + ":00Z",
                    exhibitTransferDate: data.transfer_date + ":00Z",
                    exhibitReturnDate: data.return_date + ":00Z",
                    itemID: data.item_id,
                    personInput: {
                        firstName: data.first_name,
                        lastName: data.last_name,
                        middleName: data.middle_name
                    }
                }
            }
        });
    };
    if (loading) return "loading...";
    return (
        <div>
            <form onSubmit={handleSubmit(onSubmit)}>
                <input
                    name="creation_date"
                    required
                    type="datetime-local"
                    ref={register}
                    placeholder="Accept date"
                />
                <input
                    name="transfer_date"
                    required
                    type="datetime-local"
                    ref={register}
                    placeholder="Exhibit transfer date"
                />
                <input
                    name="return_date"
                    required
                    type="datetime-local"
                    ref={register}
                    placeholder="Exhibit return date"
                />
                <br></br>
                <input name="first_name" ref={register} placeholder="First name"/>
                <input name="last_name" ref={register} placeholder="Last name"/>
                <input name="middle_name" ref={register} placeholder="Middle name"/>
                <select name="item_id" ref={register}>
                    {data.museumItems.map(i => (
                        <option key={i.id} value={i.id}>
                            {i.name}
                        </option>
                    ))}
                </select>
                <button type="submit">Create</button>
            </form>
        </div>
    );
}

const MuseumItemMovements = () => {
    const [state, dispatch] = useContext(Context);

    let {loading} = useQuery(GET_MOVEMENTS_QUERY, {
        onCompleted: data => {
            dispatch({ type: "SET_MOVEMENTS", payload: data.museumMovements });
        }
    });

    const [deleteItem] = useMutation(DELETE_MOVEMENT_MUTATION);

    const history = useHistory();

    const handleClick = (id, path) => {
        history.push("/museumItemMovement/" + path + "/" + id);
    };

    const onCreateMovement = data => {
        dispatch({ type: "ADD_MOVEMENT", payload: data.createMuseumItemMovement });
    };

    const columns = React.useMemo(
        () => [
            {
                Header: "id",
                accessor: "id",
                show: false
            },
            {
                Header: "Accept Date",
                accessor: "acceptDate"
            },
            {
                Header: "Exhibit Transfer Date",
                accessor: "exhibitTransferDate"
            },
            {
                Header: "view",
                accessor: "view",
                Cell: ({row}) => (
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
                Cell: ({row}) => (
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
                Cell: ({row}) => (
                    <button
                        onClick={() => {
                            deleteItem({variables: {id: row.original.id}});
                            dispatch({type:'REMOVE_MOVEMENT',payload:row.original.id})
                        }}
                        value={"remove"}
                    >
                        {"remove"}
                    </button>
                )
            }
        ],
        [state.movements]
    );

    if (loading) return "Loading...";

    return (
        <div>
            <CreateMovementForm addMovememtToTable={onCreateMovement}/>
            <TableContainer columns={columns} data={state.movements}/>
        </div>
    );
};

export default MuseumItemMovements;
