import React from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { paramURL } from '../urls';
import { addChoice, addWeight } from '../actions';

const NewChoice = ({ namespaceName, experimentName, paramName, isWeighted, addChoice, addWeight }) => {
  let choice;
  let weight;
  return (
    <div className="container">
      <form onSubmit={e => {
        e.preventDefault();
        if (!choice.value.trim()) {
          return;
        }
        if (isWeighted && !weight.value.trim()) {
          return;
        }
        addChoice(namespaceName, experimentName, paramName, choice.value);
        if (isWeighted) {
          addWeight(namespaceName, experimentName, paramName, weight.value);
        }
        browserHistory.push(paramURL(namespaceName, experimentName, paramName));
      }}>
        <div className="form-group">
          <label>Choice</label>
          <input type="text"
            className="form-control"
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
            ref={node => weight = node}
          />
        </div>
        <button type="submit" className="btn btn-primary">Create choice</button>
      </form>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  const exp = ns.experiments.find(e => e.name === ownProps.params.experiment);
  const param = exp.params.find(p => p.name === ownProps.params.param);
  return {
    namespaceName: ns.name,
    experimentName: exp.name,
    paramName: param.name,
    isWeighted: param.isWeighted
  }
}

const mapDispatchToProps = {
  addChoice,
  addWeight,
}

const connected = connect(mapStateToProps, mapDispatchToProps)(NewChoice);

export default connected;
