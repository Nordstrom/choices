// @flow
import React from 'react';
import { browserHistory } from 'react-router';
import { connect } from 'react-redux';

import { namespaceURL } from '../urls';
import { namespaceAdd } from '../actions';

const nameLength = 7;
const alphabet = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'

const randomName = (len: number): string => {
  let out = '';
  for (let i = 0; i < len; i++) {
    const randIndex = getRandomInt(0, alphabet.length)
    out += alphabet[randIndex];
  }
  return out;
}

function getRandomInt(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min)) + min;
}

const NewNamespace = ({ dispatch }) => {
  let input;
  return (
    <form onSubmit={e => {
      e.preventDefault();
      if (!input.value.trim()) {
        return;
      }
      dispatch(namespaceAdd(input.value));
      browserHistory.push(namespaceURL(input.value));
    }}>
      <div className="form-group">
        <label>Namespace Name</label>
        <input type="text" className="form-control" ref={node => input = node}/>
      </div>
      <button type="button" className="btn btn-primary" onClick={
        () => input.value = randomName(nameLength)
      }>Random Name</button><br />
      <button type="submit" className="btn btn-primary" >Create namespace</button>
    </form>
  );
};

const connected = connect()(NewNamespace);

export default connected;
