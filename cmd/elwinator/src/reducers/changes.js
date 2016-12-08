const changes = (state = {}, action) => {
  switch (action.type) {
  case 'CHANGES_CLEAR':
    let s = {};
    Object.keys(state)
    .filter(k => !action.changes.find(c => c === k))
    .forEach(k => s = Object.assign({}, s, state[k]));
    return s;
  case 'NAMESPACES_LOADED':
  case 'NAMESPACE_ADD':
  case 'NAMESPACE_DELETE':
  case 'NAMESPACE_NAME':
  case 'NAMESPACE_ADD_LABEL':
  case 'NAMESPACE_TOGGLE_LABEL':
  case 'EXPERIMENT_ADD':
  case 'EXPERIMENT_DELETE':
  case 'EXPERIMENT_NAME':
  case 'EXPERIMENT_NUM_SEGMENTS':
  case 'PARAM_ADD':
  case 'PARAM_DELETE':
  case 'PARAM_NAME':
  case 'PARAM_TOGGLE_WEIGHTED':
  case 'PARAM_ADD_CHOICE':
  case 'PARAM_DELETE_CHOICE':
  case 'PARAM_ADD_WEIGHT':
  case 'PARAM_CLEAR_CHOICES':
    return {
      ...state,
      [action.namespace]: state[action.namespace] ? [...state[action.namespace], action] : [action],
    };
  default:
    return state;
  }
};

export default changes;