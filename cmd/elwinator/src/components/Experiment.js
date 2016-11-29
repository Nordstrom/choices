import React from 'react';
import { Link, browserHistory } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import SegmentInput from './SegmentInput';
import Segment from './Segment';
import ParamList from './ParamList';
import { namespaceURL, paramNewURL } from '../urls';
import { experimentDelete } from '../actions';
import { availableSegments, combinedSegments } from '../reducers/experiments';
import { getParams } from '../reducers/params';

const Experiment = ({ ns, exp, freeSegments, namespaceSegments, params, dispatch }) => {
  const siProps = {
    namespaceName: ns.name,
    experimentName: exp.name,
    experimentID: exp.id,
    numSegments: exp.numSegments,
    availableSegments: freeSegments,
    namespaceSegments,
    redirectOnSubmit: false,
  }
  return (
    <div className="container">
      <div className="row">
        <div className="col-sm-9 col-sm-offset-3"><h1>{exp.name}</h1></div>
      </div>
      <div className="row">
        <NavSection>
          <Link to={namespaceURL(ns.name)} className="nav-link">{ns.name} - Namespace</Link>
          <Link to={paramNewURL(exp.id)} className="nav-link">Create param</Link>
        </NavSection>
        <div className="col-sm-9">
          <h2>Segments</h2>
          <SegmentInput {...siProps } />
          <Segment
            namespaceSegments={namespaceSegments}
            experimentSegments={exp.segments}
          />
          <h2>Params</h2>
          <ParamList params={params} />
          <Link
            to={paramNewURL(exp.id)}
            className="btn btn-default"
            role="button">Create new param</Link><br />
          <button className="btn btn-warning" onClick={() => {
            dispatch(experimentDelete(ns.name, exp.id));
            browserHistory.push(namespaceURL(ns.name));
          }}>Delete experiment {exp.name}</button>
        </div>
      </div>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.entities.namespaces.find(n => n.experiments.find(e => e === ownProps.params.experiment));
  const exp = state.entities.experiments.find(e => e.id === ownProps.params.experiment);
  const freeSegments = availableSegments(state.entities.experiments, ns.experiments.filter(eid => eid !== exp.id));
  const namespaceSegments = combinedSegments(state.entities.experiments, ns.experiments.filter(eid => eid !== exp.id));
  const params = getParams(state.entities.params, exp.params);
  return {
    ns,
    exp,
    freeSegments,
    namespaceSegments,
    params,
  }
};

const connected = connect(
  mapStateToProps,
)(Experiment);

export default connected;
