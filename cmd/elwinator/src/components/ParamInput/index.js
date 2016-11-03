import React from 'react';

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

export default ParamInput;