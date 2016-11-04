import React from 'react';
import { connect } from 'react-redux';

import { addParam } from '../actions';

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

const mapStateToProps = (state) => ({});

const mapDispatchToProps = ({
  addParam,
});

const AddParam = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ParamInput);

export default AddParam;
