import { connect } from 'react-redux';

import { addLabel } from '../actions';
import UnappliedLabelList from '../components/UnappliedLabelList';

const mapStateToProps = (state) => ({
  labels: state.labels.filter(label => !label.active),
});

const mapDispatchToProps = ({
  onLabelClick: addLabel,
});

const UnappliedLabels = connect(
  mapStateToProps,
  mapDispatchToProps,
)(UnappliedLabelList)

export default UnappliedLabels;
