import { getLabels } from './reducers/labels';
import { getExperiments } from './reducers/experiments';
import { getParams } from './reducers/params';

/**
 * isClaimed returns true when the byte at bit i is 1
 * @param {number} byte - one byte of the segments
 * @param {number} i - the bit to check
 */
export const isClaimed = (byte, i) => {
  let bit = 1 << i;
  return (byte & bit) === bit;
}

export const claimSegment = (arr, i) => {
  const index = Math.floor(i/8);
  const bit = i%8;
  const current = arr[index];
  arr[index] = current | (1 << bit);
}

const ones = [
  0,1,1,2,1,2,2,3,1,2,2,3,2,3,3,4,
  1,2,2,3,2,3,3,4,2,3,3,4,3,4,4,5,
  1,2,2,3,2,3,3,4,2,3,3,4,3,4,4,5,
  2,3,3,4,3,4,4,5,3,4,4,5,4,5,5,6,
  1,2,2,3,2,3,3,4,2,3,3,4,3,4,4,5,
  2,3,3,4,3,4,4,5,3,4,4,5,4,5,5,6,
  2,3,3,4,3,4,4,5,3,4,4,5,4,5,5,6,
  3,4,4,5,4,5,5,6,4,5,5,6,5,6,6,7,
  1,2,2,3,2,3,3,4,2,3,3,4,3,4,4,5,
  2,3,3,4,3,4,4,5,3,4,4,5,4,5,5,6,
  2,3,3,4,3,4,4,5,3,4,4,5,4,5,5,6,
  3,4,4,5,4,5,5,6,4,5,5,6,5,6,6,7,
  2,3,3,4,3,4,4,5,3,4,4,5,4,5,5,6,
  3,4,4,5,4,5,5,6,4,5,5,6,5,6,6,7,
  3,4,4,5,4,5,5,6,4,5,5,6,5,6,6,7,
  4,5,5,6,5,6,6,7,5,6,6,7,6,7,7,8,
];

const toSegment = seg => {
  if (seg === '') {
    return new Array(128).fill(0);
  }
  const segments = new Uint8Array(atob(seg).split("").map(c => c.charCodeAt(0)));
  return segments;
}

const toParam = p => ({
  id: p.id,
  experiment: p.experiment,
  name: p.name || '',
  choices: p.value.choices || [],
  weights: p.value.weights || [],
  isWeighted: p.value.weights ? p.value.weights.length > 0 : false,
});

const toExperiment = exp => {
  const name = exp.name || '';
  let segments = exp.segments || '';
  segments = toSegment(segments);
  const numSegments = segments.reduce((p, v) => p + ones[v], 0);

  let params = exp.params || [];
  params = params.map(p => p.id);
  return { id: exp.id, namespace: exp.namespace, name, segments, numSegments, params };
}

const toLabel = label => ({
  id: label.id,
  name: label.name || '',
});

export const toNamespace = ns => {
  let name = ns.name || '';
  let labels = ns.labels || [];
  labels = labels.map(l => l.id);
  let experiments = ns.experiments || [];
  experiments = experiments.map(e => e.id);
  return { name, labels, experiments }
}

const flatMap = (array, lambda) => {
  return [].concat.apply([], array.map(lambda));
}

export const toEntities = namespaces => {
  if (!namespaces) {
    return undefined;
  }
  const ns = namespaces.map(n => {
    let { labels = [], experiments = [] } = n;
    return {
      ...n,
      labels: labels.map(l => ({
        id: `${n.name}-${l}`,
        name: l,
      })),
      experiments: experiments.map(e => {
        let { params = [] } = e
        return {
          ...e,
          id: `${n.name}-${e.name}`,
          namespace: n.name,
          params: params.map(p => ({
            ...p,
            id: `${n.name}-${e.name}-${p.name}`,
            experiment: `${n.name}-${e.name}`,
          })),
        };
      }),
    };
  });
  return {
    namespaces: ns.map(n => toNamespace(n)),
    labels: flatMap(ns, n => n.labels).map(l => toLabel(l)),
    experiments: flatMap(ns, n => n.experiments).map(e => toExperiment(e)),
    params: flatMap(ns, n => flatMap(n.experiments, e => e.params)).map(p => toParam(p)),
  };
};

const fromSegments = (segments) => {
  return btoa(String.fromCharCode.apply(null, segments));
};

const fromParam = (param) => {
  const p = {
    name: param.name,
    value: param.isWeighted ? {
      choices: param.choices,
      weights: param.weights,
    } : { choices: param.choices },
  };
  return p;
}

const fromExperiment = (state, experiment) => {
  const e = {
    name: experiment.name,
    segments: fromSegments(experiment.segments),
    params: getParams(state.params, experiment.params).map(p => fromParam(p)),
  };
  return e;
}

const fromLabels = (labels) => {
  return labels.map(l => l.name);
}

export const fromNamespace = (state, namespace) => {
  const n = {
    name: namespace.name,
    labels: fromLabels(getLabels(state.labels, namespace.labels)),
    experiments: getExperiments(state.experiments, namespace.experiments).map(e => fromExperiment(state, e)),
  }
  return n;
}