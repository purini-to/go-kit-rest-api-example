# go-kit-rest-api-example

go-kit sample using chi and zap.

## Getting Started

#### Clone and run

```bash
git clone https://github.com/purini-to/go-kit-rest-api-example
cd go-kit-rest-api-example
go run cmd/api/main.go
# or debug run
# go run cmd/api/main.go -debug
```

#### Go to [http://localhost:8080/tasks](http://localhost:8080/tasks)


## Run on kubernetes

build docker image
```bash
# minikube
# eval $(minikube docker-env)
docker build . -t go-kit-rest-api-example:v1.0.0
```

deploy kubernates
```bash
kubectl apply -k ./k8s/
```