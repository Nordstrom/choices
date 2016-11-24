const paramInitialState = {
  name: '',
  isWeighted: false,
  choices: [],
  weights: [],
}

const param = (state = paramInitialState, action) => {
  switch (action.type) {
  case 'ADD_PARAM':
    return { ...state, name: action.name };
  case 'PARAM_NAME':
    return { ...state, name: action.name };
  case 'TOGGLE_WEIGHTED':
    return { ...state, isWeighted: !state.isWeighted };
  case 'ADD_CHOICE':
    return { ...state, choices: [...state.choices, action.choice] };
  case 'ADD_WEIGHT':
    return { ...state, weights: [...state.weights, action.weight] };
  case 'CLEAR_CHOICES':
    return { ...state, choices: [], weights: [] };
  default:
    return state;
  }
}

const params = (state = [], action) => {
  switch (action.type) {
  case 'ADD_PARAM':
    return [...state, param(undefined, action)];
  case 'PARAM_DELETE':
    return state.filter(p => p.name !== action.name);
  case 'PARAM_NAME':
  case 'TOGGLE_WEIGHTED':
  case 'ADD_CHOICE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
    const pars = state.map(p => {
      if (p.name !== action.param) {
        return p;
      }
      return param(p, action);
    });
    return pars;
  default:
    return state;
  }
}

export default params;
