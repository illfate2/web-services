import React, {useContext} from "react";
import {gql} from "@apollo/client";
import TableContainer from "../table/TableContainer";
import {useHistory} from "react-router-dom";
import {useForm} from "react-hook-form";
import {Context} from "../store/Store";
import {useMutationWithAuthErrHandling, useQueryWithAuthErrHandling} from "../Queries";

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

function CreateItemForm({onSubmit}) {
    const [state, dispatch] = useContext(Context);

    const {register, handleSubmit} = useForm();
    const {loading: funds_loading} = useQueryWithAuthErrHandling(GET_FUNDS_QUERY, {
        onCompleted: data => {
            dispatch({type: "SET_FUNDS", payload: data.museumFunds});
        }
    });
    const {loading: sets_loading} = useQueryWithAuthErrHandling(GET_SETS_QUERY, {
        onCompleted: data => {
            dispatch({type: "SET_SETS", payload: data.museumSets});
        }
    });

    if (funds_loading) return "loadding";
    if (sets_loading) return "loadding";
    return (
        <div>
            <form onSubmit={handleSubmit(onSubmit)}>
                <input name="name" ref={register} placeholder="Name"/>
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
                <input name="annotation" ref={register} placeholder="Annotation"/>
                <br></br>
                <input name="first_name" ref={register} placeholder="First name"/>
                <input name="last_name" ref={register} placeholder="Last name"/>
                <input name="middle_name" ref={register} placeholder="Middle name"/>
                <select name="funds" ref={register}>
                    {state.funds.map(f => (
                        <option key={f.id} value={f.id}>
                            {f.name}
                        </option>
                    ))}
                </select>
                <select name="sets" ref={register}>
                    {state.sets.map(s => (
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
    const [state, dispatch] = useContext(Context);

    const {loading} = useQueryWithAuthErrHandling(GET_ITEMS_QUERY, {
        onCompleted: data => {
            dispatch({type: "SET_ITEMS", payload: data.museumItems});
        }
    });

    const [addItem] = useMutationWithAuthErrHandling(CREATE_ITEM_QUERY, {
        onCompleted: data => {
            dispatch({type: "ADD_ITEM", payload: data.createMuseumItem});
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

    const [deleteItem] = useMutationWithAuthErrHandling(DELETE_ITEM_QUERY);

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
                Cell: ({row}) => {
                    return (
                        <div>
                            <button
                                onClick={() => {
                                    deleteItem({variables: {id: row.original.id}});
                                    dispatch({type: "REMOVE_ITEM", payload: row.original.id});
                                }}
                            >
                                delete
                            </button>
                        </div>
                    );
                }
            }
        ],
        [state.items]
    );

    if (loading) return "Loading...";

    return (
        <div>
            <CreateItemForm onSubmit={onAddItemSubmit}/>
            <TableContainer columns={columns} data={state.items}/>
        </div>
    );
};

export default MuseumItems;
