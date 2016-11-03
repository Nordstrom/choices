import React from 'react';
import './Values.css';

export const Values = props => {
  if (!props.values) {
    return false;
  }
  const values = props.values.choices.map((choice, i) => {
    const style = {
      flex: `${props.values.weights ? props.values.weights[i] : 1} 1 ${props.values.choices.length / 300}px`,
      overflow: 'hidden',
      textOverflow: 'ellipsis',
      minWidth: 0,
    }
    return (
      <span className="value" style={style}>{choice}</span>
    );
  });
  return (
    <div className="value-container">
      {values}
    </div>
  );
}
