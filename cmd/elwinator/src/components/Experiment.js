import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

const Experiment = ({ namespaceName, experimentName }) => {
  return (
    <div>
    <Link to={`/namespace/${namespaceName}/experiment/${experimentName}/param/new`}>Create param</Link>
    </div>
  );
}

const mapStateToProps = (state) => ({
});

const connected = connect(
  mapStateToProps,
)(Experiment);

export default connected;
