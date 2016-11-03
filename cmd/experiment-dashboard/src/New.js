import React from 'react';

export const New = props => {
  return (
    <div>
      <form action="/api/new" method="post">
        <label>Labels</label>
        <select name="labels" multiple>
          <option>search</option>
          <option>search-mobile</option>
          <option>rands</option>
        </select>
        <label>Experiment Name</label>
        <input name="experiment-name" type="text"/>
        <label>Number of Segments</label>
        <input name="segments" type="number" min="1" max="128"/>
        <label>Params</label>
        <input name="params" type="text"/>
      </form>
    </div>
  )
}