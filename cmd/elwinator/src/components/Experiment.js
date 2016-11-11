import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import ParamList from './ParamList';
import { paramNewURL } from '../urls';

const Experiment = ({ namespaceName, experimentName }) => {
  return (
    <div className="container">
    <h1>{experimentName} - Experiment</h1>
    <ParamList namespaceName={namespaceName} experimentName={experimentName} />
    <Link to={paramNewURL(namespaceName, experimentName)}>Create param</Link>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  const exp = ns.experiments.find(e => e.name === ownProps.params.experiment);
  return {
    namespaceName: ns.name,
    experimentName: exp.name,
  }
};

const connected = connect(
  mapStateToProps,
)(Experiment);

export default connected;
