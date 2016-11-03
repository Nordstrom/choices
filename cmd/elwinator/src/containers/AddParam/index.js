import { connect } from 'react-redux';

import { addParam } from '../../actions';
import ParamInput from '../../components/ParamInput';

const mapStateToProps = (state) => ({});

const mapDispatchToProps = ({
  addParam,
});

const AddParam = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ParamInput)

export default AddParam;
