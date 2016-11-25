import { connect } from 'react-redux';

import NewLabel from '../components/NewLabel';

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  return {
    namespaceName: ns.name,
    redirectOnSubmit: true,
  }
}

const connected = connect(mapStateToProps)(NewLabel);

export default connected;