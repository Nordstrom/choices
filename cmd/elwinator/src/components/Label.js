import React from 'react';

import AppliedLabels from '../containers/AppliedLabels';
import UnappliedLabels from '../containers/UnappliedLabels';
import AddLabel from '../containers/AddLabel';

const Label = (props) => {
  return(
    <div>
      <label>Labels: </label>
      <AppliedLabels />
      <UnappliedLabels />
      <AddLabel />
    </div>
  );
}

export default Label;
