import labels from './label';
import experiments from './experiments';

const namespaceInitialState = {
  name: '',
  labels: [],
  experiments: [],
};

const namespace = (state = namespaceInitialState, action) => {
  switch (action.type) {
  case 'NAMESPACE_NAME':
    return { ...state, name: 'action.name'};
  case 'ADD_LABEL':
  case 'TOGGLE_LABEL':
    return { ...state, labels: labels(state.labels, action)}
  case 'ADD_EXPERIMENT':
  case 'EXPERIMENT_NAME':
  case 'PARAM_NAME':
  case 'ADD_PARAM':
    return { ...state, experiments: experiments(state.experiments, action)};
  default: 
    return state;
  }
}

export default namespace;
