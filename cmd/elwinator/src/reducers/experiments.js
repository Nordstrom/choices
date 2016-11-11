import params from './params';

const experimentInitialState = {
  name: '',
  segments: '',
  params: [],
};

const experiment = (state = experimentInitialState, action) => {
  switch(action.type) {
  case 'ADD_EXPERIMENT':
    return { ...state, name: action.name };
  case 'EXPERIMENT_NAME':
    return { ...state, name: action.name };
  case 'PARAM_NAME':
  case 'ADD_PARAM':
    return { ...state, params: params(state.params, action) };
  default:
    return state;
  }
}

const experiments = (state = [], action) => {
  switch (action.type) {
  case 'ADD_EXPERIMENT':
    return [...state, experiment(undefined, action)];
  case 'EXPERIMENT_NAME':
  case 'PARAM_NAME':
  case 'ADD_PARAM':
    const exps = state.map(e => {
      if (e.name !== action.experiment) {
        return e;
      }
      return experiment(e, action);
    });
    return exps;
  default:
    return state;
  }
}

export default experiments;
