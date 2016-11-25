import React from 'react';
import { Link } from 'react-router';

import NavSection from './NavSection';
import { namespaceURL } from '../urls';
import NewLabel from '../connectors/NewLabel';

const NewLabelView = ({ params }) => {
  return (
    <div className="container">
      <div className="row"><h1>Create a new experiment</h1></div>
      <div className="row">
        <NavSection>
          <Link to={ namespaceURL(params.namespace) } className="nav-link">{params.namespace} - Namespace</Link>
        </NavSection>
        <div className="col-sm-9">
          <NewLabel params={params}/>
        </div>
      </div>
    </div>
  );
}

export default NewLabelView;
