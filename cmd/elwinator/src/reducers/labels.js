const initialLabelState = {
  id: '',
  name: '',
};

const label = (state = initialLabelState, action) => {
  switch (action.type) {
  case 'ADD_LABEL':
    return {
      id: action.id,
      name: action.name,
    };
  default:
    return state;
  }
}

const labels = (state = [], action) => {
  switch (action.type) {
  case 'ADD_LABEL':
    return [...state, label(undefined, action)];
  default:
    return state;
  }
}

export default labels;
