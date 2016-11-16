import labels from './label';
import experiments from './experiments';

const namespaceInitialState = {
  name: '',
  labels: [],
  experiments: [],
};

const namespace = (state = namespaceInitialState, action) => {
  switch (action.type) {
  case 'ADD_NAMESPACE':
  case 'NAMESPACE_NAME':
    return { ...state, name: action.name};
  case 'ADD_LABEL':
  case 'TOGGLE_LABEL':
    return { ...state, labels: labels(state.labels, action)}
  case 'ADD_EXPERIMENT':
  case 'EXPERIMENT_NAME':
  case 'EXPERIMENT_NUM_SEGMENTS':
  case 'PARAM_NAME':
  case 'ADD_PARAM':
  case 'TOGGLE_WEIGHTED':
  case 'ADD_CHOICE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
    return { ...state, experiments: experiments(state.experiments, action)};
  default: 
    return state;
  }
}

const namespaces = (state = [], action) => {
  switch (action.type) {
  case 'ADD_NAMESPACE':
    return [...state, namespace(undefined, action)];
  case 'NAMESPACE_NAME':
  case 'ADD_LABEL':
  case 'TOGGLE_LABEL':
  case 'ADD_EXPERIMENT':
  case 'EXPERIMENT_NAME':
  case 'EXPERIMENT_NUM_SEGMENTS':
  case 'PARAM_NAME':
  case 'ADD_PARAM':
  case 'TOGGLE_WEIGHTED':
  case 'ADD_CHOICE':
  case 'ADD_WEIGHT':
  case 'CLEAR_CHOICES':
    const ns = state.map(n => {
      if (n.name !== action.namespace) {
        return n;
      }
      return namespace(n, action);
    });
    return ns;
  default:
    return state;
  }
}

export default namespaces;
