/**
 * namespaceName is an action that set the name of the namespace.
 * @param {string} name - The namespace's new name.
 */
export const namespaceName = (name) => ({
  type: 'NAMESPACE_NAME',
  name,
});

let nextLabelID = 0;

/**
 * addLabel is an action that adds a label to an experiment.
 * @param {string} name - The name of the label to add.
 */
export const addLabel = (name) => ({
  type: 'ADD_LABEL',
  id: nextLabelID++,
  name,
})

/** 
 * addExperiment is an action that adds and experiment to the namespace.
 * @param {Object} experiment - The experiment to add to the namespace.
 */
export const addExperiment = (experiment) => ({
  type: 'ADD_EXPERIMENT',
  experiment,
});

/**
 * experimentName is an action that sets the name in an experiment.
 * @param {string} experiment - The experiment's original name.
 * @param {string} name - The experiment's new name.
 */
export const experimentName = (experiment, name) => ({
  type: 'EXPERIMENT_NAME',
  experiment,
  name,
});

/**
 * paramName is an action that sets the param name in an experiments param.
 * @param {string} experiment - The experiment's name.
 * @param {string} param - The param's original name.
 * @param {string} name - The param's new name.
 */
export const paramName = (experiment, param, name) => ({
  type: 'PARAM_NAME',
  experiment,
  param,
  name,
});

/**
 * addParam is an action that adds a param to an experiment.
 * @param {string} experiment - The experiment name.
 * @param {Object} param - The param you are adding.
 */
export const addParam = (experiment, param) => ({
  type: 'ADD_PARAM',
  experiment,
  param,
});
