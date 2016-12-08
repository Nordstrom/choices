import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';

import namespaces from './namespaces';
import labels from './labels';
import experiments from './experiments';
import params from './params';
import changes from './changes';

const entitiesReducer = combineReducers({
  namespaces,
  labels,
  experiments,
  params,
  changes,
});

const reducers = combineReducers({
  entities: entitiesReducer,
  routing: routerReducer,
});

export default reducers;
