import React from 'react';

const NewLabel = (props) => {
  return (
    <div className="container">
      <form>
        <div className="form-group">
          <label>Label</label>
          <input type="text" className="form-control" />
        </div>
        <button type="submit" className="btn btn-primary" />
      </form>
    </div>
  );
}

export default NewLabel;
