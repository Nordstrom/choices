import React from 'react';
import { Link } from 'react-router';

const App = ({ namespaceName, params, children }) => {
  return (
    <div className="container">
      <Link to="/n/new">Create new namespace</Link>
    </div>
  );
};

export default App;
