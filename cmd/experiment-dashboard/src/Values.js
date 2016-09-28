import React from 'react';

export const Values = props => {
  const values = props.values.map(value => {
    return (
      <div className="value" key={value.name}>
        <h2 className="value-name">{value.name}</h2>
        <span className="value-weight">{value.weight}</span>
      </div>
    );
  });
  return (
    <div>
      {values}
    </div>
  );
}