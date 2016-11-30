import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';

import namespaces from './namespaces';
import labels from './labels';
import experiments, { getExperiment } from './experiments';
import params, { getParam } from './params';

const entities = combineReducers({
  namespaces,
  labels,
  experiments,
  params,
});

const reducers = combineReducers({
  entities,
  routing: routerReducer,
});

const changeList = (state={}, namespace, action) => {
  if (state[namespace]) {
    return {
      ...state,
      [namespace]: [...state[namespace], action],
    }
  }
  return {
    ...state,
    [namespace]: [action],
  }
}

const changeLog = (state = {}, action) => {
  const previousState = { entities: state.entities, routing: state.routing };
  let nextState = reducers(previousState, action);
  if (previousState === nextState) {
    return state;
  }
  if (action.namespace) {
    return {
      ...nextState,
      changes: changeList(state.changes, action.namespace, action),
    };
  } else if (action.experiment) {
    const exp = getExperiment(nextState.entities.experiments, action.experiment);
    return {
      ...nextState,
      changes: changeList(state.changes, exp.namespace, action),
    };
  } else if (action.param) {
    const param = getParam(nextState.entities.params, action.param);
    const exp = getExperiment(nextState.entities.experiments, param.experiment);
    return {
      ...nextState,
      changes: changeList(state.changes, exp.namespace, action),
    };
  }
  return {
    ...nextState,
    changes: state.changes || {},
  };
};

export default changeLog;
