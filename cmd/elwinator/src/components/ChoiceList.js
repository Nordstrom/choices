import React from 'react';

const ChoiceList = ({ choices }) => {
  const choiceList = choices.map(c => <li key={c}>{c}</li>);
  return (
    <ul className="list-unstyled">
      {choiceList}
    </ul>
  );
};

export default ChoiceList;
