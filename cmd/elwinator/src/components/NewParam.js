import React from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { addParam } from '../actions';
import { experimentURL } from '../urls';

const NewParam = ({ experimentID, dispatch }) => {
  let input;
  return (
    <form onSubmit={e => {
      e.preventDefault();
      if (!input.value.trim()) {
        return;
      }
      dispatch(addParam(experimentID, input.value));
      browserHistory.push(experimentURL(experimentID));
    }}>
      <div className="form-group">
        <label>Param Name</label>
        <input
          type="text"
          className="form-control"
          placeholder="Enter param name"
          ref={node => { input = node }}
        />
      </div>
      <button type="submit" className="btn btn-primary">Submit</button>
    </form>
  );
};

NewParam.propTypes = {
  experimentID: React.PropTypes.string.isRequired,
  dispatch: React.PropTypes.func.isRequired,
}

const connected = connect()(NewParam);

export default connected;
