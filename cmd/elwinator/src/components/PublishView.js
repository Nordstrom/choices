import React from 'react';

import { togglePublish, namespacesLoaded } from '../actions';
import { toNamespace, fromNamespace } from '../nsconv';

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
  }).map(n => 
    <div key={n.name} className="checkbox">
      <label><input
        type="checkbox"
        checked={n.publish}
        onChange={() => dispatch(togglePublish(n.name))}
      /> {`${n.name} - ${n.experiments.map(e => e.name).join(', ')}`}
      </label>
    </div>
  );
  if (ns.length === 0) {
    return <p>No changes made</p>
  }
  return (
    <form onSubmit={e => {
      e.preventDefault();
      // filter selected
      const requests = namespaces.filter(n => n.publish)
      .map(n => {
        let url;
        let req;
        if (n.isNew) {
          url = '/api/v1/create';
          req = createRequest(n);
        } else {
          url = '/api/v1/update';
          req = updateRequest(n);
        }
        return fetch(url, req);
      });
      Promise.all(requests).then(responses => {
        const headers = new Headers({'Accept': 'application/json'});
        const req = { method: 'POST', headers: headers, body: JSON.stringify({ environment: "Staging" }) };
        const badRequest = (req, resp) =>  ({ err: "bad request", req, resp });
        fetch("/api/v1/all", req)
        .then(resp => {
          if (!resp.ok) {
            throw badRequest(req, resp);
          }
          return resp.json();
        })
        .then(json => {
          const ns = json.namespaces.map(n => {
            return toNamespace(n);
          });
          dispatch(namespacesLoaded(ns));
        })
      })
      .catch(err => console.log(err.err, err.req, err.resp));
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
