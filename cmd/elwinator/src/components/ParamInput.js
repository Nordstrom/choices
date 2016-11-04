import React from 'react';
import { connect } from 'react-redux';

import { addParam } from '../actions';
import WeightedInput from './WeightedInput';
import ChoiceEditor from './ChoiceEditor';

let ParamInput = ({ addParam }) => {
  let input;

  return (
    <div>
      <form onSubmit={e => {
        e.preventDefault();
        if (!input.value.trim()) {
          return;
        }
        addParam(input.value);
        input.value = '';
      }}>
        <label>Name</label>
        <div>
          <input ref={node => {
            input = node;
          }} />
        </div>
        <WeightedInput />
        <ChoiceEditor />
        {/* submit param */}
        <button type="submit">
          Add Param
        </button>
      </form>
    </div>
  );
}

const mapStateToProps = (state) => ({});

const mapDispatchToProps = ({
  addParam,
});

const AddParam = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ParamInput);

export default AddParam;
