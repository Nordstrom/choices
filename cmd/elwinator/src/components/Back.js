import React from 'react';
import { browserHistory } from 'react-router';

const Back = props =>
  <a className="nav-link" onClick={e => {
    e.preventDefault(); browserHistory.goBack();
  }}>Back</a>;

export default Back
