## Objective

This repo is a "draft" reference for presentation, application or a task layer component. The goal is to demonstrate how the production code and tests are physically structured and logically separated. The primary objective of the conventions followed here are -

* Clearly defined physical structures of code.
* Maintain clear separation of concerns among different layers of code within a component. For example we want to  maintain a defined interface among different business layers (transport, business logic or data layer if applicable) for a component. No layer should have access to internal state or logic of any other layer. The communication between two layers must be gated via well defined interface so that each layer is testable independently.
* Define a set of conventions that every component will adhere to so that CI/CD system can find relevant source material in each stage of the pipeline. For example, build system should rely on convention in which folder to look for integration tests.


## Prerequisite
* Go distribution
```
If you are on Go 1.5, you will need to enable experimental vendoring by setting the GO15VENDOREXPERIMENT environmental variable to 1.
If you are using Go 1.6 and above, you are all set.
```


## Structure

The folders within the repo are organized into the following structure -
```
--{repo name}
------service
----------service.go
----------internal
--------------data
--------------logic
------main
----------main.go
------integration
----------client
----------tests
----------benchmark
------vendor
```

***

* folder name should be lower case.
* folder name must be same as package name.
* Unit tests colocate with corresponding logic or go source file in the same folder. For example, tests for *foo.go* will be maintained in *foo_test.go* in the same folder.

```
main, service and integration folders have specific intention and usage.
So when you name any folder such this would indicate that you are going to use them for their intended use.
```

#### service folder
The *service* folder and its sub folders contain all production logic. The *service* package exposes an entry-level function that accepts no parameter. This entry level function which we can call either *ListenAndServe* or *BindToRoutingKeyAndServe* will be invoked so that the service can initialize itself and transition to listening mode.

```go
// service.go

package service

import (
	"flag"
	"net/http"

	"github.com/go-martini/martini"
)

// ListenAndServe allows the the service to initialize itself and
// listen on a port so that it can serve incoming request(s)
func ListenAndServe() {
	h := getRequestHander()
	setup(h)
	http.ListenAndServe(bindTo, m)
}
```

* Either *ListenAndServe* or *BindToRoutingKeyAndServe* is the only function exposed by the *service* package.
* Any other package(s) defined inside the *service* folder should be internal for the following reasons -
  * These packages are interna details of this service and meant to be used by the *service* package exclusively. Shared packages should not live under *service* folder.
  * If the integration tests reside in the same repo, we want to make sure that test code can't access these internal data structures and logic.

The *service* package refers to the internal package as follows -
```go
package service

import (
	"github.com/rue-tkashem/echo-api/service/internal/logic"
)
```

#### main folder
The *main* package will define the top level entrypoint *main* function for the go binary we build. It will reference the *service* package as shown below.

```go
// main.go
package main

import (
	"github.com/rue-tkashem/echo-api/service"
)

func main() {
	service.ListenAndServe()
}
```

When we build the binary for the given service we will do the following -
```bash
# assuming the current folder is {repo}
cd main
go build -o ${repo_name}
```

#### integration folder
* All integration tests, benchmark tests and helper functions for these tests will be maintained in *integration* folder.


#### vendor folder
* *vendor* folder will be under the root folder of the repo. This ensures that the *vendor* folder is not included when we run unit or integration tests.
* We will vendor packages required by tests and production code insde the *vendor* folder. So there will be one designated *vendor* folder for a given repo.

The following sections will make an attempt to clarify how each stage will work.


## Build

In order to build the binary image for the given repo we need to do the following -

1. change to *main* folder
2. invoke *go build* command

```
All dependencies are vendored into the vendor folder. There is no further need to retrieve any dependencies.
```
```
You can specify an image name if you want to put the image in a desired folder
or give it a special name
    go build -o ${bin_folder}/{repo_name}
This will create an image same as the repo name and put it in your bin folder.    
```

GoCD specific
```
We will create static go binaries when we build.
We will inject version number, commit hash, pipeline id into the image when we build the binary
```

## Static Analysis

Static analysis will be executed on every go file in the repository. The following folders are excluded -
* *vendor* folder as specified above

GoCD pipeline specific:
```bash
# This shows an example of how we will execute a static analysis tool that accepts a package as input
# The current folder is the root folder of the repo.
for pkg in `go list ./... | grep -v "/vendor"`
do
  # -x will print commands as they are executed
  go vet -v -x $pkg >> ${output}
  ret_code=$?

  if [ $ret_code != 0 ]; then
   printf "Error : [%d] when executing vet on package: '$pkg'" $ret_code
   exit $ret_code
  fi	  
done
```
```bash
# This shows an example of how we will execute a static analysis tool that accepts directory as input
# The current folder is the root folder of the repo.
for dir in `go list -f {{.Dir}} ./... | grep -v "vendor"`
do
  gofmt -l -d ${dir} >> fmt.log
done
```


## Run Unit Tests

One easy way to run all unit tests is to change to  *service* folder and run the unit tests from there
```bash
cd service
go test ./...
```
This will traverse all folders recursively and run unit tests in each package found.

If you want to enable coverage then you might need to run unit tests for each package separately. You can retrieve the complete list of packages by invoking
```
  cd service
  go list ./...
```

```
vendor folder is excluded from any unit tests
```

GoCD pipeline specific:
```bash
# The current working directory is root of the repo.
# All packages except for those under the vendor and integration folders are in scope for unit testing.
package_path=$(go list ./... | grep -v "vendor\|integration")
for pkg in ${package_path[@]}
do
  go test -v $pkg
end
```

## Run Integration Tests:
```
The integration tests reside in the same repo as the production code. This does present a problem,
if you change test code only this will trigger a new unnecessary build.
With reproducible builds, this should not introduce any unknown issues or side effects.
But if we do decide to put the integration tests in separate repo the following should still apply.
```

* The integration tests are located under *integration/tests* folder.
* We will pass configuration values required by the tests via *go test* command line

That's how we can run the integration tests -
```bash
# The current working directory is root of the repo.
cd integration/tests
go test -v -endpoint=http://server:3000
```

```
The flag endpoint is a custom one we have added and points to the target service endpoint.
```

The code snippet below shows how we add custom flags to our test suite.
```go
// setup_test.go

func setup() {
	endpointFlagPtr := flag.String("endpoint", "", "Service endpoint under test")
	flag.Parse()
}

func shutdown() {
}

// TestMain has custom setup and shutdown
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()

	os.Exit(code)
}
```

GoCD pipeline specific:
The integration tests will be compiled into a binary. The binary will be packaged into a docker image and then it will be executed by a GoCD agent.

```bash
# The current working directory is root of the repo.
cd integration/tests
go test -c -o integration

# execute the test binary
./integration -test.v -endpoint=http://server:3000
```

## Run Benchmarking Tests:
* Benchmark tests are located inside *integration/benchmark* folder.
* Any shared helper functions or client library are located inside *client* folder.

We will run benchmarking tests similar to integration tests -
```bash
# The current working directory is root of the repo.
cd integration/benchmark
go test -v -endpoint=http://localhost:3000 -bench=.
```

```
- If you want to run the tests concurrently, pass desired values of GOMAXPROCS
- GOMAXPROCS=2 go test -endpoint=http://localhost:3000 -bench=.
- go test will execute the benchmark tests concurrently from two go routines
```

## Coverage from Integration Tests

In order to get code coverage from integration tests we will host the service in the same process space as the integration tests. Since both tests and service code are hosted in the same process *go test* will be able to measure the code coverage statistics and generate cover profile.

In order to host the service in the same process space as tests, we will introduce another flag *inproc*. If set to *true* then the service will be hosted inside *TestMain*
```go
package tests

import (
	"testing"

	"github.com/rue-tkashem/echo-api/service"
)

func setup() {
	inprocFlagPtr := flag.String("inproc", false, "Indicates whether the service will be hosted in process")
	flag.Parse()

	if *inprocFlagPtr {
		// the service package exposes a function that allows us to host the service in process
		go service.ListenAndServe()
	}
}

func shutdown() {
}

// TestMain has custom setup and shutdown
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()

	os.Exit(code)
}
```

Once in-process hosting is setup, we need to run the following commands

```bash
# The current working directory is root of the repo.
cd integration/tests
go test -covermode=set -coverprofile="$GOPATH/bin/cover.out" -endpoint=http://localhost:3000 -inproc=true -coverpkg="a,b,c"
go tool cover -html="$GOPATH/bin/cover.out"
```

The *-coverpkg* option allows you to specify a comma separated list of packages for which we want to mesaure code coverage. It is very easy to retrieve a complete list of packages under the service folder -
```bash
# The current working directory is root of the repo.
cd service
go list ./...
```
The integration test code does not directly reference the server code, it uses an URL to invoke the target APIs. For this reason we need to explicitly specify the list of production code packages for which we want to generate coverage profile.

```
go test will generate warning that these packages are not directly referenced, but we can safely ignore them.
```

## Data Race Detection Tests
* We will reuse integration test suite and invoke these tests in order to detect data race conditions. This is as simple as
building a race enabled image and point your integration tests to the race enabled build.
```bash
# The current working directory is root of the repo.
cd main
go build -race -o datarace
```
This will generate an instrumented build with data race detection enabled.


go doc
------
1. Run go doc server locally, it automatically detects your code in GOPATH.
```
godoc -http=":6060"
```
2. You can specify your desired package by opening a browser
```
http://localhost:6060/pkg/github.com/rue-tkashem/echo-api/integration/client/      
```
