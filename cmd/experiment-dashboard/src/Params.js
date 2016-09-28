import React from 'react';
import { Values } from './Values.js';

export const Params = props => {
  const params = props.params.map(param => {
    return (
      <div className="param" key={param.name}>
        <h2 className="param-name">{param.name}</h2>
        <Values values={param.values} />
      </div>
    );
  });
  return (
    <div>
      {params}
    </div>
  );
}
