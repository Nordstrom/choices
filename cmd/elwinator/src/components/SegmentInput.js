import React, { PropTypes } from 'react';
import { browserHistory } from 'react-router';

import { experimentURL } from '../urls';

const percents = ['1', '25', '50', '100'];

const SegmentInput = ({ namespaceName, experimentName, numSegments, percent, redirectOnSubmit, experimentNumSegments, experimentPercent }) => {
  let numSeg;

  const radio = percents.map(p => {
    return (      
      <div className="radio" key={p}>
        <label>
        <input type="radio"
          name="percent"
          checked={ p === "" + Math.floor((numSegments/128) * 100) }
          onChange={() => experimentPercent(namespaceName, experimentName, p)}
        /> {p}%
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
      experimentNumSegments(namespaceName, experimentName, numSeg.value);
      if (!redirectOnSubmit) {
        return;
      }
      browserHistory.push(experimentURL(namespaceName, experimentName));
    }}>
    {radio}
    <div className="form-group">
      <label>Number of segments</label>
      <input
        type="number"
        min="1"
        max="128"
        className="form-control"
        value={numSegments}
        onChange={(e) => experimentNumSegments(namespaceName, experimentName, e.target.value)}
        ref={ node => numSeg = node }
      />
      <p className="help-block">The number of segments to use for this experiment</p>
    </div>
    <button type="submit" className="btn btn-primary">Set and Randomize Segments</button>
    </form>
  );
};

SegmentInput.propTypes = {
  namespaceName: PropTypes.string.isRequired,
  experimentName: PropTypes.string.isRequired,
  numSegments: PropTypes.string.isRequired,
  percent: PropTypes.string,
  redirectOnSubmit: PropTypes.bool.isRequired,
  experimentNumSegments: PropTypes.func.isRequired,
  experimentPercent: PropTypes.func.isRequired,
};

export default SegmentInput;
