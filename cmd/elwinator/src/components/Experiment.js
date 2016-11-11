import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import ParamList from './ParamList';
import { namespaceURL, paramNewURL } from '../urls';

const Experiment = ({ namespaceName, experimentName }) => {
  return (
    <div className="container">
      <div className="row"><div className="col-sm-9 col-sm-offset-3"><h1>{experimentName}</h1></div></div>
      <div className="row">
        <div className="col-sm-3">
          <nav>
            <ul className="nav nav-pills nav-stacked">
              <li className="nav-item"><Link to={namespaceURL(namespaceName)} className="nav-link">{namespaceName} - Namespace</Link></li>
              <li className="nav-item"><Link to={paramNewURL(namespaceName, experimentName)}>Create param</Link></li>
            </ul>
          </nav>
        </div>
        <div className="col-sm-9">
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
