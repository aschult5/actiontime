# Concurrent Actions
[![Build Status](https://travis-ci.com/aschult5/go-action-time.svg?branch=master)](https://travis-ci.com/aschult5/go-action-time)
[![codecov](https://codecov.io/gh/aschult5/go-action-time/branch/master/graph/badge.svg)](https://codecov.io/gh/aschult5/go-action-time)

A simple library that allows concurrent use of the following methods:

1. `addAction (string) returning error`
    This function accepts a json serialized string of the form below and maintains an average time for each action.

```json
{"action":"jump", "time":100}
{"action":"run", "time":75}
{"action":"jump", "time":200}
```

2. `getStats () returning string`
    Accepts no input and returns a serialized json array of the average time for each action that has been provided to the `addAction` function.
    Output after the sample calls above would be:

```json
{"action":"jump", "avg":150},
{"action":"run", "avg":75}
```

## Design
### Assumptions
* `time` values may only be json numbers
* `action` values may only be json strings
* Case-insensitive keys; duplicate normalized keys will follow [go's preference](https://blog.golang.org/json-and-go)
* Case-sensitive values
* Empty or missing fields should produce an error
* Extra fields can be ignored
* The set of valid `action` values is reasonably small, i.e. will fit into memory
* No need to persist inputs
* No need to track of the sums of `time` values
* Result of `getStats` does not need to be strictly chronologically-accurate
  * Rationale:
    1. No indication that `addAction` should be treated as a sensitive transaction, e.g. a bank deposit or withdrawal
    2. Averages are fuzzy by nature and in the long run individual calls to `addAction` will have little effect
* Caller is responsible for formatting result of `GetStats`

## Considerations
### Language
I have never before written *anything* in GoLang, but it has elegant support for concurrency and json, which are core aspects of this project.
