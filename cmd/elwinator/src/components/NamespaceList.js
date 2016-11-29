import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import { namespaceURL } from '../urls';
import { namespaceDelete } from '../actions';
import { getExperiments } from '../reducers/experiments';

const NamespaceList = ({ namespaces, dispatch }) => {
  const namespaceList = namespaces
  .filter(n => !n.delete)
  .map((n, i) => 
    <tr key={n.name}>
      <td>{i + 1}</td>
      <td><Link to={namespaceURL(n.name)}>{n.name}</Link></td>
      <td>{n.experiments.map(e => e.name).join(', ')}</td>
      <td>
        <button
          className="btn btn-default btn-xs"
          onClick={() => dispatch(namespaceDelete(n.name))}
        >&times;</button>
      </td>
    </tr>
  )
  return (
    <div className="row">
      <div className="col-sm-12">
        <table className="table table-striped">
          <thead>
            <tr>
              <th>#</th>
              <th>Namespace</th>
              <th>Experiment(s)</th>
              <th>Delete</th>
            </tr>
          </thead>
          <tbody>
            {namespaceList}
          </tbody>
        </table>
      </div>
    </div>
  );
};

const mapStateToProps = (state) => {
  const ns = state.entities.namespaces.map(n => ({
    ...n,
    experiments: getExperiments(state.entities.experiments, n.experiments),
  }));
  return {
    namespaces: ns,
  }
};

const connected = connect(mapStateToProps)(NamespaceList);

export default connected;
