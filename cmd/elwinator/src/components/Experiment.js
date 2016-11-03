import React from 'react';

const Experiment = ({ edit, experiments, updateName, createExperiment }) => {
  const experimentsList = experiments.map(exp => <li key={exp.name}>{exp.name}</li>)
  return (
    <div>
      <h2>Experiments</h2>
      <ul>
        {experimentsList}
      </ul>
      <table>
      <tbody>
        <tr>
          <td>Name:</td>
          <td>
            <input type="text" value={edit.name} onChange={(ev) => updateName(ev.target.value)}/>
          </td>
        </tr>
      </tbody>
      </table>
      <button onClick={() => createExperiment(edit)}>Create Experiment</button>
    </div>
  )
}

export default Experiment;
