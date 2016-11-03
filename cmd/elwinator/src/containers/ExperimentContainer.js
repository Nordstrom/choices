import { connect } from 'react-redux';

import { updateName, createExperiment } from '../actions';
import Experiment from '../components/Experiment';

const mapStateToProps = (state) => ({
  edit: state.experiment.edit,
  experiments: state.experiment.experiments,
});

const mapDispatchToProps = ({
  updateName,
  createExperiment,
});

const ExperimentContainer = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Experiment)

export default ExperimentContainer;
