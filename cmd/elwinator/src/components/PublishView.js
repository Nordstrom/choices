import React from 'react';

import { togglePublish } from '../actions';

function fromSegments(segments) {
  return btoa(String.fromCharCode.apply(null, segments));
};

function fromParam(param) {
  const p = {
    name: param.name,
    value: param.isWeighted ? {
      choices: param.choices,
      weights: param.weights,
    } : { choices: param.choices },
  };
  return p;
}

function fromExperiment(experiment) {
  const e = {
    name: experiment.name,
    segments: fromSegments(experiment.segments),
    params: experiment.params.map(p => fromParam(p)),
  };
  return e;
}

function fromLabels(labels) {
  return labels.filter(l => l.active)
  .map(l => l.name);
}

function fromNamespace(namespace) {
  const n = {
    name: namespace.name,
    labels: fromLabels(namespace.labels),
    experiments: namespace.experiments.map(e => fromExperiment(e)),
  }
  return n;
}

function createRequest(namespace) {
  return {
    method: 'POST',
    body: JSON.stringify({ namespace: fromNamespace(namespace) }),
  }
}

function updateRequest(namespace) {
  return {
    method: 'POST',
    body: JSON.stringify({ namespace: fromNamespace(namespace) }),
  }
}

const PublishView = ({ namespaces, dispatch }) => {
  const ns = namespaces.filter(n => {
    if (n.isDirty) {
      return true;
    }
    return false;
  }).map(n => <div key={n.name} className="checkbox"><label><input type="checkbox" checked={n.publish} onChange={() => dispatch(togglePublish(n.name))} /> {n.name}</label></div>);
  if (ns.length === 0) {
    return <p>No changes made</p>
  }
  return (
    <form onSubmit={e => {
      e.preventDefault();
      // filter selected
      namespaces.filter(n => n.publish)
      .forEach(n => {
        let url;
        let req;
        if (n.isNew) {
          url = '/api/v1/create';
          req = createRequest(n);
        } else {
          url = '/api/v1/update';
          req = updateRequest(n);
        }
        fetch(url, req);
      });
    }}>
    {ns}
    <button type="submit" className="btn btn-primary">Publish Changes</button>
    </form>
  );
}

PublishView.propTypes = {
  namespaces: React.PropTypes.array.isRequired,
}

PublishView.defaultProps = {
  namespaces: [],
}

export default PublishView;
