import { combineReducers } from 'redux';

import params from './params';
import labels from './labels';
import experiment from './experiment';
import choices from './choices';

const reducers = combineReducers({
  params,
  labels,
  experiment,
  choices,
});

export default reducers;
