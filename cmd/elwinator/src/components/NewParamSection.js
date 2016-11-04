import React from 'react';

import WeightedInput from './WeightedInput';
import ParamList from './ParamList';
import AddParam from '../containers/AddParam';

const NewParamSection = (props) => {
  return (
    <div>
      <h2>Pararms</h2>
      <WeightedInput />
      <ParamList />
      <AddParam />
    </div>
  );
}

export default NewParamSection;
