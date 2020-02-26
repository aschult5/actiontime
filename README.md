# Concurrent Actions

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

