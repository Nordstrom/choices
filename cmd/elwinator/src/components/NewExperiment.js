// @flow
import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import { experimentAdd } from '../actions';
import { namespaceURL } from '../urls';

const NewExperiment = ({ namespaceName, dispatch }) => {
  let input;
  return (
    <form onSubmit={e => {
      e.preventDefault();
      if (!input.value.trim()) {
        return;
      }
      dispatch(experimentAdd(namespaceName, input.value));
      browserHistory.push(namespaceURL(namespaceName));
    }}>
      <div className="form-group">
        <label>Experiment Name</label>
        <input
          type="text"
          className="form-control"
          placeholder="Enter experiment name"
          ref={node => { input = node }}
        />
      </div>
      <button type="submit" className="btn btn-primary">Submit</button>
    </form>
  );
};

NewExperiment.propTypes = {
  namespaceName: React.PropTypes.string.isRequired,
  dispatch: React.PropTypes.func.isRequired,
}

const connected = connect()(NewExperiment);

export default connected;
