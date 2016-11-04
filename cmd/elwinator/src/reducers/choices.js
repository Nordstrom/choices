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
      choices: [...state.choices, state.choice],
      weights: [...state.weights, state.weight],
    });
  default:
    return state;
  }
}

export default choices;
