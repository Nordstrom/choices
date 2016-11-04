import React from 'react';
import { connect } from 'react-redux';

import { removeParam } from '../actions';

const ParamList = ({ params, onParamClick }) => (
  <ul>
    {params.map(param => 
      <li key={param.id} onClick={()=> onParamClick(param.id)}>
        {param.name}
      </li>
    )}
  </ul>
);

const mapStateToProps = (state) => ({
  params: state.params.params,
});

const mapDispatchToProps = ({
  onParamClick: removeParam,
});

const ParamListContainer = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ParamList);

export default ParamListContainer;
