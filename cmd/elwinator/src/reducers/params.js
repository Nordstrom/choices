const param = (state, action) => {
  switch (action.type) {
    case 'ADD_PARAM':
      return {
        id: action.id,
        name: action.name,
      };
    default:
      return state;
  }
}

const params = (state, action) => {
  switch (action.type) {
    case 'ADD_PARAM':
      return [
        ...state,
        param(undefined, action)
      ];
    case 'REMOVE_PARAM':
      return state.filter(p => p.id !== action.id);
    default:
      return state;
  }
}

const initialState = {
  params: [],
  isWeighted: false,
}

const paramContainer = (state = initialState, action) => {
  switch (action.type) {
  case 'ADD_PARAM':
    return Object.assign({}, state, { params: params(state.params, action) });
  case 'REMOVE_PARAM':
    return Object.assign({}, state, { params: params(state.params, action) });
  case 'TOGGLE_WEIGHTED':
    return Object.assign({}, state, { isWeighted: !state.isWeighted });
  default:
    return state;
  }
}

export default paramContainer
