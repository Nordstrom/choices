const initialLabelState = {
  name: '',
  active: false,
};

const label = (state = initialLabelState, action) => {
  switch (action.type) {
  case 'ADD_LABEL':
    return {
      name: action.name,
      active: true,
    };
  case 'TOGGLE_LABEL':
    return {...state, active: !state.active};
  default:
    return state;
  }
}

const labels = (state = [], action) => {
  switch (action.type) {
  case 'ADD_LABEL':
    return [...state, label(undefined, action)];
  case 'TOGGLE_LABEL':
    return state.map(l => {
      if (l.name !== action.name) {
        return l;
      }
      return label(l, action);
    });
  default:
    return state;
  }
}

export default labels;
