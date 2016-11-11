import React from 'react';
import { browserHistory } from 'react-router';

import { paramURL } from '../urls';

const NewChoice = ({ namespaceName, experimentName, paramName, isWeighted, addChoice, addWeight, redirectOnSubmit }) => {
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
      addChoice(namespaceName, experimentName, paramName, choice.value);
      if (isWeighted) {
        addWeight(namespaceName, experimentName, paramName, weight.value);
      }
      if (!redirectOnSubmit) {
        choice.value = '';
        weight.value = '';
        return;
      }
      browserHistory.push(paramURL(namespaceName, experimentName, paramName));
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

export default NewChoice;
