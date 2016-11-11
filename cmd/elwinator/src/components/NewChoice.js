import React from 'react';

const NewChoice = (props) => {
  return (
    <div className="container">
      <form>
        <div className="form-group">
          <label>Choice</label>
          <input type="text" className="form-control" />
        </div>
        <div className="form-group">
          <label>Weight</label>
          <input type="number" className="form-control" min="1" max="100" disabled/>
        </div>
        <button type="submit" className="btn btn-primary">Create choice</button>
      </form>
    </div>
  );
}

export default NewChoice;
