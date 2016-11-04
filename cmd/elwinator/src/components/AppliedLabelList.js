import React from 'react';
import { connect } from 'react-redux';

import { removeLabel } from '../actions';

const AppliedLabelList = ({ labels, onLabelClick }) => {
  const labelList = labels.map(label => 
    <li key={label.id} onClick={() => onLabelClick(label.id)}>{label.name}</li>
  )
  return (
    <div>
      <label>Applied Labels:</label>
      <ul>
        {labelList}
      </ul>
    </div>
  )
};

const mapStateToProps = (state) => ({
  labels: state.labels.filter(label => label.active),
});

const mapDispatchToProps = ({
  onLabelClick: removeLabel,
});

const AppliedLabels = connect(
  mapStateToProps,
  mapDispatchToProps,
)(AppliedLabelList)

export default AppliedLabels;
