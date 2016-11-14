import { connect } from 'react-redux';

import NewLabel from '../components/NewLabel';
import { addLabel } from '../actions';

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  return {
    namespaceName: ns.name,
    redirectOnSubmit: true,
  }
}

const mapDispatchToProps = {
  addLabel,
}

const connected = connect(mapStateToProps, mapDispatchToProps)(NewLabel);

export default connected;