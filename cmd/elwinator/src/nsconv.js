const toLabel = label => ({
  name: label || '',
  active: true,
});

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
  params = params.map(toParam);
  return { name, segments, numSegments, params };
}

export const toNamespace = ns => {
  let name = ns.name || '';
  let labels = ns.labels || [];
  labels = labels.map(toLabel);
  let experiments = ns.experiments || [];
  experiments = experiments.map(toExperiment);
  return { name, labels, experiments }
}

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

const fromExperiment = (experiment) => {
  const e = {
    name: experiment.name,
    segments: fromSegments(experiment.segments),
    params: experiment.params.map(p => fromParam(p)),
  };
  return e;
}

const fromLabels = (labels) => {
  return labels.filter(l => l.active)
  .map(l => l.name);
}

export const fromNamespace = (namespace) => {
  const n = {
    name: namespace.name,
    labels: fromLabels(namespace.labels),
    experiments: namespace.experiments.map(e => fromExperiment(e)),
  }
  return n;
}