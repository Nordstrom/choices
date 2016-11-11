import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import LabelList from './LabelList';
import ExperimentList from './ExperimentList';

const Namespace = ({ namespaceName }) => {
  return (
    <div className="container">
      <h1>{namespaceName}</h1>
      <LabelList namespaceName={namespaceName} />
      <Link to={`/n/${namespaceName}/l/new`}>New label</Link>
      <ExperimentList namespaceName={namespaceName} />
      <Link to={`/n/${namespaceName}/e/new`}>New Experiment</Link>
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