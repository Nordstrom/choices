/**
 * rootURL is the root url.
 */
export const rootURL = () => "/";

/**
 * namespaceNewURL is the url for creating new namespaces.
 */
export const namespaceNewURL = () => "/n/new";

/**
 * namespaceURL is the url for a specified namespace.
 * @param {string} n - The namespace name.
 */
export const namespaceURL = (n) => `/n/${encodeURIComponent(n)}`;

/**
 * labelNewURL is the url for creating labels in a namespace.
 * @param {string} n - The namespace name.
 */
export const labelNewURL = (n) => `/n/${encodeURIComponent(n)}/l/new`;

/**
 * experimentNewURL is the url for creating experiments in a namespace.
 * @param {string} n - The namespace name.
 */

export const experimentNewURL = (n) => `/n/${encodeURIComponent(n)}/e/new`;

/**
 * experimentURL is the url for a specified experiment.
 * @param {string} n - The namespace name.
 * @param {string} e - The experiment name.
 */
export const experimentURL = (n, e) => `/n/${encodeURIComponent(n)}/e/${encodeURIComponent(e)}`;

/**
 * paramNewURL is the url for creating params in an experiment.
 * @param {string} n - The namespace name.
 * @param {string} e - The experiment name.
 */
export const paramNewURL = (n, e) => `/n/${encodeURIComponent(n)}/e/${encodeURIComponent(e)}/p/new`;

/**
 * paramURL is the url for a specified param.
 * @param {string} n - The namespace name.
 * @param {string} e - The experiment name.
 * @param {string} p - The param name.
 */
export const paramURL = (n, e, p) => `/n/${encodeURIComponent(n)}/e/${encodeURIComponent(e)}/p/${encodeURIComponent(p)}`;

/**
 * choiceNewURL is the url for creating a new choice.
 * @param {string} n - The namespace name.
 * @param {string} e - The experiment name.
 * @param {string} p - The param name.
 */
export const choiceNewURL = (n, e, p) => `/n/${encodeURIComponent(n)}/e/${encodeURIComponent(e)}/p/${encodeURIComponent(p)}/c/new`;
