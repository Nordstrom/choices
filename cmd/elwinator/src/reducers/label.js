const initialLabelState = {
  id: 0,
  name: '',
  active: false,
};

const label = (state = initialLabelState, action) => {
  switch (action.type) {
  case 'ADD_LABEL':
    return {
      id: action.id,
      name: action.name,
      active: true,
    };
  case 'TOGGLE_LABEL':
    return {...state,active: !state.active};
  default:
    return state;
  }
}

export default label;
