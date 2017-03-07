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

> We are in the process of migrating to a supported Kubernetes
> cluster. When that migration occurs we will update our _prod_
> endpoint to `http://prod.elwin.aws.cloud.nordstrom.net`.

`http://dev.elwin.aws.cloud.nordstrom.net` is our current _dev_
endpoint.

> There are other dev endpoints but they should pointed to the same
> service as the above. They are
> `http://elwin.k8s-a.ttoapps.aws.cloud.nordstrom.net` and
> `http://elwin-test.ttoapps.aws.cloud.nordstrom.net`.

## Elwin request

Elwin's request structure is very simple. In it's current state,
clients will make a GET request to one of the endpoints. The request
requires two params. One is specifies the team making the request. The
second specifies the user's id.

The team name can currently be specified in query params `team`,
`label`, `teamid`, or `group-id`. If any of those params are not
blank, they will be used to filter the experiments.

The user id is specified in the param `userid`. In most cases this
should be the `ExperimentID` from the `experiments` cookie on web
requests.

You can also supply other query params to match labels on your
experiments. For example if you are running a test that should only be
shown to desktop users, you could set the label `platform` with the
value `desktop` on your experiment. When you query for desktop
experiments you would then include the `platform=desktop` query param.
Another example is if you wanted to run a test internal only before
deploying it to customers you could add a label
`traffic-source=internal` and query for it the same way. You can query
for multiple values for the same label key by repeating the query in
the request. For example, `?env=prod&env=dev`. The query params create
an *and* selection on the labels of your experiments. The results
returned will be the union of all experiments whose labels match.

The full request for experiments that are for the ato team in the dev
and prod _env_ironment for customers browsing on desktop platform for
the userid `andrew` would look like the following.

```
http://dev.elwin.aws.cloud.nordstrom.net/?team=ato&env=prod&env=dev&platform=desktop&userid=andrew
```

## Elwin response

All the experiments a user qualifies for will be returned in the
experiment response. The response is JSON.

```javascript
{
  "experiments": {
    "personalized-header-experiment": {
      "namespace": "aaaaaa",
      "params": {
        "personalized": "default"
      }
    },
    "button-experiment": {
      "namespace": "bbbbbb",
      "params": {
        "button-color": "blue",
        "button-size": "large"
      }
    }
  }
}
```

> In this example there are two experiments:
> `personalized-header-experiment` and `button-experiment`. In
> `personalized-header-experiment` the user was hashed into the
> `default` experience. In `button-experiment` they were hashed into
> `blue` + `large` MVT experience.

The top level of the object will contain only a single key,
`experiments`. The `experiments` key contains a map of experiment
names to experiment values, represented as an object. 

Experiment names will be keys in the `experiments` object. The values
for the experiment names will be objects with two keys. `namespace`,
the first key, contains a string value. It is not essential to the
experiment but is required for the data collection. `params`, the
second key, contains an object of param names and param values.

The `params` object has keys that correspond to param names. The
values of these keys will be the experience the user has been hashed
into. If you are running an A/B/N test there will only be one key. If
you are running a multivariate test then there be keys for each arm.
The values for params will always be returned as strings.