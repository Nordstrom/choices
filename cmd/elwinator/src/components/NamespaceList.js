import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import { namespaceURL } from '../urls';

const NamespaceList = ({ namespaces }) => {
  const namespaceList = namespaces.map(n => <li key={n.name}><Link to={namespaceURL(n.name)}>{n.name}</Link></li>)
  return (
    <div className="row">
      <div className="col-sm-12">
        <ul className="list-unstyled">
          {namespaceList}
        </ul>
      </div>
    </div>
  );
};

const mapStateToProps = (state) => ({
  namespaces: state.namespaces,
});

const connected = connect(mapStateToProps)(NamespaceList);

export default connected;
