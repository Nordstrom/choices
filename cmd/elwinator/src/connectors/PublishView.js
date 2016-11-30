import { connect } from 'react-redux';

import PublishView from '../components/PublishView';
import { getNamespaces } from '../reducers/namespaces';

const mapStateToProps = (state) => ({
  namespaces: getNamespaces(state.entities.namespaces, Object.keys(state.changes)),
  changes: state.changes,
});

const connected = connect(mapStateToProps)(PublishView);

export default connected;
