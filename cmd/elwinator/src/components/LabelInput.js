import React from 'react';

let LabelInput = ({ createLabel }) => {
  let input;

  return (
    <div>
      <form onSubmit={e => {
        e.preventDefault();
        if (!input.value.trim()) {
          return;
        }
        createLabel(input.value);
        input.value = '';
      }}>
        <input ref={node => {
          input = node;
        }} />
        <button type="submit">
          Create Label
        </button>
      </form>
    </div>
  );
}

export default LabelInput;