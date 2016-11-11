import React from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { addParam } from '../actions';

const NewParam = ({ namespaceName, experimentName, addParam }) => {
  let input;
  return (
    <div className="container">
      <form onSubmit={e => {
        e.preventDefault();
        if (!input.value.trim()) {
          return;
        }
        addParam(namespaceName, experimentName, input.value);
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
    </div>
  );
};

const mapStateToProps = (state, ownProps) => {
  return {
    namespaceName: ownProps.params.namespace,
    experimentName: ownProps.params.experiment,
  }
};

const mapDispatchToProps = {
  addParam,
};

const connected = connect(mapStateToProps, mapDispatchToProps)(NewParam);

export default connected;
