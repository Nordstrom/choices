import React from 'react';
import { connect } from 'react-redux';

import NewParamSection from './NewParamSection';
import { updateName, createExperiment } from '../actions';

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
      <NewParamSection />
      <button onClick={() => createExperiment(edit)}>Create Experiment</button>
    </div>
  );
}

const mapStateToProps = (state) => ({
  edit: state.experiments.edit,
  experiments: state.experiments.experiments,
});

const mapDispatchToProps = ({
  updateName,
  createExperiment,
});

const ExperimentContainer = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Experiment);

export default ExperimentContainer;
