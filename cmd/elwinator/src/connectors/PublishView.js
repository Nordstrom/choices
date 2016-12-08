import { connect } from 'react-redux';

import PublishView from '../components/PublishView';
import { getNamespaces } from '../reducers/namespaces';

const mapStateToProps = (state) => ({
  namespaces: getNamespaces(state.entities.namespaces, Object.keys(state.entities.changes)),
  changes: state.entities.changes,
  entities: state.entities,
});

const connected = connect(mapStateToProps)(PublishView);

export default connected;
