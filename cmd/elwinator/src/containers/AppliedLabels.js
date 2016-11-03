import { connect } from 'react-redux';

import { removeLabel } from '../actions';
import AppliedLabelList from '../components/AppliedLabelList';

const mapStateToProps = (state) => ({
  labels: state.labels.filter(label => label.active),
});

const mapDispatchToProps = ({
  onLabelClick: removeLabel,
});

const AppliedLabels = connect(
  mapStateToProps,
  mapDispatchToProps,
)(AppliedLabelList)

export default AppliedLabels;
