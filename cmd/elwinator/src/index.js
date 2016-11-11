import React from 'react';
import ReactDOM from 'react-dom';
import { createStore } from 'redux';
import { Provider } from 'react-redux';
import 'bootstrap/dist/css/bootstrap.css';
// import 'bootstrap/dist/css/bootstrap-theme.css';
import { Router,  Route, browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';

// import { loadState, saveState } from './Storage';
import App from './App';
import reducers from './reducers';
import NewNamespace from './components/NewNamespace';
import Namespace from './components/Namespace';
import NewExperiment from './components/NewExperiment';
import Experiment from './components/Experiment';
import NewParam from './components/NewParam';
import Param from './components/Param';
import NewChoice from './components/NewChoice';

// const persistedState = loadState();
const store = createStore(reducers);
  // persistedState,
// );

// store.subscribe(() => {
//   saveState(store.getState());
// });

const history = syncHistoryWithStore(browserHistory, store)

ReactDOM.render(
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={App} />
      <Route path="/namespace/new" component={NewNamespace} />
      <Route path="/namespace/:namespace" component={Namespace} />
      <Route path="/namespace/:namespace/experiment/new" component={NewExperiment} />
      <Route path="/namespace/:namespace/experiment/:experiment" component={Experiment} />
      <Route path="/namespace/:namespace/experiment/:experiment/param/new" component={NewParam} />
      <Route path="/namespace/:namespace/experiment/:experiment/param/:param/choice/new" component={NewChoice} />
      <Route path="/namespace/:namespace/experiment/:experiment/param/:param" component={Param} />
    </Router>
  </Provider>,
  document.getElementById('root')
);
