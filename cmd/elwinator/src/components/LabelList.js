import React, { PropTypes } from 'react';
import { connect } from 'react-redux';
import classNames from 'classnames';

import { toggleLabel } from '../actions';

const LabelList = ({ namespaceName, labels, dispatch }) => {
  const labelList = labels.map(l => {
    const spanClassName = classNames({
      'btn': true,
      'btn-default': !l.active,
      'btn-success': l.active,
    });
    return (
    <li key={l.name}>
      <button
        className={spanClassName}
        onClick={() => dispatch(toggleLabel(namespaceName, l.name))}
       >{l.name}</button>
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

LabelList.propTypes = {
  namespaceName: PropTypes.string.isRequired,
  labels: PropTypes.arrayOf(PropTypes.object).isRequired,
};

const connected = connect()(LabelList);

export default connected;
