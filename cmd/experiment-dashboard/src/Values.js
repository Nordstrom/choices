import React from 'react';

export const Values = props => {
  if (!props.values) {
    return false;
  }
  const values = props.values.choices.map((choice, i) => {
    return (
      <div className="value" key={choice}>
        <h2 className="value-name">{choice}</h2>
        <span className="value-weight">{props.values.weights ? props.values.weights[i] : 1}</span>
      </div>
    );
  });
  return (
    <div>
      {values}
    </div>
  );
}
