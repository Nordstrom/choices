import React from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { namespaceURL } from '../urls';
import { addNamespace } from '../actions';

const NewNamespace = ({ dispatch }) => {
  let input;
  return (
    <div className="container">
      <form onSubmit={e => {
        e.preventDefault();
        if (!input.value.trim()) {
          return;
        }
        dispatch(addNamespace(input.value));
        browserHistory.push(namespaceURL(input.value));
      }}>
        <div className="form-group">
          <label>Namespace Name</label>
          <input type="text" className="form-control" ref={node => input = node}/>
        </div>
        <button type="submit" className="btn btn-primary" >Create namespace</button>
      </form>
    </div>
  );
};

const connected = connect()(NewNamespace);

export default connected;
