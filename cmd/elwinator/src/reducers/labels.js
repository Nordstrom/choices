const label = (state, action) => {
  switch (action.type) {
  case 'CREATE_LABEL':
    return {
      id: action.id,
      name: action.name,
      active: true,
    };
  default:
    return state;
  }
}

const labels = (state = [], action) => {
  switch (action.type) {
  case 'REMOVE_LABEL':
    return state.map(label => {
      if (label.id === action.id) {
        return { id: label.id, name: label.name, active: false };
      }
      return label;
    });
  case 'ADD_LABEL':
    return state.map(label => {
      if (label.id === action.id) {
        return { id: label.id, name: label.name, active: true };
      }
      return label;
    });
  case 'CREATE_LABEL':
    return [...state, label(undefined, action)];
  default:
    return state;
  }
}

export default labels;
