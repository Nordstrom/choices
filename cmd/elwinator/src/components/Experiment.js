import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import SegmentInput from './SegmentInput';
import Segment from './Segment';
import ParamList from './ParamList';
import { namespaceURL, paramNewURL } from '../urls';
import { experimentNumSegments, experimentPercent } from '../actions';

const Experiment = ({ ns, exp, experimentNumSegments, experimentPercent }) => {
  const siProps = {
    namespaceName: ns.name,
    experimentName: exp.name,
    numSegments: exp.numSegments,
    availableSegments: ns.experiments.reduce((prev, e) => {
      if (e.name === exp.name) {
        return prev;
      }
      return prev - e.numSegments;
    }, 128),
    namespaceSegments: ns.experiments.reduce((prev, e) => {
      if (e.name === exp.name) {
        return prev;
      }
      e.segments.forEach((seg, i) => {
        prev[i] |= seg
      });
      return prev;
    }, new Uint8Array(16).fill(0)),
    redirectOnSubmit: false,
    experimentNumSegments,
  }
  return (
    <div className="container">
      <div className="row"><div className="col-sm-9 col-sm-offset-3"><h1>{exp.name}</h1></div></div>
      <div className="row">
        <NavSection>
          <Link to={namespaceURL(ns.name)} className="nav-link">{ns.name} - Namespace</Link>
          <Link to={paramNewURL(ns.name, exp.name)} className="nav-link">Create param</Link>
        </NavSection>
        <div className="col-sm-9">
          <h2>Segments</h2>
          <SegmentInput {...siProps } />
          <Segment namespaceSegments={siProps.namespaceSegments} experimentSegments={exp.segments} />
          <h2>Params</h2>
          <ParamList namespaceName={ns.name} experimentName={exp.name} />  
        </div>
      </div>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  const exp = ns.experiments.find(e => e.name === ownProps.params.experiment);
  return {
    ns,
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
