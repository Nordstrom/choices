import React from 'react';

import Label from './Label';
import ExperimentContainer from './ExperimentContainer';

const Namespace = (props) => {
  return (
    <div>
    <h1>{props.Name}</h1>
    <Label />
    <ExperimentContainer />
    </div>
  );
}

export default Namespace;
