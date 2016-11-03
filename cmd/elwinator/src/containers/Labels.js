import { connect } from 'react-redux';

import { removeLabel } from '../actions';
import ParamList from '../components/ParamList';

const mapStateToProps = (state) => ({
  params: state.labels.filter(label => label.active),
});

const mapDispatchToProps = ({
  onParamClick: removeLabel,
});

const Labels = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ParamList)

export default Labels;
