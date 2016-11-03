import React from 'react';

import Label from './Label';

const Namespace = (props) => {
  return (
    <div>
    <h1>{props.Name}</h1>
    <Label />
    </div>
  );
}

export default Namespace;
