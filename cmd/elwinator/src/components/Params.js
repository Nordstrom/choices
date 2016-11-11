import React from 'react'; 

import Choices from './Choices';

const Params = ({params}) => {
  const paramList = params.map(p => {
    return (
      <div key={p.name}>
        <h6>{p.name}</h6>
        <Choices choices={p.choices} weights={p.weights}/>
      </div>
    );
  })
  return (
    <div>
      {paramList}
    </div>
  );
}

Params.defaultProps = {
  params: [],
}

export default Params;
