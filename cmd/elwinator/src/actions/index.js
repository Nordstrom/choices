/**
 * addNamespace is an action that adds a namespace to the namespace list.
 * @param {string} name - The name of the namespace.
 */
export const addNamespace = (name) => ({
  type: 'ADD_NAMESPACE',
  name,
});

/**
 * namespaceName is an action that set the name of the namespace.
 * @param {string} namespace - The original name of the namespace.
 * @param {string} name - The new name of the namespace.
 */
export const namespaceName = (namespace, name) => ({
  type: 'NAMESPACE_NAME',
  namespace,
  name,
});

let nextLabelID = 0;

/**
 * addLabel is an action that adds a label to an experiment.
 * @param {string} namespace - The namespace for the label.
 * @param {string} name - The name of the label to add.
 */
export const addLabel = (namespace, name) => ({
  type: 'ADD_LABEL',
  namespace,
  id: nextLabelID++,
  name,
})

/** 
 * addExperiment is an action that adds and experiment to the namespace.
 * @param {string} namespace - The namespace for the experiment.
 * @param {string} name - The name of the experiment to add to the namespace.
 */
export const addExperiment = (namespace, name) => ({
  type: 'ADD_EXPERIMENT',
  namespace,
  name,
});

/**
 * experimentName is an action that sets the name in an experiment.
 * @param {string} namespace - The namespace for the experiment.
 * @param {string} experiment - The experiment's original name.
 * @param {string} name - The experiment's new name.
 */
export const experimentName = (namespace, experiment, name) => ({
  type: 'EXPERIMENT_NAME',
  namespace,
  experiment,
  name,
});

/**
 * paramName is an action that sets the param name in an experiments param.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment's name.
 * @param {string} param - The param's original name.
 * @param {string} name - The param's new name.
 */
export const paramName = (namespace, experiment, param, name) => ({
  type: 'PARAM_NAME',
  namespace,
  experiment,
  param,
  name,
});

/**
 * addParam is an action that adds a param to an experiment.
 * @param {string} experiment - The experiment name.
 * @param {Object} param - The param you are adding.
 */
export const addParam = (namespace, experiment, name) => ({
  type: 'ADD_PARAM',
  namespace,
  experiment,
  name,
});

/**
 * toggleWeighted is an action that toggles whether a param is weighted or
 * uniform.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment's name.
 * @param {string} param - The name of the param.
 */
export const toggleWeighted = (namespace, experiment, param) => ({
  type: 'TOGGLE_WEIGHTED',
  namespace,
  experiment,
  param,
});

/**
 * addChoice is an action that adds a choice to a param. You must also call
 * addWeight if the param is a weighted param.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment's name.
 * @param {string} param - The name of the param.
 * @param {string} choice - The choice to add to the param.
 */
export const addChoice = (namespace, experiment, param, choice) => ({
  type: 'ADD_CHOICE',
  namespace,
  experiment,
  param,
  choice,
});

/**
 * addWeight is an action that adds a weight to a param. You should always
 * call addChoice before calling this.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment's name.
 * @param {string} param - The name of the param.
 * @param {string} weight - The weight to add to the param.
 */
export const addWeight = (namespace, experiment, param, weight) => ({
  type: 'ADD_WEIGHT',
  namespace,
  experiment,
  param,
  weight,
});