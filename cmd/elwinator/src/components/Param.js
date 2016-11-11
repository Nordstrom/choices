import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import { choiceNewURL } from '../urls';

const Param = ({ namespaceName, experimentName, p }) => {
  return (
    <div className="container">
      <Link to={choiceNewURL(namespaceName, experimentName, p.name)}>Create new choice</Link>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  const exp = ns.experiments.find(e => e.name === ownProps.params.experiment);
  const p = exp.params.find(p => p.name === ownProps.params.param);
  return {
    namespaceName: ns.name,
    experimentName: exp.name,
    p,
  }
}

const connected = connect(mapStateToProps)(Param);

export default connected;
