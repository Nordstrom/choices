import React from 'react';
import { Link } from 'react-router';

import NamespaceList from './components/NamespaceList';
import PublishView from './connectors/PublishView';
import { namespaceNewURL } from './urls';

const App = ({ namespaceName, params, children }) => {
  return (
    <div className="container">
      <h1>Elwinator</h1>
      <h2>Namespaces</h2>
      <NamespaceList />
      <Link to={namespaceNewURL()}>Create new namespace</Link>
      <h2>Publish Changes</h2>
      <PublishView />
    </div>
  );
};

export default App;
