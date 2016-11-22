import React from 'react';
import ReactDOM from 'react-dom';
import { createStore } from 'redux';
import { Provider } from 'react-redux';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap/dist/css/bootstrap-theme.css';
import { Router,  Route, browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';

import { toNamespace } from './nsconv';
import { namespacesLoaded } from './actions';
import { loadState, saveState } from './Storage';
import App from './App';
import reducers from './reducers';
import NewNamespaceView from './components/NewNamespaceView';
import Namespace from './components/Namespace';
import NewLabelView from './components/NewLabelView';
import NewExperimentView from './components/NewExperimentView';
import Experiment from './components/Experiment';
import NewParamView from './components/NewParamView';
import Param from './components/Param';
import NewChoiceView from './components/NewChoiceView';


const persistedState = loadState();
const store = createStore(
  reducers,
  persistedState,
);


store.subscribe(() => {
  saveState(store.getState());
});

const history = syncHistoryWithStore(browserHistory, store)

const render = () => {
ReactDOM.render(
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={App} />
      <Route path="/n/new" component={NewNamespaceView} />
      <Route path="/n/:namespace" component={Namespace} />
      <Route path="/n/:namespace/l/new" component={NewLabelView} />
      <Route path="/n/:namespace/e/new" component={NewExperimentView} />
      <Route path="/n/:namespace/e/:experiment" component={Experiment} />
      <Route path="/n/:namespace/e/:experiment/p/new" component={NewParamView} />
      <Route path="/n/:namespace/e/:experiment/p/:param" component={Param} />
      <Route path="/n/:namespace/e/:experiment/p/:param/c/new" component={NewChoiceView} />
    </Router>
  </Provider>,
  document.getElementById('root')
);
}

const headers = new Headers({'Accept': 'application/json'});
const req = { method: 'POST', headers: headers, body: JSON.stringify({ environment: "Staging" }) };
const badRequest =  { err: "bad request" };
fetch("/api/v1/all", req)
.then(resp => {
  if (!resp.ok) {
    throw badRequest;
  }
  return resp.json();
})
.then(json => {
  const ns = json.namespaces.map(n => {
    return toNamespace(n);
  });
  store.dispatch(namespacesLoaded(ns));
})
.then(() => {
  render()
})
.catch(e => {
  console.log(e);
  render();
});
