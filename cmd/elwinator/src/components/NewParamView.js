import React from 'react';
import { Link } from 'react-router';

import NavSection from './NavSection';
import { namespaceURL, experimentURL } from '../urls';
import NewParam from './NewParam';

const NewParamView = ({ params }) => {
  return (
    <div className="container">
      <div className="row"><h1>Create a new experiment</h1></div>
      <div className="row">
        <NavSection>
          <Link
            to={ namespaceURL(params.namespace) }
            className="nav-link"
          >{params.namespace} - Namespace</Link>
          <Link
            to={ experimentURL(params.namespace, params.experiment) }
            className="nav-link"
          >{params.experiment} - Experiment </Link>
        </NavSection>
        <div className="col-sm-9">
          <NewParam namespaceName={params.namespace} experimentName={params.experiment} />
        </div>
      </div>
    </div>
  );
}

export default NewParamView;
