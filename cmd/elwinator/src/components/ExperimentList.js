import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import { experimentURL } from '../urls';

const ExperimentList = ({ namespaceName, experiments }) => {
  const exps = experiments.map(e => <li key={e.name}><Link to={experimentURL(namespaceName, e.name)}>{e.name}</Link></li>)
  return (
    <ul>
      {exps}
    </ul>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.namespaceName);
  return {
    namespaceName: ns.name,
    experiments: ns.experiments,
  }
}

const connected = connect(mapStateToProps)(ExperimentList);

export default connected;
