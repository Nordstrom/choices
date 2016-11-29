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
  case 'ADD_NAMESPACE':
      return {
        ...state,
        name: action.name,
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
  case 'ADD_LABEL':
    return {
      ...state,
      labels: [...state.labels, action.id],
    };
  case 'TOGGLE_LABEL':
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
  case 'ADD_EXPERIMENT':
    return {
      ...state,
      experiments: [...state.experiments, action.id],
    };
  case 'TOGGLE_PUBLISH':
    return {
      ...state,
      publish: !state.publish,
    };
  default:
    return state;
  }
};

const namespaces = (state = [], action) => {
  switch (action.type) {
  case 'NAMESPACES_LOADED':
    return action.namespaces.map(n => Object.assign({}, namespace(undefined, action), n));
  case 'ADD_NAMESPACE':
    return [...state, namespace(undefined, action)];
  case 'NAMESPACE_DELETE':
  case 'NAMESPACE_NAME':
  case 'ADD_LABEL':
  case 'TOGGLE_LABEL':
  case 'ADD_EXPERIMENT':
  case 'TOGGLE_PUBLISH':
    const ns = state.map(n => {
      if (n.name !== action.namespace) {
        return n;
      }
      return namespace(n, action);
    });
    return ns;
  default:
    if (action.entities && action.entities.namespaces) {
      return Object.assign({}, state, action.entities.namespaces);
    }
    return state;
  }
};

export default namespaces;
