# Something Broke

Steps for checking what's broken.

## A test I created is not showing up when I query elwin.

1. Check if your labels match the query params

```bash
curl prod.elwin.aws.cloud.nordstrom.net/?userid=andrew&team=blah&env=dev
```

Your experiment should have labels that match all the query params
except `userid`. In this example, your experiment would need labels
`team` and `env` with the values `blah` and `dev`.

2. Check the test was actually created.

```bash
# dev environment
curl -X POST -d '{"environment": 0}' json-gateway.ato.platform.prod.aws.cloud.nordstrom.net/api/v1/all
# prod environment
curl -X POST -d '{"environment": 1}' json-gateway.ato.platform.prod.aws.cloud.nordstrom.net/api/v1/all
```

If you do not see the experiment in either output then it needs to be
recreated. Use the [form][neo-form] to create an experiment.

If your test is in the wrong environment, you should use
[houston][houston] to launch or delete tests from dev and prod.

3. Check if elwin is failing to update

If you verfied the test is in storage and the query is correct then
elwin may not be able to contact the storage server. You can first
check the logs.

```bash
$ kubectl get po -l run=elwin
NAME                         READY     STATUS    RESTARTS   AGE
elwin-3750452436-4vshr       1/1       Running   0          10d
elwin-3750452436-6m78g       1/1       Running   0          11d
elwin-3750452436-7rnmx       1/1       Running   0          10d
elwin-3750452436-ffnq0       1/1       Running   0          11d
elwin-3750452436-gh607       1/1       Running   0          10d
elwin-3750452436-jlkcd       1/1       Running   0          11d
elwin-3750452436-p298j       1/1       Running   0          11d
elwin-3750452436-rqh9l       1/1       Running   0          11d
elwin-3750452436-sp9c0       1/1       Running   0          11d
elwin-3750452436-znccr       1/1       Running   0          11d
elwin-dev-3169672050-0sq6s   1/1       Running   0          10d
elwin-dev-3169672050-3r0w2   1/1       Running   0          14d
elwin-dev-3169672050-f3kcd   1/1       Running   0          11d
elwin-dev-3169672050-gxlqj   1/1       Running   0          10d
elwin-dev-3169672050-lzf50   1/1       Running   0          14d
```

Select a pod and check it's logs.

```bash
kubectl logs --tail=50 elwin-3750452436-4vshr
```

Check for lines like the following.

```bash
2017/03/13 17:14:17 grpc: addrConn.resetTransport failed to create client transport: connection error: desc = "transport: dial tcp: lookup elwin-storage: no such host"; Reconnecting to {elwin-storage:80 <nil>}
2017/03/13 17:14:19 grpc: addrConn.resetTransport failed to create client transport: connection error: desc = "transport: dial tcp: lookup elwin-storage: no such host"; Reconnecting to {elwin-storage:80 <nil>}
```

This means grpc has lost it's connection to storage. You can try restarting the pod, but it should restart automatically after `UPDATE_FAIL_TIMEOUT` passes (default is 15m).

[neo-form]: http://127.0.0.1 "neo form"
[houston]: http://houston.ato.platform.prod.aws.cloud.nordstrom.net "Launch Control"
