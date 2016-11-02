import { connect } from 'react-redux';

import { removeParam } from '../actions';
import ParamList from '../ParamList';

const mapStateToProps = (state) => ({
    params: state,
});

const mapDispatchToProps = ({
    onParamClick: removeParam,
});

const Params = connect(
    mapStateToProps,
    mapDispatchToProps,
)(ParamList)

export default Params;
