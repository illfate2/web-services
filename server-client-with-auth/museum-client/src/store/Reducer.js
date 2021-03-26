const Reducer = (state, action) => {
    switch (action.type) {
        case "SET_SETS":
            return {
                ...state,
                sets: action.payload
            };
        case "ADD_SET":
            let new_sets = [...state.sets];
            new_sets.push(action.payload);
            return {
                ...state,
                sets: new_sets
            };
        case "REMOVE_SET":
            return {
                ...state,
                sets: state.sets.filter(function (value, index, arr) {
                    return value.id != action.payload;
                })
            };
        case "UPDATE_SET":
            let updated_sets = [...state.sets];
            const idx = updated_sets.findIndex(obj => obj.id === action.payload.id);
            updated_sets[idx] = action.payload;
            return {
                ...state,
                sets: updated_sets
            };
        case "SET_FUNDS":
            return {
                ...state,
                funds: action.payload
            };
        case "ADD_FUND":
            let new_funds = [...state.funds];
            new_funds.push(action.payload);
            return {
                ...state,
                funds: new_funds
            };
        case "REMOVE_FUND":
            return {
                ...state,
                funds: state.funds.filter(function (value, index, arr) {
                    return value.id != action.payload;
                })
            };
        case "UPDATE_FUND":
            let updated_funds = [...state.funds];
            const fund_idx = updated_funds.findIndex(obj => obj.id == action.payload.id);
            updated_funds[fund_idx] = action.payload;
            return {
                ...state,
                funds: updated_funds
            };
        case "SET_ITEMS":
            return {
                ...state,
                items: action.payload
            };
        case "ADD_ITEM":
            let new_items = [...state.items];
            new_items.push(action.payload);
            return {
                ...state,
                items: new_items
            };
        case "REMOVE_ITEM":
            return {
                ...state,
                items: state.items.filter(function (value, index, arr) {
                    return value.id != action.payload;
                })
            };
        case "UPDATE_ITEM":
            let updated_items = [...state.items];
            const item_idx = updated_items.findIndex(obj => obj.id === action.payload.id);
            updated_items[item_idx] = action.payload;
            return {
                ...state,
                items: updated_items
            };
        case "SET_MOVEMENTS":
            return {
                ...state,
                movements: action.payload
            };
        case "ADD_MOVEMENT":
            let new_movements = [...state.movements];
            new_movements.push(action.payload);
            return {
                ...state,
                movements: new_movements
            };
        case "REMOVE_MOVEMENT":
            return {
                ...state,
                movements: state.movements.filter(function (value, index, arr) {
                    return value.id !== action.payload;
                })
            };
        case "UPDATE_MOVEMENT":
            let updated_movements = [...state.movements];
            const movement_idx = updated_movements.findIndex(obj => obj.id === action.payload.id);
            updated_movements[movement_idx] = action.payload;
            return {
                ...state,
                movements: updated_movements
            };
        default:
            return state;
    }
};

export default Reducer;
