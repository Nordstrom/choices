import React from 'react';
import { connect } from 'react-redux';

const NewNamespace = ({ namespaceName, params, children }) => {
  return (
    <div className="container">
      <form>
        <div className="form-group">
          <label>Namespace Name</label>
          <input type="text" className="form-control" />
        </div>
        <button type="submit" className="btn btn-primary" >Create namespace</button>
      </form>
    </div>
  );
};

const mapStateToProps = (state) => ({
  namespaceName: state.namespace.name,
});

const connected = connect(mapStateToProps)(NewNamespace);

export default connected;
