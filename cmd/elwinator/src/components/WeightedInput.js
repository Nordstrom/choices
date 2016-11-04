import React from 'react';
import { connect } from 'react-redux';

import { toggleWeighted } from '../actions';

const WeightedInput = ({ isWeighted, onWeightedClick }) => {
  return (
    <label>
      Use Weighted Params
      <input type="checkbox" checked={isWeighted} onChange={() => onWeightedClick()} />
    </label>
  );
}

const mapStateToProps = (state) => ({
  isWeighted: state.params.isWeighted,
});

const mapDispatchToProps = ({
  onWeightedClick: toggleWeighted,
});

const wiContainer = connect(
  mapStateToProps,
  mapDispatchToProps,
)(WeightedInput);

export default wiContainer;
