import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';

import namespaces from './namespaces';
import labels from './labels';
import experiments from './experiments';
import params from './params';

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

export default reducers;
