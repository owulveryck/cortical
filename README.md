[![Coverage Status](https://coveralls.io/repos/github/owulveryck/cortical/badge.svg?branch=master)](https://coveralls.io/github/owulveryck/cortical?branch=master)
[![Build Status](https://travis-ci.org/owulveryck/cortical.svg?branch=master)](https://travis-ci.org/owulveryck/cortical)
[![](https://godoc.org/github.com/owulveryck/cortical?status.svg)](http://godoc.org/github.com/owulveryck/cortical)

![Picture](https://github.com/owulveryck/cortical/raw/master/doc/cortical.png)

# Cortical

```go
import "github.com/owulveryck/cortical"
```

## What is Cortical?

Cortical is a go ~~framework~~ ~~middleware~~ piece of code that acts as a message dispatcher. Then messages are transmitted in full duplex over a websocket.
Cortical is therefore a very convenient way to distribute messages to "processing units" (other go functions) and to get the responses back in a concurrent way.

The "processing units" are called _Cortexes_ and do not need to be aware of any web mechanism.

### Cortexes

A cortex is any function that provides the two methods:

* A "send" function that returns a channel of `[]byte`. The content of the channel is sent to the websocket once available (cf [`GetInfoFromCortexFunc`](https://godoc.org/github.com/owulveryck/cortical#GetInfoFromCortexFunc))
* A "receive" method that take a pointer of `[]byte`. This function is called each time a message is received (cf [`SendInfoToCortex`](https://godoc.org/github.com/owulveryck/cortical#SendInfoToCortex))

Cortical take care of extracting and sending the `[]byte` to the websocket and dispatches them through all the cortexes.

### Registering the cortexes and creating a [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)

By now, the registration of cortexes is done at the creation of the Cortical object.

*WARNING* This will probably change in the future

```go
brain := &cortical.Cortical{
     Upgrader: websocket.Upgrader{},
     Cortexs:  []func(context.Context) (cortical.GetInfoFromCortexFunc, cortical.SendInfoToCortex){NewCortex}, 
}
http.HandleFunc("/ws", brain.ServeWS)
log.Fatal(http.ListenAndServe(":8080", nil))
```

### Examples of cortex

See the [example in godoc](https://godoc.org/github.com/owulveryck/cortical#example-package)

## Why is it coded in go?

1. Because I love go
2. Because I only know how to code in go
3. *The real reason*: concurrency (take a look at [Rob Pike - 'Concurrency Is Not Parallelism' on youtube](https://www.youtube.com/watch?v=cN_DpYBzKso&t=680s) to truly understand why)

# Caution

The API may change a lot; use it at your own risks. PR are welcome.
