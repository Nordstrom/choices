import React from 'react';

const AppliedLabelList = ({ labels, onLabelClick }) => {
  const labelList = labels.map(label => 
    <li key={label.id} onClick={() => onLabelClick(label.id)}>{label.name}</li>
  )
  return (
    <div>
      <label>Applied Labels:</label>
      <ul>
        {labelList}
      </ul>
    </div>
  )
};

export default AppliedLabelList;
