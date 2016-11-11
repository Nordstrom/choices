import React from 'react';
import { Link } from 'react-router';

import NamespaceList from './components/NamespaceList';

const App = ({ namespaceName, params, children }) => {
  return (
    <div className="container">
      <h1>Namespaces</h1>
      <NamespaceList />
      <Link to="/n/new">Create new namespace</Link>
    </div>
  );
};

export default App;
