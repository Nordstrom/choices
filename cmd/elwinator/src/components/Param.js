import React from 'react';
import { Link, browserHistory } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import NewChoice from './NewChoice';
import ChoiceList from './ChoiceList';
import { namespaceURL, experimentURL, choiceNewURL } from '../urls';
import { toggleWeighted, clearChoices, paramDelete } from '../actions';
import { getExperiments } from '../reducers/experiments';
import { getParams } from '../reducers/params'; 

const Param = ({ namespaceName, experimentID, experimentName, p, dispatch }) => {
  return (
    <div className="container">
      <div className="row"><div className="col-sm-9 col-sm-offset-3"><h1>{p.name}</h1></div></div>
      <div className="row">
        <NavSection>
          <Link
            to={namespaceURL(namespaceName)}
            className="nav-link"
          >{namespaceName} - Namespace</Link>
          <Link
            to={experimentURL(namespaceName, experimentName)}
            className="nav-link"
          >{experimentName} - Experiment</Link>
          <Link
            to={choiceNewURL(namespaceName, experimentName, p.name)}
            className="nav-link"
          >Create a new choice</Link>
        </NavSection>
        <div className="col-sm-9">
          <h2>Weight</h2>
          <form>
            <div className="form-group">
              <label>
                <input type="checkbox"
                  onChange={() => {
                    dispatch(toggleWeighted(p.id));
                    dispatch(clearChoices(p.id));
                  }}
                  checked={p.isWeighted} /> Weighted choices
              </label>
	      <p className="help-block">
                If you select this then you will need to add weights to your choices.
              </p>
            </div>
          </form>
          <h2>Choices</h2>
          <ChoiceList
            namespaceName={namespaceName}
            experimentName={experimentName}
            paramName={p.name}
            choices={p.choices}
            weights={p.weights}
          />
          <NewChoice
            namespaceName={namespaceName}
            experimentName={experimentName}
            paramName={p.name}
            isWeighted={p.isWeighted}
            redirectOnSubmit={false}
            dispatch={dispatch}
          />
          <button className="btn btn-warning" onClick={() => {
            dispatch(paramDelete(experimentID, p.id));
            browserHistory.push(experimentURL(namespaceName, experimentName));
          }}>Delete param {p.name}</button>
        </div>
      </div>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  const ns = state.entities.namespaces.find(n => n.name === ownProps.params.namespace);
  const exp = getExperiments(state.entities.experiments, ns.experiments).find(e => e.name === ownProps.params.experiment);
  const p = getParams(state.entities.params, exp.params).find(p => p.name === ownProps.params.param);
  return {
    namespaceName: ns.name,
    experimentID: exp.id,
    experimentName: exp.name,
    p,
  }
}

const connected = connect(mapStateToProps)(Param);

export default connected;
