// @flow
import React from 'react';

import { changesClear, namespaceTogglePublish, entitiesLoaded, namespaceLocalDelete } from '../actions';
import { toEntities, fromNamespace } from '../nsconv';

function createRequest(entities, namespace) {
  return {
    method: 'POST',
    body: JSON.stringify({ namespace: fromNamespace(entities, namespace) }),
  }
}

function deleteRequest(namespace) {
  return {
    method: 'POST',
    body: JSON.stringify({name: namespace.name}),
  }
}

function updateRequest(entities, namespace) {
  return {
    method: 'POST',
    body: JSON.stringify({ namespace: fromNamespace(entities, namespace) }),
  }
}

const changeList = (changes) => changes.map((c, ci) => {
  const details = Object.keys(c)
  .filter(k => k !== 'type')
  .map((k, ki) => <p key={ki}><strong>{k}</strong>={c[k]}</p>);
  return <li key={ci}>{c.type}: {details}</li>
});

const PublishView = ({ namespaces, changes, entities, dispatch }: {namespaces: Array<Object>, changes: Array<Object>, entities: Object, dispatch: Function}) => {
  const ns = namespaces.map(n => 
    <div key={n.name} className="checkbox">
      <label><input
        type="checkbox"
        checked={n.publish}
        onChange={() => dispatch(namespaceTogglePublish(n.name))}
      /> {n.name}
      </label>
      <ol>{changeList(changes[n.name])}</ol>
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
        if (n.isNew && n.delete) {
          return Promise.resolve(dispatch(namespaceLocalDelete(n.name)));
        }
        else if (n.isNew) {
          url = '/api/v1/create';
          req = createRequest(entities, n);
        } else if (n.delete) {
          url = '/api/v1/delete';
          req = deleteRequest(n);
        } else {
          url = '/api/v1/update';
          req = updateRequest(entities, n);
        }
        return fetch(url, req);
      });
      Promise.all(requests).then(responses => {
        dispatch(changesClear(namespaces.filter(n => n.publish).map(n => n.name)));

        const headers = new Headers({'Accept': 'application/json'});
        const req = {
          method: 'POST',
          headers: headers,
          body: JSON.stringify({ environment: "Staging" })
        };
        const badRequest = (req, resp) =>  ({ err: "bad request", req, resp });
        fetch("/api/v1/all", req)
        .then(resp => {
          if (!resp.ok) {
            throw badRequest(req, resp);
          }
          return resp.json();
        })
        .then(json => {
          const entities = toEntities(json.namespaces)
          dispatch(entitiesLoaded(entities));
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
  changes: React.PropTypes.object.isRequired,
  entities: React.PropTypes.object.isRequired,
}

export default PublishView;
