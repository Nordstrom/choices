import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';

import namespaces from './namespaces';

const reducers = combineReducers({
  namespaces,
  routing: routerReducer,
});

export default reducers;
