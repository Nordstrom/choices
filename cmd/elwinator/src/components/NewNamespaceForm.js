import React from 'react';
import { connect } from 'react-redux';

import * as actions from '../actions';

const NewNamespaceForm = ({
  label,
  createLabel,
  updateLabel,
  toggleLabel,
  experimentName,
  updateExperimentName,
  createExperiment,
  param,
  params,
  updateParam,
  addParam,
  choice,
  choices,
  updateChoice,
  weight,
  weights,
  updateWeight,
  addChoice,
  isWeighted,
  toggleWeighted,
  createNamespace,
}) => {
return (
  <div className="container">
  <div className="row">
    <form className="form-horizontal">
    <div className="form-group">
      <label className="col-sm-2 control-label">Labels</label>
      <div className="col-sm-10">
        <input className="form-control"
          type="text"
          value={label}
          placeholder="New label"
          onChange={(e) => {updateLabel(e.target.value)}}
        />
      </div>
    </div>
    <div className="form-group">
      <div className="col-sm-offset-2 col-sm-10">
        <button className="btn btn-default" onClick={(e) => {e.preventDefault(); createLabel(label)}}>Create New Label</button>
      </div>
    </div>
    <div className="form-group">
      <label className="col-sm-2 control-label">Experiment Name</label>
      <div className="col-sm-10">
        <input className="form-control"
          type="text"
          value={experimentName}
          placeholder="Experiment name"
          onChange={(e) => updateExperimentName(e.target.value)}
        />
      </div>
    </div>
    <div className="form-group">
      <label className="col-sm-2 control-label">Param</label>
      <div className="col-sm-10">
        <input className="form-control"
          type="text"
          value={param}
          placeholder="Param"
          onChange={(e) => updateParam(e.target.value)}
        />
      </div>
    </div>
    <div className="form-group">
      <div className="col-sm-offset-2 col-sm-10">
        <label>
          <input type="checkbox" checked={isWeighted} onChange={() => toggleWeighted()} /> Weighted Param
        </label>
      </div>
    </div>
    {/* choices */}
    <div className="form-group">
      <label className="col-sm-2 control-label">Choice</label>
      <div className="col-sm-10">
        <input className="form-control"
          type="text"
          value={choice}
          placeholder="Value"
          onChange={(e) => updateChoice(e.target.value)}
        />
      </div>
    </div>
    <div className="form-group">
      <label className="col-sm-2 control-label">Weight</label>
      <div className="col-sm-10">
        <input className="form-control"
          type="number"
          value={weight}
          min="1"
          max="100"
          disabled={!isWeighted}
          placeholder="Weight"
          onChange={(e) => updateWeight(e.target.value)}
        />
      </div>
    </div>
    <div className="form-group">
      <div className="col-sm-offset-2 col-sm-10">
        <button className="btn btn-default" onClick={(e) => {e.preventDefault(); addChoice(choice, weight)}}>Add Choice</button>
      </div>
    </div>
    <div className="form-group">
      <div className="col-sm-offset-2 col-sm-10">
        <button className="btn btn-default" onClick={(e) => {e.preventDefault(); addParam(param, isWeighted, choices, weights)}}>Add Param</button>
      </div>
    </div>
    <div className="form-group">
      <div className="col-sm-offset-2 col-sm-10">
        <button className="btn btn-default" onClick={(e) => {e.preventDefault(); createExperiment(experimentName, params)}}>Add Experiment</button>
      </div>
    </div>
    <div className="form-group">
      <div className="col-sm-offset-2 col-sm-10">
        <button className="btn btn-default" onClick={(e) => {e.preventDefault(); createNamespace()}}>Create Namespace</button>
      </div>
    </div>
    </form>
  </div>
  </div>
);
}

const mapStateToProps = (state) => {
  return {
    label: state.labels.name,
    experimentName: state.experiments.edit.name,
    param: state.params.name,
    params: state.params.params,
    choice: state.choices.choice,
    choices: state.choices.choices,
    weight: state.choices.weight,
    weights: state.choices.weights,
    isWeighted: state.choices.isWeighted,  
  }
}

const mapDispatchToProps = {
  ...actions,
  updateExperimentName: actions.updateName,
}

const nnfContainer = connect(
  mapStateToProps,
  mapDispatchToProps,
)(NewNamespaceForm)

export default nnfContainer;
