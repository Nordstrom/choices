import { v4 } from 'node-uuid';

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
  id: v4(),
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
  id: v4(),
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
  id: v4(),
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
 * @param {string} experiment - The experiment's original name.
 * @param {string} name - The experiment's new name.
 */
export const experimentName = (experiment, name) => ({
  type: 'EXPERIMENT_NAME',
  experiment,
  name,
});

/**
 * experimentNumSegments is an action that sets the number of segments in an
 * expermient.
 * @param {string} experiment - The experiment that is being changed.
 * @param {Array} namespaceSegments - The segments claimed by the namespace.
 * @param {number} numSegments - The number of segments the experiment
 * should have.
 */
export const experimentNumSegments = (experiment, namespaceSegments, numSegments) => ({
  type: 'EXPERIMENT_NUM_SEGMENTS',
  experiment,
  namespaceSegments,
  numSegments,
});

/**
 * paramName is an action that sets the param name in an experiments param.
 * @param {string} param - The param's id.
 * @param {string} name - The param's new name.
 */
export const paramName = (param, name) => ({
  type: 'PARAM_NAME',
  param,
  name,
});

/**
 * addParam is an action that adds a param to an experiment.
 * @param {string} experiment - The experiment id.
 * @param {Object} param - The param you are adding.
 */
export const addParam = (experiment, name) => ({
  type: 'ADD_PARAM',
  id: v4(),
  experiment,
  name,
});

/**
 * paramDelete is an action that deletes a param from an experiment.
 * @param {string} namespace - The namespace that the param is in.
 * @param {string} experiment - The experiment that the param is in.
 * @param {string} name - The name of the param to delete.
 */
export const paramDelete = (experiment, param) => ({
  type: 'PARAM_DELETE',
  experiment,
  param,
});

/**
 * toggleWeighted is an action that toggles whether a param is weighted or
 * uniform.
 * @param {string} param - The name of the param.
 */
export const toggleWeighted = (param) => ({
  type: 'TOGGLE_WEIGHTED',
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
export const addChoice = (param, choice) => ({
  type: 'ADD_CHOICE',
  param,
  choice,
});

/**
 * choiceDelete is an action that deletes a choice from a param.
 * @param {string} param - The param name.
 * @param {string} index - The index of the choice to delete.
 */
export const choiceDelete = (param, index) => ({
  type: 'CHOICE_DELETE',
  param,
  index,
});

/**
 * addWeight is an action that adds a weight to a param. You should always
 * call addChoice before calling this.
 * @param {string} param - The name of the param.
 * @param {string} weight - The weight to add to the param.
 */
export const addWeight = (param, weight) => ({
  type: 'ADD_WEIGHT',
  param,
  weight,
});

/**
 * clearChoices is an action that removes all the choices and weights.
 * @param {string} param - The name of the param.
 */
export const clearChoices = (param) => ({
  type: 'CLEAR_CHOICES',
  param,
});
