# contactgraph

This project aims to provide an example of golang architecture for a simple project around contacts.

## Build

To do a local build, just do:

```
go build
```

To generate a docker image:

```
docker build -t contatcgraph .
```

## Test

All the test are based on golang test package:

```
go test .
```

## Run

The service can be launch directly based on default values

```
$> ./contactgraph
```

and can be configured through env variable:

* `HTTP_ADDR`: the listen string representation like ":8080"
* `LOG_LEVEL`: define the level of log. default is `info`

## Architecture principles

This repository try to provide a possible golang service architecture which is describe [here](./doc/ddd.md)

## API

The documentation about the API can be found [here](./doc/api.md)