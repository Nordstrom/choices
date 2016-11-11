import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';

import namespace from './namespace';

const reducers = combineReducers({
  namespace,
  routing: routerReducer,
});

export default reducers;
