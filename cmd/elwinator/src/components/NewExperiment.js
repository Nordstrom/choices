import React from 'react';
import { connect } from 'react-redux';
import { browserHistory } from 'react-router';

import { addExperiment } from '../actions';

const NewExperiment = ({ namespaceName, addExperiment }) => {
  let input;
  return (
    <div>
      <form onSubmit={e => {
        e.preventDefault();
        if (!input.value.trim()) {
          return;
        }
        addExperiment(namespaceName, input.value);
        browserHistory.push(`/namespace/${namespaceName}/experiment/${input.value}`);
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

const mapStateToProps = (state, ownProps) => ({
  namespaceName: ownProps.params.namespace,
});

const mapDispatchToProps = ({
  addExperiment,
});

const connected = connect(mapStateToProps, mapDispatchToProps)(NewExperiment);

export default connected;
