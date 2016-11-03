import React from 'react';

const UnappliedLabelList = ({ labels, onLabelClick }) => {
  const labelsList = labels.map(label => 
      <li key={label.id} onClick={() => onLabelClick(label.id)}>{label.name}</li>
  )
  return (
    <div>
      <label>Unapplied labels:</label>
      <ul>
        {labelsList}
      </ul>
    </div>
  );
};

export default UnappliedLabelList;
