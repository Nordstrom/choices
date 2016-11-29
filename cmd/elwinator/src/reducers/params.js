const paramInitialState = {
  id: '',
  experiment: '',
  name: '',
  isWeighted: false,
  choices: [],
  weights: [],
}

const param = (state = paramInitialState, action) => {
  switch (action.type) {
  case 'PARAM_ADD':
    return { ...state, id: action.id, experiment: action.experiment, name: action.name };
  case 'PARAM_NAME':
    return { ...state, name: action.name };
  case 'PARAM_TOGGLE_WEIGHTED':
    return { ...state, isWeighted: !state.isWeighted };
  case 'PARAM_ADD_CHOICE':
    return { ...state, choices: [...state.choices, action.choice] };
  case 'PARAM_DELETE_CHOICE':
    return {
      ...state,
      choices: state.choices.filter((_, i) => i !== action.index),
      weights: state.weights.filter((_, i) => i !== action.index)
    };
  case 'PARAM_ADD_WEIGHT':
    return { ...state, weights: [...state.weights, action.weight] };
  case 'PARAM_CLEAR_CHOICES':
    return { ...state, choices: [], weights: [] };
  default:
    return state;
  }
};

/**
 * getParam returns the param object for the id supplied
 * @param {Object} state - the param state object.
 * @param {string} param - the param id.
 */
export const getParam = (state, param) => {
  return state.find(p => p.id === param);
};

/**
 * getParams returns the param objects for the id's supplied
 * @param {Object} state - the param state object.
 * @param {Array} paramIDs - an array of param id's.
 */
export const getParams = (state, paramIDs) => {
  return paramIDs.map(pid => state.find(p => pid === p.id));
};

const params = (state = [], action) => {
  switch (action.type) {
  case 'PARAM_ADD':
    return [...state, param(undefined, action)];
  case 'PARAM_DELETE':
    return state.filter(p => p.id !== action.param);
  case 'PARAM_NAME':
  case 'PARAM_TOGGLE_WEIGHTED':
  case 'PARAM_ADD_CHOICE':
  case 'PARAM_DELETE_CHOICE':
  case 'PARAM_ADD_WEIGHT':
  case 'PARAM_CLEAR_CHOICES':
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
