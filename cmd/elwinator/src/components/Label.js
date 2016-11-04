import React from 'react';

import AppliedLabels from './AppliedLabelList';
import UnappliedLabels from './UnappliedLabelList';
import LabelInput from './LabelInput';

const Label = (props) => {
  return(
    <div>
      <label>Labels: </label>
      <AppliedLabels />
      <UnappliedLabels />
      <LabelInput />
    </div>
  );
}

export default Label;
