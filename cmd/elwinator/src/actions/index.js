export const namespacesLoaded = (namespaces) => ({
  type: 'NAMESPACES_LOADED',
  namespaces,
});

/**
 * addNamespace is an action that adds a namespace to the namespace list.
 * @param {string} name - The name of the namespace.
 */
export const addNamespace = (name) => ({
  type: 'ADD_NAMESPACE',
  name,
});

/**
 * namespaceDelete is an action that marks a namespace to be deleted.
 * @param {string} namespace - The namespace to delete.
 */
export const namespaceDelete = (namespace) => ({
  type: 'NAMESPACE_DELETE',
  namespace,
});

/**
 * namespaceLocalDelete is an action that removes the given namespace from the
 * local state. This would be used when you create a new namespace don't
 * publish it then decide to delete it.
 * @param {string} namespace - The namespace you want to delete.
 */
export const namespaceLocalDelete = (namespace) => ({
  type: 'NAMESPACE_LOCAL_DELETE',
  namespace,
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

/**
 * toggleLabel is an action that toggles the specified label.
 * @param {string} namespace - name of the namespace the label is in.
 * @param {string} name - name of the label to toggle.
 */
export const toggleLabel = (namespace, name) => ({
  type: 'TOGGLE_LABEL',
  namespace,
  name,
})

/**
 * addLabel is an action that adds a label to an experiment.
 * @param {string} namespace - The namespace for the label.
 * @param {string} name - The name of the label to add.
 */
export const addLabel = (namespace, name) => ({
  type: 'ADD_LABEL',
  namespace,
  name,
});

/**
 * togglePublish is an action that toggles a namespace for publishing.
 * @param {string} namespace - The namespace to publish.
 */
export const togglePublish = (namespace) => ({
  type: 'TOGGLE_PUBLISH',
  namespace,
});

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
 * experimentDelete is an action that marks an experiment for deletion.
 * @param {string} namespace - The namespace the experiment is in.
 * @param {string} experiment - The name of the experiment to delete.
 */
export const experimentDelete = (namespace, experiment) => ({
  type: 'EXPERIMENT_DELETE',
  namespace,
  experiment,
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
 * experimentNumSegments is an action that sets the number of segments in an
 * expermient.
 * @param {string} namespace - The namespace that the experiment is in.
 * @param {string} experiment - The experiment that is being changed.
 * @param {Array} namespaceSegments - The segments claimed by the namespace.
 * @param {number} numSegments - The number of segments the experiment
 * should have.
 */
export const experimentNumSegments = (namespace, experiment, namespaceSegments, numSegments) => ({
  type: 'EXPERIMENT_NUM_SEGMENTS',
  namespace,
  experiment,
  namespaceSegments,
  numSegments,
});

/**
 * experimentPercent is an action that claims segments based on a percentage.
 * @param {string} namespace - The namespace that the experiment is in.
 * @param {string} experiment - The experiment that is being changed.
 * @param {string} percent - The percentage of segments to claim.
 */
export const experimentPercent = (namespace, experiment, percent) => ({
  type: 'EXPERIMENT_PERCENT',
  namespace,
  experiment,
  percent,
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
 * @param {string} namespace - The namespace that the param is in.
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
 * paramDelete is an action that deletes a param from an experiment.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment that the param is in.
 * @param {string} name - The name of the param to delete.
 */
export const paramDelete = (namespace, experiment, name) => ({
  type: 'PARAM_DELETE',
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
 * choiceDelete is an action that deletes a choice from a param.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment that the param is in.
 * @param {string} param - The param name.
 * @param {string} index - The index of the choice to delete.
 */
export const choiceDelete = (namespace, experiment, param, index) => ({
  type: 'CHOICE_DELETE',
  namespace,
  experiment,
  param,
  index,
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

/**
 * clearChoices is an action that removes all the choices and weights.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment's name.
 * @param {string} param - The name of the param.
 */
export const clearChoices = (namespace, experiment, param) => ({
  type: 'CLEAR_CHOICES',
  namespace,
  experiment,
  param,
});
