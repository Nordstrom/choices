import { combineReducers } from 'redux';

import params from './params';
import labels from './labels';
import experiment from './experiment';

const reducers = combineReducers({
  params,
  labels,
  experiment,
});

export default reducers;
