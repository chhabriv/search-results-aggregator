## API

The API documentation can be found as follows:
- [API Specs](https://github.com/chhabriv/search-results-aggregator/blob/main/spec/openapi-specs.yaml)
- [Human Friendly Documentation](https://chhabriv.github.io/search-results-aggregator/)

#### Pre-requisites

- Go 1.19.4
- Docker 20.10.21

##### To generate mocks for interface

Note: the mocks are checked in, the pre-requisites need to be run if the interface is changed and mocks need to re-generated.

1. Install mockgen

```shell
go install github.com/golang/mock/mockgen@v1.6.0
go get github.com/golang/mock/mockgen/model
```

2. Generate mocks
```shell
go generate ./...
```

##### To generate human friendly page for coverage report

Install go tool cover

```shell
go install golang.org/x/tools/cmd/cover@latest
```

## Running local builds

1. Start the search-results-aggregator service

```shell
go build
./search-results-aggregator
```

2. Perform health check to confirm whether the http service is up and running.

```shell
curl -X GET http://localhost:8080/healthcheck
```

You should see:

```json
{
    "status": "serving"
}
```

### Unit Tests

To run unit tests:

```shell
go test ./...
```

Generate static page with unit test coverage:

```shell
go test ./... -covermode=count -coverprofile=cov.out
go tool cover -html cov.out -o cover.html
open cover.html
```

## Running containerized builds

To run the service in a Docker container, do the following:

1. Create the application container using the [Dockerfile](https://github.com/chhabriv/search-results-aggregator/blob/main/Dockerfile)

```shell
docker build . -t local/search-results-aggregator:latest
```

2. To run the container

```shell
docker run --rm -p 8080:8080 local/search-results-aggregator:latest
```

## Sample cURL commands

_Can be run against the application running directly from the executable or from the container._

1. Healthcheck

```shell
curl --location --request GET 'http://localhost:8080/health'
```

2. Sample request to retrieve 5 links sorted by views 

```shell
curl --location --request GET 'localhost:8080/links?sortKey=views&limit=5'
```
