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

export const getLabels = (state, labelIDs) => {
  return labelIDs.map(lid => state.find(l => lid === l.id));
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
