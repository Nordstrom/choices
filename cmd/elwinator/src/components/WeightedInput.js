import React from 'react';
import { connect } from 'react-redux';

import { toggleWeighted } from '../actions';

const WeightedInput = ({ isWeighted, onWeightedClick }) => {
  return (
    <label>
      <input type="checkbox" checked={isWeighted} onChange={() => onWeightedClick()} />
      Weighted Param
    </label>
  );
}

const mapStateToProps = (state) => ({
  isWeighted: state.choices.isWeighted,
});

const mapDispatchToProps = ({
  onWeightedClick: toggleWeighted,
});

const wiContainer = connect(
  mapStateToProps,
  mapDispatchToProps,
)(WeightedInput);

export default wiContainer;
