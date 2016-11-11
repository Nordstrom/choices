import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import { removeParam } from '../actions';

const ParamList = ({ namespaceName, experimentName, params, onParamClick }) => (
  <ul>
    {params.map(param => 
      <li key={param.name} >
        <Link to={`/namespace/${namespaceName}/experiment/${experimentName}/param/${param.name}`}>{param.name}</Link>
      </li>
    )}
  </ul>
);

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.namespaceName);
  const exp = ns.experiments.find(e => e.name === ownProps.experimentName);
  return {
    params: exp.params,
  }
};

const mapDispatchToProps = ({
  onParamClick: removeParam,
});

const connected = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ParamList);

export default connected;
