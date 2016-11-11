import React from 'react';
import { Link } from 'react-router';

import NamespaceList from './components/NamespaceList';
import { namespaceNewURL } from './urls';

const App = ({ namespaceName, params, children }) => {
  return (
    <div className="container">
      <h1>Namespaces</h1>
      <NamespaceList />
      <Link to={namespaceNewURL()}>Create new namespace</Link>
    </div>
  );
};

export default App;
