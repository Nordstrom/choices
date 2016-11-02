import React from 'react';
import { connect } from 'react-redux';
import { addParam } from '../actions';

let AddParam = ({ dispatch }) => {
  let input;

  return (
    <div>
      <form onSubmit={e => {
        e.preventDefault();
        if (!input.value.trim()) {
          return;
        }
        dispatch(addParam(input.value));
        input.value = '';
      }}>
        <input ref={node => {
          input = node;
        }} />
        <button type="submit">
          Add Param
        </button>
      </form>
    </div>
  );
}
AddParam = connect()(AddParam);

export default AddParam;