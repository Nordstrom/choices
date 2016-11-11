import React from 'react';
import { Link } from 'react-router';

import { choiceNewURL } from '../urls';

const Param = ({ namespaceName, experimentName, paramName }) => {
  return (
    <div className="container">
      <Link to={choiceNewURL(namespaceName, experimentName, paramName)}>Create new choice</Link>
    </div>
  );
}

export default Param;
