// @flow
import React from 'react';
import { Link, browserHistory } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import LabelList from './LabelList';
import NewLabel from './NewLabel';
import Segment from './Segment';
import ExperimentList from './ExperimentList';
import { rootURL, labelNewURL, experimentNewURL } from '../urls';
import { namespaceDelete } from '../actions';
import { getNamespace } from '../reducers/namespaces';
import { combinedSegments, getExperiments } from '../reducers/experiments';
import { getLabels } from '../reducers/labels';

const Namespace = ({ ns, labels, namespaceSegments, dispatch }) => {
  const llProps = { namespaceName: ns.name, labels: ns.labels.map(lid => labels.find(l => lid === l.id)) };
  const newLabelProps = { namespaceName: ns.name, dispatch, redirectOnSubmit: false };
  return (
    <div className="container">
      <div className="row"><h1>{ ns.name }</h1></div>
      <div className="row">
        <NavSection>
          <Link to={ labelNewURL(ns.name) } className="nav-link">New label</Link>
          <Link to={ experimentNewURL(ns.name) } className="nav-link">New Experiment</Link>
        </NavSection>
        <div className="col-sm-9">
          <h2>Labels</h2>
          <LabelList { ...llProps } />
          <NewLabel { ...newLabelProps } />
          <h2>Segments</h2>
          <Segment namespaceSegments={namespaceSegments} />
          <h2>Experiments</h2>
          <ExperimentList namespaceName={ ns.name } />
          <Link
            to={experimentNewURL(ns.name)}
            className="btn btn-default"
            role="button"
          >Create new experiment</Link><br />
          <button className="btn btn-warning" onClick={() => {
            dispatch(namespaceDelete(ns.name));
            browserHistory.push(rootURL());
          }}>Delete namespace {ns.name}</button>
        </div>
      </div>
    </div>
  );
};

const mapStateToProps = (state, ownProps) => {
  const ns = getNamespace(state.entities.namespaces, ownProps.params.namespace);
  const namespaceSegments = combinedSegments(state.entities.experiments, ns.experiments);
  const experiments = getExperiments(state.entities.experiments, ns.experiments);
  const labels = getLabels(state.entities.labels, ns.labels);
  return {
    ns,
    namespaceSegments,
    experiments,
    labels,
  }
};

const connected = connect(mapStateToProps)(Namespace);
 
export default connected;
