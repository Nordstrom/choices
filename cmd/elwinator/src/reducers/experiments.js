import params from './params';
import shuffle from '../shuffle';

const experimentInitialState = {
  name: '',
  segments: new Array(128).fill(1),
  numSegments: 128,
  params: [],
  dirtySegments: true,
};

const experiment = (state = experimentInitialState, action) => {
  switch(action.type) {
  case 'ADD_EXPERIMENT':
    return { ...state, name: action.name };
  case 'EXPERIMENT_NAME':
    return { ...state, name: action.name };
  case 'EXPERIMENT_NUM_SEGMENTS':
    const ns = parseInt(action.numSegments, 10);
    const ensSegments = new Array(128).fill(0).fill(1, 0, ns);
    return { ...state, numSegments: ns, segments: shuffle(ensSegments) };
  case 'EXPERIMENT_PERCENT':
    const p = Math.floor((parseFloat(action.percent) / 100)*128);
    const epSegments = new Array(128).fill(0).fill(1, 0, p);
    return { ...state, numSegments: p, segments: shuffle(epSegments) };
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
  case 'EXPERIMENT_NUM_SEGMENTS':
  case 'EXPERIMENT_PERCENT':
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
