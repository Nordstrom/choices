import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import { addExperiment } from '../actions';

const NewExperiment = ({ namespaceName, dispatch }) => {
  let input;
  return (
    <div className="container">
      <form onSubmit={e => {
        e.preventDefault();
        if (!input.value.trim()) {
          return;
        }
        dispatch(addExperiment(namespaceName, input.value));
        browserHistory.push(`/n/${namespaceName}/e/${input.value}`);
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
    </div>
  );
};

NewExperiment.propTypes = {
  namespaceName: React.PropTypes.string.isRequired,
  dispatch: React.PropTypes.func.isRequired,
}

const connected = connect()(NewExperiment);

export default connected;
