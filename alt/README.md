# aslsp ALT

This is the _older_ version of the app where I had written two *separate* 'microservices'.

There is also some Istio stuff here and some guidance for how to compile this for ARM processors and other Operation hints.

<HR>

An extremely simple front-back HTTP app written in Go, packaged as Docker containers, and deployable in Kubernetes.

Something I've been using for a few years as materials for experiments and learnings from Raspberry Pi to Chaos Engineering.

Some of these instructions are particular to my own uses and some of them are probably out of date.

## Build

* Replace Cv0XX with the current version of **front/craque**
* Replace Bv0XX with the current version of **back/bacque**

```
cd ../front
docker build -t craque:Cv0XX .
docker tag craque:Cv0XX docker.io/maroda/craque:Cv0XX
docker push docker.io/maroda/craque:Cv0XX
cd ../back
docker build -t craque:Bv0XX .
docker tag craque:Bv0XX docker.io/maroda/craque:Bv0XX
docker push docker.io/maroda/craque:Bv0XX
```

## ARM

For running on Raspberry Pi 3 as a binary:
```
GOARCH=arm GOARM=7 GOOS=linux go build -o craque-arm front/
GOARCH=arm GOARM=7 GOOS=linux go build -o bacque-arm back/
```

## Operation

The app **front/craque** requires the environment variable `BACQUE` be set to the endpoint serving **back/bacque**.

* In kubernetes (front.craque.yaml) this is usually `"http://bacque:9999/fetch"`.
* Running the go app directly, use `export BACQUE="http://localhost:9999/fetch"`.
* With raw docker:
  * Start **bacque** first: `docker run --rm --name bacque -p 9999:9999 craque:Bv006`
  * Get the IP: `>>> export BCQ="http://$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' bacque):9999"`
  * Start **craque**: `docker run -e BACQUE=$BCQ --name craque -p 8888:8888 craque:Cv005`

Once deployed, *craque* will access the *bacque* server to display some dynamically retrieved data, including a datetime stamp.
Hitting [http://app.craq.io/dt]() will return something like this:

```
DateTime=201902202006
RequestIP=192.168.192.65
LocalIP=192.168.142.193
```

As of version Cv012, Craque will fall back to a local retrieval of DateTime if the endpoint set with BACQUE is unavailable (and returns with response code 418). It does *not* return the same "enriched sender/receiver IP data" that BACQUE does.

## Deploy to New Kubernetes Cluster (LoadBalancer service)

1. Configure context: `export KUBECONFIG=<ABS_PATH_CONFIG>`
2. Create the namespace: `kubectl apply -f 01_cluster/craque-ns.yaml`
3. Add docker registry private repo creds: `kubectl -n crq create secret docker-registry regcred --docker-server='https://index.docker.io/v1/' --docker-username='maroda' --docker-password='<REDACTED>' --docker-email='maroda@gmail.com'`
4. Deploy backend: `kubectl -n crq apply -f 02_back/lb-bacque.yaml`
5. Deploy frontend: `kubectl -n crq apply -f 03_front/lb-craque.yaml`
6. Get DNS for LoadBalancer: `export CLB=$(kubectl -n crq get svc craque -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')`
7. Apply DNS: `pushd tf && terraform apply -var "simple_lb=$CLB" -auto-approve && popd`

The last step requires AWS auth and a DNS zone already configured.

