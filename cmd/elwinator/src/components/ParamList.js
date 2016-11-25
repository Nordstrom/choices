import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import { paramDelete } from '../actions';
import { paramURL } from '../urls';

const ParamList = ({ namespaceName, experimentName, params, dispatch }) => (
  <table className="table table-striped">
    <thead>
      <tr>
        <th>#</th>
        <th>Param</th>
        <th>Choices</th>
        <th>Delete</th>
      </tr>
    </thead>
    <tbody>
    {params.map((param, i) => 
      <tr key={param.name} >
        <td>{i + 1}</td>
        <td><Link to={paramURL(namespaceName, experimentName, param.name)}>{param.name}</Link></td>
        <td>{param.choices.join(', ')}</td>
        <td><button className="btn btn-default btn-xs" onClick={
          () => dispatch(paramDelete(namespaceName, experimentName, param.name))
        }>&times;</button></td>
      </tr>
    )}
    </tbody>
  </table>
);

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.namespaceName);
  const exp = ns.experiments.find(e => e.name === ownProps.experimentName);
  return {
    params: exp.params,
  }
};

const connected = connect(
  mapStateToProps,
)(ParamList);

export default connected;
