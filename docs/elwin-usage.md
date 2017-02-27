# Elwin Usage

* [Elwin endpoints](#elwin-endpoints)
* [Elwin request](#elwin-request)
* [Elwin response](#elwin-response)

## Elwin endpoints

Elwin uses two endpoints for experiments: _dev_ and _prod_. The _prod_
endpoint will serve tests that are considered live to customers. The
_dev_ endpoint can be used to preview experiments before enabling it
for your customers.

`http://elwin.ttoapps.aws.cloud.nordstorm.net` is our current _prod_
endpoint.

`http://dev.elwin.aws.cloud.nordstrom.net` is our current _dev_
endpoint.

We are in the process of migrating to a supported Kubernetes cluster.
When that migration occurs we will update our _prod_ endpoint to
`http://prod.elwin.aws.cloud.nordstrom.net`.

## Elwin request

Elwin's request structure is very simple. In it's current state,
clients will make a GET request to one of the endpoints. The request
requires two params. One is specifies the team making the request. The
second specifies the user's id.

The team name can currently be specified in query params `label`,
`teamid`, or `group-id`. If any of those are params are not blank,
they will be used to filter the experiments. If more than one is used,
it will prefer them in the order specified before.

The user id is specified in the param `userid`. In most cases this
should be the `ExperimentID` from the `experiments` cookie on web
requests.

## Elwin response

All the experiments a user qualifies for will be returned in the
experiment response. The response is JSON.

```javascript
{
  "experiments": {
    "some-experiment": {
      "namespace": "aaaaaa",
      "params": {
        "my-param-name": "value-for-my-user"
      }
    },
    "another-experiment": {
      "namespace": "bbbbbb",
      "params": {
        "first-param": "some-value",
        "second-param": "another-value"
      }
    }
  }
}
```

The top level of the object will contain only a single key,
`experiments`. The `experiments` key contains a map of experiment
names to experiment values, represented as an object. Experiment names
will be keys in the object. Experiment values will be objects with two
keys. `namespace`, the first key, contains a string value. It is not
essential to the experiment but is required for the data collection.
`params`, the second key, contains an object of param names and param
values. If you are running an A/B/N test there will only be one key.
If you are running a multivariate test then there be keys for each
arm. The values for params will always be returned as strings.