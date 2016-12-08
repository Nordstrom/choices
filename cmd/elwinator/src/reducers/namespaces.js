const namespaceInitialState = {
  name: '',
  labels: [],
  experiments: [],
  isDirty: false,
  isNew: false,
  delete: false,
  publish: false,
};

const namespace = (state = namespaceInitialState, action) => {
  switch (action.type) {
  case 'NAMESPACE_ADD':
      return {
        ...state,
        name: action.namespace,
        isDirty: true,
        isNew: true,
      };
  case 'NAMESPACE_DELETE':
    return {
      ...state,
      delete: true,
      isDirty: true,
    };
  case 'NAMESPACE_NAME':
    return {
      ...state,
      name: action.name,
      isDirty: true,
    };
  case 'NAMESPACE_ADD_LABEL':
    return {
      ...state,
      labels: [...state.labels, action.id],
    };
  case 'NAMESPACE_TOGGLE_LABEL':
    if (state.labels.find(id => id === action.id)) {
      return {
        ...state,
        labels: state.labels.filter(n => n !== action.id),
        isDirty: true,
      };
    }
    return {
      ...state,
      labels: [...state.labels, action.id],
      isDirty: true,
    };
  case 'EXPERIMENT_ADD':
    return {
      ...state,
      experiments: [...state.experiments, action.id],
    };
  case 'EXPERIMENT_DELETE':
    return {
      ...state,
      experiments: state.experiments.filter(eid => eid !== action.experiment),
    }
  case 'NAMESPACE_TOGGLE_PUBLISH':
    return {
      ...state,
      publish: !state.publish,
    };
  default:
    return state;
  }
};

/**
 * getNamespace returns the namespace requested
 * @param {Object} state - the namespaces state Object
 * @param {string} name - the name of the namespace.
 */
export const getNamespace = (state, name) => {
  return state.find(n => n.name === name);
}

/**
 * getNamespaces returns the namespaces listed in names array
 * @param {Object} state - the namespaces state object.
 * @param {string} names - the names of the namespaces you want.
 */
export const getNamespaces = (state, names) => {
  return names.map(name => state.find(n => n.name === name));
}

const namespaces = (state = [], action) => {
  switch (action.type) {
  case 'NAMESPACES_LOADED':
    return action.namespaces.map(n => Object.assign({}, namespace(undefined, action), n));
  case 'NAMESPACE_ADD':
    return [...state, namespace(undefined, action)];
  case 'NAMESPACE_DELETE':
  case 'NAMESPACE_NAME':
  case 'NAMESPACE_ADD_LABEL':
  case 'NAMESPACE_TOGGLE_LABEL':
  case 'EXPERIMENT_ADD':
  case 'EXPERIMENT_DELETE':
  case 'NAMESPACE_TOGGLE_PUBLISH':
    const ns = state.map(n => {
      if (n.name !== action.namespace) {
        return n;
      }
      return namespace(n, action);
    });
    return ns;
  default:
    if (action.entities && action.entities.namespaces) {
      return action.entities.namespaces;
    }
    return state;
  }
};

export default namespaces;
