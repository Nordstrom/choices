import params from './params';
import { sample } from '../shuffle';

const experimentInitialState = {
  name: '',
  segments: [],
  numSegments: 0,
  params: [],
};

const experiment = (state = experimentInitialState, action) => {
  switch(action.type) {
  case 'ADD_EXPERIMENT':
    return { ...state, name: action.name };
  case 'EXPERIMENT_NAME':
    return { ...state, name: action.name };
  case 'EXPERIMENT_NUM_SEGMENTS':
    const ns = parseInt(action.numSegments, 10);
    return { ...state, numSegments: ns, segments: sample(action.namespaceSegments, action.numSegments) };
  case 'PARAM_NAME':
  case 'ADD_PARAM':
  case 'PARAM_DELETE':
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
  case 'EXPERIMENT_DELETE':
    return state.filter(e => e.name !== action.experiment);
  case 'EXPERIMENT_NAME':
  case 'EXPERIMENT_NUM_SEGMENTS':
  case 'PARAM_NAME':
  case 'ADD_PARAM':
  case 'PARAM_DELETE':
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
