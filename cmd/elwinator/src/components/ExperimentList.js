import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import { experimentURL } from '../urls';
import { experimentDelete } from '../actions';
import { getExperiments } from '../reducers/experiments';

const ExperimentList = ({ namespaceName, experiments, dispatch }) => {
  const exps = experiments.map((e, i) =>
    <tr key={e.name}>
      <td>{i + 1}</td>
      <td><Link to={experimentURL(namespaceName, e.name)}>{e.name}</Link></td>
      <td>{e.params.map(p => p.name).join(', ')}</td>
      <td>
        <button
          className="btn btn-default btn-xs"
          onClick={() => dispatch(experimentDelete(namespaceName, e.name))}
        >&times;</button>
      </td>
    </tr>
  );
  return (
    <table className="table table-striped">
    <thead>
      <tr>
        <th>#</th>
        <th>Experiment</th>
        <th>Param(s)</th>
        <th>Delete</th>
      </tr>
    </thead>
    <tbody>
      {exps}
    </tbody>
    </table>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.entities.namespaces.find(n => n.name === ownProps.namespaceName);
  const experiments = getExperiments(state.entities.experiments, ns.experiments);
  return {
    namespaceName: ns.name,
    experiments,
  }
}

const connected = connect(mapStateToProps)(ExperimentList);

export default connected;
