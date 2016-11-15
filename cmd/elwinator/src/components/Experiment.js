import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import SegmentInput from './SegmentInput';
import ParamList from './ParamList';
import { rootURL, namespaceURL, paramNewURL } from '../urls';
import { experimentNumSegments, experimentPercent } from '../actions';

const Experiment = ({ namespaceName, exp, experimentNumSegments, experimentPercent }) => {
  const siProps = {
    namespaceName,
    experimentName: exp.name,
    numSegments: exp.numSegments,
    dirtySegments: exp.dirtySegments,
    redirectOnSubmit: false,
    experimentNumSegments,
    experimentPercent,
  }
  return (
    <div className="container">
      <div className="row"><div className="col-sm-9 col-sm-offset-3"><h1>{exp.name}</h1></div></div>
      <div className="row">
        <NavSection>
          <Link to={ rootURL() } className="nav-link">Home</Link>
          <Link to={namespaceURL(namespaceName)} className="nav-link">{namespaceName} - Namespace</Link>
          <Link to={paramNewURL(namespaceName, exp.name)} className="nav-link">Create param</Link>
        </NavSection>
        <div className="col-sm-9">
          <h2>Segments</h2>
          <SegmentInput {...siProps } />
          <h2>Params</h2>
          <ParamList namespaceName={namespaceName} experimentName={exp.name} />  
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
    exp,
  }
};

const mapDispatchToProps = {
  experimentNumSegments,
  experimentPercent,
}
const connected = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Experiment);

export default connected;
