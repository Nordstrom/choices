import React from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { addParam } from '../actions';

const NewParam = ({ namespaceName, experimentID, experimentName, dispatch }) => {
  let input;
  return (
    <form onSubmit={e => {
      e.preventDefault();
      if (!input.value.trim()) {
        return;
      }
      dispatch(addParam(experimentID, input.value));
      browserHistory.push(`/n/${namespaceName}/e/${experimentName}/p/${input.value}`)
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
  namespaceName: React.PropTypes.string.isRequired,
  experimentName: React.PropTypes.string.isRequired,
  dispatch: React.PropTypes.func.isRequired,
}

const mapStateToProps = (state, ownProps) => {
  const experimentID = state.entities.experiments.find(e => e.name === ownProps.experimentName).id;
  return {
    experimentID,
  };
};

const connected = connect(mapStateToProps)(NewParam);

export default connected;
