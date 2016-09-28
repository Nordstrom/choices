import React from 'react';
import { Segments } from './Segments.js';
import { Params } from './Params.js';
import './Experiments.css';

export const Experiments = props => {
  const experiments = props.experiments.map(experiment => {
    return (
      <div className="experiment" key={experiment.name}>
        <h2 className="experiment-name">{experiment.name}</h2>
        <Segments segments={experiment.segments} />
        <Params params={experiment.params} />
      </div>
    );
  });
  return (
    <div>
      {experiments}
    </div>
  )
}
