import React from 'react';
import { Link } from 'react-router';

import NavSection from './NavSection';
import { rootURL, namespaceURL } from '../urls';
import NewExperiment from './NewExperiment';

const NewExperimentView = ({ params }) => {
  return (
    <div className="container">
      <div className="row"><h1>Create a new experiment</h1></div>
      <div className="row">
        <NavSection>
          <Link to={ rootURL() } className="nav-link">Home</Link>
          <Link to={ namespaceURL(params.namespace) } className="nav-link">{params.namespace} - Namespace</Link>
        </NavSection>
        <div className="col-sm-9">
          <NewExperiment namespaceName={params.namespace} />
        </div>
      </div>
    </div>
  );
}

export default NewExperimentView;
