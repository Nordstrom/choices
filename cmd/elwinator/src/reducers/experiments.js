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
  case 'TOGGLE_WEIGHTED':
  case 'ADD_CHOICE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
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
  case 'TOGGLE_WEIGHTED':
  case 'ADD_CHOICE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
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
