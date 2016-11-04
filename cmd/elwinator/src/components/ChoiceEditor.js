import React from 'react';
import { connect } from 'react-redux';

import { updateChoice, updateWeight, addChoice } from '../actions';

const ChoiceEditor = ({ choice, weight, isWeighted, updateChoice, updateWeight, addChoice }) => {
  return (
    <div className="form-group">
      <label>Choices</label>
      <div>
        <input
          type="text"
          value={choice}
          placeholder="value"
          onChange={(e) => updateChoice(e.target.value)}
        />
        <input
          type="number"
          value={weight}
          min="1"
          max="100"
          disabled={!isWeighted}
          placeholder="weight"
          onChange={(e) => updateWeight(e.target.value)}
          />
        <button onClick={() => addChoice()}>Add Choice</button>
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
