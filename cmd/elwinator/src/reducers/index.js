import { combineReducers } from 'redux';

import params from './params';
import labels from './labels';

const reducers = combineReducers({
  params,
  labels,
});

export default reducers;
