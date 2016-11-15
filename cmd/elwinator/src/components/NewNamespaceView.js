import React from 'react';

import NavSection from './NavSection';
import NewNamespace from './NewNamespace';

const NewNamespaceView = (props) => {
  return (
    <div className="container">
      <div className="row"><h1>Create a new namespace</h1></div>
      <div className="row">
        <NavSection />
        <div className="col-sm-9">
          <NewNamespace />
        </div>
      </div>
    </div>
  );
}

export default NewNamespaceView;
