import React from 'react';

const UnappliedLabelList = ({ labels, onLabelClick }) => (
  <ul>
    {labels.map(label => 
      <li key={label.id} onClick={() => onLabelClick(label.id)}>{label.name}</li>
    )}
  </ul>
);

export default UnappliedLabelList;
