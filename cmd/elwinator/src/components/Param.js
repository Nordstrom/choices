import React, { PropTypes } from 'react';
import { Link, browserHistory } from 'react-router';
import { connect } from 'react-redux';

import NavSection from './NavSection';
import NewChoice from './NewChoice';
import ChoiceList from './ChoiceList';
import { namespaceURL, experimentURL, choiceNewURL } from '../urls';
import { paramToggleWeighted, paramClearChoices, paramDelete } from '../actions';
import { getNamespace } from '../reducers/namespaces';
import { getExperiment } from '../reducers/experiments';
import { getParam } from '../reducers/params';

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
            to={experimentURL(experimentID)}
            className="nav-link"
          >{experimentName} - Experiment</Link>
          <Link
            to={choiceNewURL(p.id)}
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
                    dispatch(paramToggleWeighted(p.id));
                    dispatch(paramClearChoices(p.id));
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
            paramID={p.id}
            choices={p.choices}
            weights={p.weights}
          />
          <NewChoice
            paramID={p.id}
            paramName={p.name}
            isWeighted={p.isWeighted}
            redirectOnSubmit={false}
            dispatch={dispatch}
          />
          <button className="btn btn-warning" onClick={() => {
            dispatch(paramDelete(experimentID, p.id));
            browserHistory.push(experimentURL(experimentID));
          }}>Delete param {p.name}</button>
        </div>
      </div>
    </div>
  );
}

Param.propTypes = {
  namespaceName: PropTypes.string.isRequired,
  experimentID: PropTypes.string.isRequired,
  experimentName: PropTypes.string.isRequired,
}

const mapStateToProps = (state, ownProps) => {
  const p = getParam(state.entities.params, ownProps.params.param);
  const exp = getExperiment(state.entities.experiments, p.experiment);
  const ns = getNamespace(state.entities.namespaces, exp.namespace);
  return {
    namespaceName: ns.name,
    experimentID: exp.id,
    experimentName: exp.name,
    p,
  }
}

const connected = connect(mapStateToProps)(Param);

export default connected;
