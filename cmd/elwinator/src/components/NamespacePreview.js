import React from 'react';

import LabelList from './LabelList';
import Experiments from './Experiments';

const NamespacePreview = (props) => {
  return (
    <div className="container">
    <h2>{props.Name}</h2>
    <h3>Labels</h3>
    <LabelList />
    <h3>Experiments</h3>
    <Experiments />
    </div>
  );
}

export default NamespacePreview;
