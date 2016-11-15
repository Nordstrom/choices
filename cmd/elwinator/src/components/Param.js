import React from 'react';
import { Link } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import NewChoice from './NewChoice';
import ChoiceList from './ChoiceList';
import { namespaceURL, experimentURL, choiceNewURL } from '../urls';
import { addChoice, addWeight, toggleWeighted, clearChoices } from '../actions';

const Param = ({ namespaceName, experimentName, p, addChoice, addWeight, toggleWeighted, clearChoices }) => {
  return (
    <div className="container">
      <div className="row"><div className="col-sm-9 col-sm-offset-3"><h1>{p.name}</h1></div></div>
      <div className="row">
        <NavSection>
          <Link to={namespaceURL(namespaceName)} className="nav-link">{namespaceName} - Namespace</Link>
          <Link to={experimentURL(namespaceName, experimentName)} className="nav-link">{experimentName} - Experiment</Link>
          <Link to={choiceNewURL(namespaceName, experimentName, p.name)} className="nav-link">Create a new choice</Link>
        </NavSection>
        <div className="col-sm-9">
          <h2>Weight</h2>
          <form>
            <div className="form-group">
              <label>
                <input type="checkbox"
                  onChange={() => {
                    toggleWeighted(namespaceName, experimentName, p.name);
                    clearChoices(namespaceName, experimentName, p.name);
                  }}
                  checked={p.isWeighted} /> Weighted choices
              </label>
              <p className="help-block">If you select this then you will need to add weights to your choices.</p>
            </div>
          </form>
          <h2>Choices</h2>
          <ChoiceList choices={p.choices} />
          <NewChoice
            namespaceName={namespaceName}
            experimentName={experimentName}
            paramName={p.name}
            isWeighted={p.isWeighted}
            addChoice={addChoice}
            addWeight={addWeight}
            redirectOnSubmit={false}
          />
        </div>
      </div>
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
  toggleWeighted,
  clearChoices,
}

const connected = connect(mapStateToProps, mapDispatchToProps)(Param);

export default connected;
