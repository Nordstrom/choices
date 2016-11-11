import React from 'react';
import { connect } from 'react-redux';

import { updateChoice, updateWeight, addChoice } from '../actions';

const ChoiceEditor = ({ choice, weight, isWeighted, updateChoice, updateWeight, addChoice }) => {
  return (
    <div>
    <div className="form-group">
      <label className="col-sm-2 control-label">Choice</label>
      <div className="col-sm-10">
        <input
          type="text"
          value={choice}
          placeholder="value"
          onChange={(e) => updateChoice(e.target.value)}
        />
      </div>
    </div>
    <div className="form-group">
      <label className="col-sm-2 control-label">Weight</label>
      <div className="col-sm-10">
        <input
          type="number"
          value={weight}
          min="1"
          max="100"
          disabled={!isWeighted}
          placeholder="weight"
          onChange={(e) => updateWeight(e.target.value)}
        />
      </div>
    </div>
    <div className="form-group">
      <div className="col-sm-offset-2 col-sm-10">
        <button onClick={() => addChoice()}>Add Choice</button>
      </div>
    </div>
    </div>
  );
}

const mapStateToProps = (state) => ({
  choice: state.choices.choice,
  weight: state.choices.weight,
  isWeighted: state.choices.isWeighted,
});

const mapDispatchToProps = {
  updateChoice,
  updateWeight,
  addChoice,
}

const ChoiceEditorContainer = connect(
  mapStateToProps,
  mapDispatchToProps, 
)(ChoiceEditor);

export default ChoiceEditorContainer;
