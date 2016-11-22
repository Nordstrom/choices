import { connect } from 'react-redux';

import PublishView from '../components/PublishView';

const mapStateToProps = (state) => ({
  namespaces: state.namespaces,
});

const connected = connect(mapStateToProps)(PublishView);

export default connected;
