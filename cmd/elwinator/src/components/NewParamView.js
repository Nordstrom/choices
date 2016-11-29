import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import { namespaceURL, experimentURL } from '../urls';
import NewParam from './NewParam';
import { getExperiment } from '../reducers/experiments';

const NewParamView = ({ namespaceName, experimentID, experimentName }) => {
  return (
    <div className="container">
      <div className="row"><h1>Create a new param</h1></div>
      <div className="row">
        <NavSection>
          <Link
            to={ namespaceURL(namespaceName) }
            className="nav-link"
          >{namespaceName} - Namespace</Link>
          <Link
            to={ experimentURL(experimentID) }
            className="nav-link"
          >{experimentName} - Experiment </Link>
        </NavSection>
        <div className="col-sm-9">
          <NewParam experimentID={experimentID} />
        </div>
      </div>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const exp = getExperiment(state.entities.experiments, ownProps.params.experiment);
  return {
    namespaceName: exp.namespace,
    experimentID: exp.id,
    experimentName: exp.name,
  }
}

const connected = connect(mapStateToProps)(NewParamView);

export default connected;
