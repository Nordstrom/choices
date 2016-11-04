import React from 'react';
import { connect } from 'react-redux';

import { createLabel } from '../actions';

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

const mapStateToProps = (state) => ({});

const mapDispatchToProps = ({
  createLabel,
});

const AddLabel = connect(
  mapStateToProps,
  mapDispatchToProps,
)(LabelInput)

export default AddLabel;
