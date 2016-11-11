import React from 'react';

const Choices = ({choices, weights}) => {
  const cumWeight = weights.reduce((prev, cur) => {
    return prev + cur;
  }, 0);
  const choicesList = choices.map((c, i) => {
    return (
      <li key={c}>{c} - {100*weights[i]/cumWeight}%
      </li>
    );
  });
  return (
    <ul>
      {choicesList}
    </ul>
  );
}

Choices.defaultProps = {
  choices: [],
  weights: [],
}

export default Choices;
