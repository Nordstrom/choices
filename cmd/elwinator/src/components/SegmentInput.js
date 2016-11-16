import React, { PropTypes } from 'react';
import { browserHistory } from 'react-router';

import { experimentURL } from '../urls';

const percents = [1, 25, 50, 100];

const SegmentInput = ({ namespaceName, experimentName, namespaceSegments, numSegments, availableSegments, redirectOnSubmit, experimentNumSegments }) => {
  let numSeg;

  const radio = percents.map(p => {
    return (      
      <div className="radio" key={p}>
        <label>
        <input type="radio"
          name="percent"
          checked={ Math.floor((p/100)*availableSegments) === numSegments }
          onChange={() => experimentNumSegments(namespaceName, experimentName, namespaceSegments,  Math.floor((p/100)*availableSegments))}
        /> {p}% of available segments
        </label>
      </div>
    );
  });
  return (
    <form onSubmit={e => {
      e.preventDefault();
      if (!numSeg.value.trim()) {
        return;
      }
      experimentNumSegments(namespaceName, experimentName, namespaceSegments, numSeg.value);
      if (!redirectOnSubmit) {
        return;
      }
      browserHistory.push(experimentURL(namespaceName, experimentName));
    }}>
    Segments available: <strong>{availableSegments}</strong>
    {radio}
    <div className="form-group">
      <label>Number of segments</label>
      <input
        type="number"
        min="1"
        max={availableSegments}
        className="form-control"
        value={numSegments}
        onChange={(e) => experimentNumSegments(namespaceName, experimentName, namespaceSegments, e.target.value)}
        ref={ node => numSeg = node }
      />
      <p className="help-block">The number of segments to use for this experiment</p>
    </div>
    <button type="submit" className="btn btn-default">Randomize Segments</button>
    </form>
  );
};

SegmentInput.propTypes = {
  namespaceName: PropTypes.string.isRequired,
  experimentName: PropTypes.string.isRequired,
  namespaceSegments: PropTypes.array.isRequired,
  numSegments: PropTypes.number.isRequired,
  availableSegments: PropTypes.number.isRequired,
  redirectOnSubmit: PropTypes.bool.isRequired,
  experimentNumSegments: PropTypes.func.isRequired,
};

export default SegmentInput;
