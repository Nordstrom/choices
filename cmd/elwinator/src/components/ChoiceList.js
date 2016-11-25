import React from 'react';
import { connect } from 'react-redux';

import { choiceDelete } from '../actions';

const ChoiceList = ({ namespaceName, experimentName, paramName, choices, weights, dispatch }) => {
  const choiceList = choices.map((c, i) =>
    <tr key={c}>
      <td>{i+1}</td>
      <td>{c}</td>
      <td>{weights ? weights[i] : 1}</td>
      <td><button className="btn btn-default btn-xs" onClick={() => {
        dispatch(choiceDelete(namespaceName, experimentName, paramName, i));
      }}>&times;</button></td>
    </tr>
  );
  return (
    <table className="table table-striped">
      <thead>
        <tr>
          <th>#</th>
          <th>Choice</th>
          <th>Weight</th>
          <th>Delete</th>
        </tr>
      </thead>
      <tbody>
        {choiceList}
      </tbody>
    </table>
  );
};

ChoiceList.propTypes = {
  namespaceName: React.PropTypes.string.isRequired,
  experimentName: React.PropTypes.string.isRequired,
  paramName: React.PropTypes.string.isRequired,
  choices: React.PropTypes.arrayOf(React.PropTypes.string),
  weights: React.PropTypes.arrayOf(React.PropTypes.number),
}

const connected = connect()(ChoiceList);

export default connected;
