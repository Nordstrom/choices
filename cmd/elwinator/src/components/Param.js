import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import ChoiceList from './ChoiceList';
import { choiceNewURL } from '../urls';
import { addChoice, addWeight } from '../actions';

const Param = ({ namespaceName, experimentName, p, addChoice, addWeight }) => {
  return (
    <div className="container">
      <h1>{p.name}</h1>
      <h2>Choices</h2>
      <ChoiceList choices={p.choices} />
      <Link to={choiceNewURL(namespaceName, experimentName, p.name)}>Create new choice</Link>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.namespaces.find(n => n.name === ownProps.params.namespace);
  const exp = ns.experiments.find(e => e.name === ownProps.params.experiment);
  const p = exp.params.find(p => p.name === ownProps.params.param);
  return {
    namespaceName: ns.name,
    experimentName: exp.name,
    p,
  }
}

const mapDispatchToProps = {
  addChoice,
  addWeight,
}

const connected = connect(mapStateToProps, mapDispatchToProps)(Param);

export default connected;
