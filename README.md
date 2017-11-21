# go-service-narada4d-example
[![CircleCI](https://circleci.com/gh/powerman/go-service-narada4d-example.svg?style=svg)](https://circleci.com/gh/powerman/go-service-narada4d-example)

Example how to build stateful Go microservice based on Narada4D framework
as Docker image.

**WARNING!** This is initial quick-n-dirty version for experimenting.


## Building

### Method 1

Build, test and deploy to local docker swarm: `EXAMPLE_PORT=8080 ./build`.

Additionally it'll add generated files into current directory to make it
easier to test/run/investigate on host, without docker.

### Method 2

Build and test: `./circleci_build`

This method require CircleCI command-line tool `circleci` but otherwise
doesn't depend on using CircleCI.

### Method 3

Build and test all branches, push built image to Docker Hub after building
`master` branch: connect CircleCI and setup environment vars
`$DOCKER_USER`, `$DOCKER_PASS`.
