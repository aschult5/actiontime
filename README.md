# Action Time
[![Build Status](https://travis-ci.com/aschult5/actiontime.svg?branch=master)](https://travis-ci.com/aschult5/actiontime)
[![codecov](https://codecov.io/gh/aschult5/actiontime/branch/master/graph/badge.svg)](https://codecov.io/gh/aschult5/actiontime)

A simple Golang library that accepts a json serialized string of the form below and maintains an average time for each action.
```json
{"action":"jump", "time":100}
{"action":"run", "time":75}
{"action":"jump", "time":200}
```

## Disclaimer
I have never before written *anything* in GoLang, but it has elegant support for concurrency and json, which are core aspects of this project. It is also used by the folks that will be inspecting this project. You've been warned!

## Usage
This module is intended to be imported from a go program.

### Dependencies
* [go 1.13](https://golang.org/dl/) has been tested
* python3 is required to generate test cases
* [golint](https://github.com/golang/lint) is required for contributions

### Installation
`go get github.com/aschult5/actiontime`

### Documentation
From `go doc`
```go
package actiontime // import "github.com/aschult5/actiontime"

Package actiontime takes actions and times as json, tracking average times.
Input is received as a json string, per requirements.

var ErrMissingInput = errors.New("actiontime: Missing input data")
type Stats struct{ ... }
```

From `go doc Stats`
```go
package actiontime // import "github.com/aschult5/actiontime"

type Stats struct {
        // Has unexported fields.
}
    Stats tracks passed actions' average times.

func (a *Stats) AddAction(input string) error
func (a *Stats) GetStats() string
```

### Example
```go
package main

import (
        "fmt"
        "github.com/aschult5/actiontime"
)

func main() {
        var a actiontime.Stats
        a.AddAction(`{"action":"jump", "time":100}`)
        a.AddAction(`{"action":"fall", "time":100}`)
        a.AddAction(`{"action":"jump", "time":200}`)
        a.AddAction(`{"action":"fall", "time":200}`)
        a.AddAction(`{"action":"sit", "time":500}`)
        a.AddAction(`{"action":"stand", "time":700}`)

        fmt.Println(a.GetStats())
}
```
Possible output:
```json
[{"action":"stand","avg":700},{"action":"jump","avg":150},{"action":"fall","avg":150},{"action":"sit","avg":500}]
```

## Testing
### Running Tests
```bash
go test [-race]
```

Some test case files will need to be manually generated, as they create large files that probably don't belong in revision control.

### Running Benchmarks
```bash
go test -bench . -benchmem -benchtime 10s -run=^$
```

See [travis-ci](https://travis-ci.com/aschult5/actiontime) for latest benchmark results on a single core.

### Checking code coverage
```bash
go test -coverprofile=coverage.txt -covermode=atomic
go tool cover -func coverage.txt
```
See [codecov](https://codecov.io/gh/aschult5/actiontime) for code coverage history.

### Generating Tests
See `python3 ./tools/testgenerator.py --help`  
Generated tests will have to be manually integrated by adding a new Test\* case to `statsimpl_test.go`

## Design
### Decisions
* RWMutex was chosen due to this problem's close relationship with the [readers-writers problem](https://en.m.wikipedia.org/wiki/Readers%E2%80%93writers_problem).
  * Some experimenting was done with using a channel instead, but initial benchmarks showed that the RWMutex solution generally outperformed. More testing is needed.
* `actiontime.statsImpl` exists to separate the interface from its implementation. This reduces `actiontime.Stats` to a simple message transcoder and input validator.
*  testgenerator.py was written before I discovered how easy it is to write test tables in Go. Under the same circumstances, I would choose the test tables next time. Benchmarks were written as tables.
*  statsWrapper was written to offer a common test interface between Stats and statsImpl, which allows the same test cases to be run against each.
  * Ideally I would unit test each method independently, using mocked dependencies to simulate behavior, and then integrate test against the exported methods. However, I have limited time and nothing to integrate.
  * Having proper unit tests and integration tests would eliminate the usefulness of statsWrapper


### Assumptions

* `time` may only be a json number greater than 0 that fits into a float64
* `action` may only be a json string of length between 1 and `actiontime.MaxActionLen` characters
* Case-insensitive keys; duplicate normalized keys will follow [go's preference](https://blog.golang.org/json-and-go)
* Case-sensitive values
* Extra fields can be ignored
* The set of valid `action` values is reasonably small, i.e. will fit into memory
* No need to persist inputs
* No need to track of the sums of `time` values
* Result of `GetStats` does not need to be strictly chronologically-accurate
  * Rationale:
    1. No indication that `AddAction` should be treated as a sensitive transaction, e.g. a bank deposit or withdrawal
    2. Averages are fuzzy by nature and in the long run individual calls to `AddAction` will have little effect
* Caller is responsible for formatting result of `GetStats`

## TODO
See github.com/aschult5/actiontime/issues

### Notable improvements
* Improve performance by...
  1. Asynchronously handling calls to AddAction (Issue #23)
* Better test coverage under load (Issue #28)
