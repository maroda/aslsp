# Simple

An extremely simple front-back HTTP app written in Go, packaged as Docker containers, and deployable in Kubernetes.

The build process so far is simple:

```
cd ./cluster
kubectl apply -f craque-ns.yaml
cd ../front
docker build -t craque:Bv003 .
docker tag craque:Bv003 docker.io/maroda/craque:Bv003
docker push docker.io/maroda/craque:Bv003
kubectl -n crq apply -f craque.yaml
cd ../back
docker build -t craque:Cv003 .
docker tag craque:Cv003 docker.io/maroda/craque:Cv003
docker push docker.io/maroda/craque:Cv003
kubectl -n crq apply -f bacque.yaml
```

If it's a private repo, make sure to add the secret to the namespace:

	kubectl -n crq create secret docker-registry crqregcred --docker-server='https://index.docker.io/v1/' --docker-username='USER' --docker-password='PASS' --docker-email='EMAIL'

Once deployed, *craque* will access the *bacque* server to retrieve the local time.

	http://craque_loadbalancer_url/dt

Both apps also have a `/ping` endpoint for configuring liveliness tests (not yet configured).

Going to any invalid endpoint (e.g.: `/`, `/foo`, `/pickles`) will simply return "Hello. `/<endpoint>`"
