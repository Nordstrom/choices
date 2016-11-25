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

const Namespace = ({ ns, dispatch }) => {
  const llProps = { namespaceName: ns.name, labels: ns.labels };
  const newLabelProps = { namespaceName: ns.name, dispatch, redirectOnSubmit: false };
  const nsSegments = ns.experiments.reduce((prev, e) => {
      e.segments.forEach((seg, i) => {
        prev[i] |= seg;
      });
      return prev;
    }, new Uint8Array(16).fill(0));
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
          <Segment namespaceSegments={nsSegments} />
          <h2>Experiments</h2>
          <ExperimentList namespaceName={ ns.name } />
          <Link to={experimentNewURL(ns.name)} className="btn btn-default" role="button">Create new experiment</Link><br />
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
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  return {
    ns,
  }
};

const connected = connect(mapStateToProps)(Namespace);
 
export default connected;