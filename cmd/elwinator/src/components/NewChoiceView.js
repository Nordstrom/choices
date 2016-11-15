import React from 'react';
import { Link } from 'react-router';

import NavSection from './NavSection';
import { rootURL, namespaceURL, experimentURL, paramURL } from '../urls';
import NewChoice from '../connectors/NewChoice';

const NewChoiceView = ({ params }) => {
  return (
    <div className="container">
      <div className="row"><h1>Create a new experiment</h1></div>
      <div className="row">
        <NavSection>
          <Link to={ rootURL() } className="nav-link">Home</Link>
          <Link to={ namespaceURL(params.namespace) } className="nav-link">{params.namespace} - Namespace</Link>
          <Link to={ experimentURL(params.namespace, params.experiment) } className="nav-link">{params.experiment} - Experiment </Link>
          <Link to={ paramURL(params.namespace, params.experiment, params.param) } className="nav-link">{params.param} - Param </Link>
        </NavSection>
        <div className="col-sm-9">
          <NewChoice params={params}/>
        </div>
      </div>
    </div>
  );
}

export default NewChoiceView;
