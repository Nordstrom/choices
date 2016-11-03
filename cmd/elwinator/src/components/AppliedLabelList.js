import React from 'react';

const AppliedLabelList = ({ labels, onLabelClick }) => (
  <ul>
    {labels.map(label => 
      <li key={label.id} onClick={() => onLabelClick(label.id)}>{label.name}</li>
    )
    }
  </ul>
);

export default AppliedLabelList;
