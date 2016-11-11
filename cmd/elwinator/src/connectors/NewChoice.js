import { connect } from 'react-redux';

import NewChoice from '../components/NewChoice';
import { addChoice, addWeight } from '../actions';

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  const exp = ns.experiments.find(e => e.name === ownProps.params.experiment);
  const param = exp.params.find(p => p.name === ownProps.params.param);
  return {
    namespaceName: ns.name,
    experimentName: exp.name,
    paramName: param.name,
    isWeighted: param.isWeighted,
    redirectOnSubmit: true,
  }
}

const mapDispatchToProps = {
  addChoice,
  addWeight,
}

const connected = connect(mapStateToProps, mapDispatchToProps)(NewChoice);

export default connected;