import React from 'react';
import { Link } from 'react-router';

const App = ({ namespaceName, params, children }) => {
  return (
    <div>
      <Link to="/new">Create new namespace</Link>
    </div>
  );
};

export default App;
