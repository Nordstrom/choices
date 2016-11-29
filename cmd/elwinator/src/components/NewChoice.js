import React, { PropTypes } from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { paramAddChoice, paramAddWeight } from '../actions';
import { paramURL } from '../urls';

const NewChoice = ({
  paramID,
  paramName,
  isWeighted,
  dispatch,
  redirectOnSubmit
}) => {
  let choice;
  let weight;
  return (
    <form onSubmit={e => {
      e.preventDefault();
      if (!choice.value.trim()) {
        return;
      }
      if (isWeighted && !weight.value.trim()) {
        return;
      }
      dispatch(paramAddChoice(paramID, choice.value));
      if (isWeighted) {
        dispatch(paramAddWeight(paramID, parseInt(weight.value, 10)));
      }
      if (!redirectOnSubmit) {
        choice.value = '';
        weight.value = '';
        return;
      }
      browserHistory.push(paramURL(paramID));
    }}>
      <div className="form-group">
        <label>Choice</label>
        <input type="text"
          className="form-control"
          placeholder="Enter a choice"
          ref={node => choice = node}
        />
      </div>
      <div className="form-group">
        <label>Weight</label>
        <input type="number"
          className="form-control"
          min="1"
          max="100"
          disabled={!isWeighted}
          placeholder="Enter a weight"
          ref={node => weight = node}
        />
      </div>
      <button type="submit" className="btn btn-primary">Create choice</button>
    </form>
  );
}

NewChoice.propTypes = {
  paramID: PropTypes.string.isRequired,
  paramName: PropTypes.string.isRequired,
  isWeighted: PropTypes.bool.isRequired,
  dispatch: PropTypes.func.isRequired,
  redirectOnSubmit: PropTypes.bool.isRequired,
}

const connected = connect()(NewChoice);

export default connected;
