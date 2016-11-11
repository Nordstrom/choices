import React from 'react';
import { connect } from 'react-redux';
import classNames from 'classnames';

import { toggleLabel } from '../actions';

const LabelList = ({ labels, toggleLabel }) => {
  const labelList = labels.map(l => {
    const spanClassName = classNames({
      'btn': true,
      'btn-default': !l.active,
      'btn-success': l.active,
    });
    return (
    <li key={l.id}>
      <button className={spanClassName} onClick={() => toggleLabel(l.id)}>{l.name}</button>
    </li>
    );
  }
  );

  return (
    <ul className="list-inline">
      {labelList}
    </ul>
  )
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.namespaceName);
  return {
    labels: ns.labels,
  }
};

const LabelListContainer = connect(
  mapStateToProps,
  {toggleLabel},
)(LabelList);

export default LabelListContainer;
