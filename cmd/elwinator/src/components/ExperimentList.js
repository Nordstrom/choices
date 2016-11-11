import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

const ExperimentList = ({ namespaceName, experiments }) => {
  const exps = experiments.map(e => <li><Link to={`/namespace/${namespaceName}/experiment/${e.name}`}>{e.name}</Link></li>)
  return (
    <ul>
      {exps}
    </ul>
  );
}

const mapStateToProps = (state, ownProps) => {
  return {
    namespaceName: state.namespace.name,
    experiments: state.namespace.experiments,
  }
}

const connected = connect(mapStateToProps)(ExperimentList);

export default connected;
