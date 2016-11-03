import React from 'react';

const ParamList = ({ params, onParamClick }) => (
  <ul>
    {params.map(param => 
      <li key={param.id} onClick={()=> onParamClick(param.id)}>
        {param.name}
      </li>
    )}
  </ul>
);

export default ParamList;
