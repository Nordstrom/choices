import React from 'react';
import { connect } from 'react-redux';

import LabelList from './LabelList';
import ExperimentList from './ExperimentList';

const Namespace = ({ namespaceName, params, children }) => {
  return (
    <div>
      <h1>{namespaceName}</h1>
      <LabelList />
      <ExperimentList />
    </div>
  );
};

const mapStateToProps = (state) => ({
  namespaceName: state.namespace.name,
});

const connected = connect(mapStateToProps)(Namespace);
 
export default connected;