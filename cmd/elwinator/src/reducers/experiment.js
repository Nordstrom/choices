const expInitialState = {
  name: '',
  params: [],
};

const experiment = (state = expInitialState, action) => {
  switch (action.type) {
  case 'UPDATE_NAME':
    return Object.assign({}, state, { name: action.name });
  case 'ADD_PARAM':
    return Object.assign({}, state, { params: [...state.params, action.param] });
  default:
    return state;
  }
};

const experiments = (state = [], action) => {
  switch (action.type) {
  case 'CREATE_EXPERIMENT':
    return [...state, action.experiment];
  default:
    return state;
  }
}

const ecInitialState = {
  edit: expInitialState,
  experiments: [],
}
const experimentContainer = (state = ecInitialState, action) => {
  switch (action.type) {
  case 'CREATE_EXPERIMENT':
    action.experiment = state.edit;
    return {
      edit: expInitialState,
      experiments: experiments(state.experiments, action),
    };
  case 'UPDATE_NAME':
    return Object.assign({}, state, {
      edit: experiment(state.edit, action),
    });
  default:
    return state;
  }
}

export default experimentContainer;
