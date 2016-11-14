import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import ParamList from './ParamList';
import { rootURL, namespaceURL, paramNewURL } from '../urls';

const Experiment = ({ namespaceName, experimentName }) => {
  return (
    <div className="container">
      <div className="row"><div className="col-sm-9 col-sm-offset-3"><h1>{experimentName}</h1></div></div>
      <div className="row">
        <NavSection>
          <Link to={ rootURL() }>Home</Link>
          <Link to={namespaceURL(namespaceName)} className="nav-link">{namespaceName} - Namespace</Link>
          <Link to={paramNewURL(namespaceName, experimentName)}>Create param</Link>
        </NavSection>
        <div className="col-sm-9">
          <h2>Params</h2>
          <ParamList namespaceName={namespaceName} experimentName={experimentName} />  
        </div>
      </div>
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
