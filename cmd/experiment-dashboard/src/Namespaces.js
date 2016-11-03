import React from 'react';
import { Experiments } from './Experiments.js';
import { Segments } from './Segments.js';
import './Namespaces.css';

export const Namespaces = props => {
  if (!props.namespaces) {
    return false;
  }
  const namespaces = props.namespaces.map(namespace => {
    return (
      <div className="namespace" key={namespace.name}>
        <h2 className="namespace-name">{namespace.name}</h2>
        <Segments segments={namespace.segments} />
        <Experiments experiments={namespace.experiments} />
      </div>
    );
  });
  return (
    <div className="namespace-container">
      {namespaces}
    </div>
  )
}