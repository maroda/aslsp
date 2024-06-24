# aslsp

**ASLSP** (... originally a piece by John Cage that means: _As SLow aS Possible_ ...) is an extremely simple HTTP app written in Go, packaged with Docker, and deployable in Kubernetes.

Originally used for testing Chaos Engineering tests, now a general purpose app that can operate in two modes: client or server.

The client on port 8888 is called: **CRAQUE**
The server on port 9999 is called: **BACQUE**

## Local Startup

Open two terminals. Start up CRAQUE, the client, on the first:

1. `go build aslsp`
2. `export BACQUE="http://localhost:9999/fetch"`
3. `./aslsp`

Start up BACQUE, the server, on the second:

1. `./aslsp -nofetch`

## Local Test

### Client

- Test the client in a browser at: [http://localhost:8888/dt](http://localhost:8888/dt)
- In a terminal: `curl localhost:8888/dt`

The result will look like (in macOS the IPv6 localhost shows up as the requestor):

```
DateTime=202406211513
RequestIP=::1
LocalIP=10.10.10.28
```

#### When the Server is Down

The client will fall-back in the case of BACQUE being unavailable. It will respond with HTTP Error Code 418 to indicate: "I'm not supposed to be showing you this, but the server is down so I am." In other words: I am a teapot, not a coffee maker. I know you want coffee, but here's some tea. This is a real response code that can be tracked.

```
I'm a teapot
DateTime=202406211524
```

### Server

- Test the server in a browser at: [http://localhost:9999/fetch](http://localhost:9999/fetch)
- In a terminal: `curl localhost:9999/fetch`
- Prometheus scrape: `curl localhost:9999/metrics`

The result will look identical to the client, because this is what the client fetches:

```
DateTime=202406211513
RequestIP=::1
LocalIP=10.10.10.28
```

## Kubernetes Deploy

Using `kubectl` install these yaml files in order:

1. `aslsp.yaml` (namespace)
2. `bacque.yaml` (server should be up first)
3. `craque.yaml` (BACQUE is set in the yaml)

## Kubernetes Test

### Client

- Get the server LoadBalancer address with something like: `kubectl -n aslsp get svc craque`
- Test the server in a browser, my local example is: [http://198.19.249.2/dt](http://198.19.249.2/dt)

The result:

```
>>> curl 198.19.249.2/dt
DateTime=202406212249
RequestIP=192.168.194.20
LocalIP=192.168.194.16
```

### Server

- Shell into one of the `craque` containers.
- Test from the commandline: `wget -q -O - http://bacque/fetch`

The result will be identical to what the client gets:

```
/go/src/aslsp # wget -q -O - http://bacque/fetch
DateTime=202406212356
RequestIP=192.168.194.20
LocalIP=192.168.194.15
```

## TODO

- Tests
- Can Codefresh do latest?
- Push to internal reg?
- Sync to public aslsp & publish to github via codefresh, using the same code

