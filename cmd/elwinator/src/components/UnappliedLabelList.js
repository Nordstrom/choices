import React from 'react';
import { connect } from 'react-redux';

import { addLabel } from '../actions';

const UnappliedLabelList = ({ labels, onLabelClick }) => {
  const labelsList = labels.map(label => 
      <li key={label.id} onClick={() => onLabelClick(label.id)}>{label.name}</li>
  )
  return (
    <div>
      <label>Unapplied labels:</label>
      <ul>
        {labelsList}
      </ul>
    </div>
  );
};

const mapStateToProps = (state) => ({
  labels: state.labels.filter(label => !label.active),
});

const mapDispatchToProps = ({
  onLabelClick: addLabel,
});

const UnappliedLabels = connect(
  mapStateToProps,
  mapDispatchToProps,
)(UnappliedLabelList)

export default UnappliedLabels;
