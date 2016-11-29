import React, { PropTypes } from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { addChoice, addWeight } from '../actions';
import { paramURL } from '../urls';

const NewChoice = ({
  namespaceName,
  experimentName,
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
      dispatch(addChoice(paramID, choice.value));
      if (isWeighted) {
        dispatch(addWeight(paramID, parseInt(weight.value, 10)));
      }
      if (!redirectOnSubmit) {
        choice.value = '';
        weight.value = '';
        return;
      }
      browserHistory.push(paramURL(paramName));
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
   namespaceName: PropTypes.string.isRequired,
   experimentName: PropTypes.string.isRequired,
   paramName: PropTypes.string.isRequired,
   isWeighted: PropTypes.bool.isRequired,
   dispatch: PropTypes.func.isRequired,
   redirectOnSubmit: PropTypes.bool.isRequired,
}

const mapStateToProps = (state, ownProps) => {
  const paramID = state.entities.params.find(p => p.name === ownProps.paramName).id;
  return {
    paramID,
  };
};

const connected = connect(mapStateToProps)(NewChoice);

export default connected;
