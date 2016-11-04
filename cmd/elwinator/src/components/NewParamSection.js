import React from 'react';

import ParamList from './ParamList';
import ParamInput from './ParamInput';

const NewParamSection = (props) => {
  return (
    <div>
      <h2>Pararms</h2>
      <ParamList />
      <ParamInput />
    </div>
  );
}

export default NewParamSection;
