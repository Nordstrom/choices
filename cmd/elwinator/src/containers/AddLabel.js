import { connect } from 'react-redux';

import { createLabel } from '../actions';
import LabelInput from '../components/LabelInput';

const mapStateToProps = (state) => ({});

const mapDispatchToProps = ({
  createLabel,
});

const AddLabel = connect(
  mapStateToProps,
  mapDispatchToProps,
)(LabelInput)

export default AddLabel;
