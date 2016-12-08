/**
 * rootURL is the root url.
 */
export const rootURL = () => "/";

/**
 * namespaceNewURL is the url for creating new namespaces.
 */
export const namespaceNewURL = () => "/new-namespace";

/**
 * namespaceURL is the url for a specified namespace.
 * @param {string} n - The namespace id.
 */
export const namespaceURL = (n) => `/n/${encodeURIComponent(n)}`;

/**
 * labelNewURL is the url for creating labels in a namespace.
 * @param {string} n - The namespace id.
 */
export const labelNewURL = (n) => `/n/${encodeURIComponent(n)}/new-label`;

/**
 * experimentNewURL is the url for creating experiments in a namespace.
 * @param {string} n - The namespace id.
 */

export const experimentNewURL = (n) => `/n/${encodeURIComponent(n)}/new-experiment`;

/**
 * experimentURL is the url for a specified experiment.
 * @param {string} e - The experiment id.
 */
export const experimentURL = (e) => `/e/${encodeURIComponent(e)}`;

/**
 * paramNewURL is the url for creating params in an experiment.
 * @param {string} e - The experiment id.
 */
export const paramNewURL = (e) =>
  `/e/${encodeURIComponent(e)}/new-param`;

/**
 * paramURL is the url for a specified param.
 * @param {string} p - The param id.
 */
export const paramURL = (p) =>
  `/p/${encodeURIComponent(p)}`;

/**
 * choiceNewURL is the url for creating a new choice.
 * @param {string} p - The param id.
 */
export const choiceNewURL = (p) =>
  `/p/${encodeURIComponent(p)}/new-choice`;
