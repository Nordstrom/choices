const initialState = {
  choices: [],
  weights: [],
  choice: '',
  weight: 1,
  isWeighted: false,
}

const choices = (state = initialState, action) => {
  switch (action.type) {
  case 'TOGGLE_WEIGHTED':
    // if we are toggling from uniform to weighted we need to reset the
    // choices since the saved choices will have incorrect weights.
    if (!state.isWeighted) {
      return { choices: [], weights: [], isWeighted: true};
    }
    return Object.assign({}, state, { isWeighted: !state.isWeighted });
  case 'UPDATE_CHOICE':
    return Object.assign({}, state, { choice: action.choice });
  case 'UPDATE_WEIGHT':
    return Object.assign({}, state, { weight: action.weight });
  case 'ADD_CHOICE':
    return Object.assign({}, state, {
      choice: '',
      choices: [...state.choices, action.choice],
      weights: [...state.weights, action.weight],
    });
  default:
    return state;
  }
}

export default choices;
