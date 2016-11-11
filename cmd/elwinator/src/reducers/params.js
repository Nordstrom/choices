const paramInitialState = {
  name: '',
  isWeighted: false,
  choices: [],
  weights: [],
}

const param = (state = paramInitialState, action) => {
  switch (action.type) {
  case 'ADD_PARAM':
    return {...state, ...action};
  case 'PARAM_NAME':
    return { ...state, name: action.name };
  default:
    return state;
  }
}

const params = (state = [], action) => {
  switch (action.type) {
  case 'ADD_PARAM':
    return [...state, param(undefined, action)];
  case 'PARAM_NAME':
    const pars = state.params.map(p => {
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
