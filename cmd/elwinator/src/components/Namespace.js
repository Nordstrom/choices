import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import LabelList from './LabelList';
import ExperimentList from './ExperimentList';
import { labelNewURL, experimentNewURL } from '../urls';

const Namespace = ({ namespaceName }) => {
  return (
    <div className="container">
      <h1>{namespaceName}</h1>
      <h2>Labels</h2>
      <LabelList namespaceName={namespaceName} />
      <Link to={labelNewURL(namespaceName)}>New label</Link>
      <h2>Experiments</h2>
      <ExperimentList namespaceName={namespaceName} />
      <Link to={experimentNewURL(namespaceName)}>New Experiment</Link>
    </div>
  );
};

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  return {
    namespaceName: ns.name,
  }
};

const connected = connect(mapStateToProps)(Namespace);
 
export default connected;