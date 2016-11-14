import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import LabelList from './LabelList';
import ExperimentList from './ExperimentList';
import { labelNewURL, experimentNewURL } from '../urls';
import { toggleLabel } from '../actions';

const Namespace = ({ ns, toggleLabel }) => {
  const llProps = {namespaceName: ns.name, labels: ns.labels, toggleLabel, };
  return (
    <div className="container">
      <h1>{ns.name}</h1>
      <h2>Labels</h2>
      <LabelList {...llProps} />
      <Link to={labelNewURL(ns.name)}>New label</Link>
      <h2>Experiments</h2>
      <ExperimentList namespaceName={ns.name} />
      <Link to={experimentNewURL(ns.name)}>New Experiment</Link>
    </div>
  );
};

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  return {
    ns,
  }
};

const mapDispatchToProps = {
  toggleLabel,
}

const connected = connect(mapStateToProps, mapDispatchToProps)(Namespace);
 
export default connected;