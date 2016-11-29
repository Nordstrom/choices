const paramInitialState = {
  id: '',
  name: '',
  isWeighted: false,
  choices: [],
  weights: [],
}

const param = (state = paramInitialState, action) => {
  switch (action.type) {
  case 'ADD_PARAM':
    return { ...state, id: action.id, name: action.name };
  case 'PARAM_NAME':
    return { ...state, name: action.name };
  case 'TOGGLE_WEIGHTED':
    return { ...state, isWeighted: !state.isWeighted };
  case 'ADD_CHOICE':
    return { ...state, choices: [...state.choices, action.choice] };
  case 'CHOICE_DELETE':
    return {
      ...state,
      choices: state.choices.filter((_, i) => i !== action.index),
      weights: state.weights.filter((_, i) => i !== action.index)
    };
  case 'ADD_WEIGHT':
    return { ...state, weights: [...state.weights, action.weight] };
  case 'CLEAR_CHOICES':
    return { ...state, choices: [], weights: [] };
  default:
    return state;
  }
};

/**
 * getParams retuns the param objects for the id's supplied
 * @param {Object} state - the param state object.
 * @param {Array} paramIDs - an array of param id's.
 */
export const getParams = (state, paramIDs) => {
  return paramIDs.map(pid => state.find(p => pid === p.id));
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
  case 'CHOICE_DELETE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
    const pars = state.map(p => {
      if (p.id !== action.param) {
        return p;
      }
      return param(p, action);
    });
    return pars;
  default:
    return state;
  }
};

export default params;
