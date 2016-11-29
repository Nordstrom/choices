import React from 'react';
import ReactDOM from 'react-dom';
import { createStore } from 'redux';
import { Provider } from 'react-redux';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap/dist/css/bootstrap-theme.css';
import { Router,  Route, browserHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';

// import { toNamespace } from './nsconv';
// import { entitiesLoaded } from './actions';
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
      <Route path="/new-namespace" component={NewNamespaceView} />
      <Route path="/n/:namespace" component={Namespace} />
      <Route path="/n/:namespace/new-label" component={NewLabelView} />
      <Route path="/n/:namespace/new-experiment" component={NewExperimentView} />
      <Route path="/e/:experiment" component={Experiment} />
      <Route path="/e/:experiment/new-param" component={NewParamView} />
      <Route path="/p/:param" component={Param} />
      <Route path="/p/:param/new-choice" component={NewChoiceView} />
    </Router>
  </Provider>,
  document.getElementById('root')
);
}

// const headers = new Headers({'Accept': 'application/json'});
// const req = { method: 'POST', headers: headers, body: JSON.stringify({ environment: "Staging" }) };
// const badRequest =  { err: "bad request" };
// fetch("/api/v1/all", req)
// .then(resp => {
//   if (!resp.ok) {
//     throw badRequest;
//   }
//   return resp.json();
// })
// .then(json => {
//   const ns = json.namespaces.map(n => {
//     return toNamespace(n);
//   });
//   store.dispatch(entitiesLoaded(ns));
// })
// .then(() => {
//   render()
// })
// .catch(e => {
//   console.log(e);
//   render();
// });
render();
