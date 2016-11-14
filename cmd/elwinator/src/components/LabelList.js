import React, { PropTypes } from 'react';
import classNames from 'classnames';

const LabelList = ({ namespaceName, labels, toggleLabel }) => {
  const labelList = labels.map(l => {
    const spanClassName = classNames({
      'btn': true,
      'btn-default': !l.active,
      'btn-success': l.active,
    });
    return (
    <li key={l.name}>
      <button className={spanClassName} onClick={() => toggleLabel(namespaceName, l.name)}>{l.name}</button>
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
  toggleLabel: PropTypes.func.isRequired,
};

export default LabelList;
