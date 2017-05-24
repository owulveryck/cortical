[![Coverage Status](https://coveralls.io/repos/github/owulveryck/cortical/badge.svg?branch=master)](https://coveralls.io/github/owulveryck/cortical?branch=master)
[![Build Status](https://travis-ci.org/owulveryck/cortical.svg?branch=master)](https://travis-ci.org/owulveryck/cortical)
[![](https://godoc.org/github.com/owulveryck/cortical?status.svg)](http://godoc.org/github.com/owulveryck/cortical)
[![Report card](https://goreportcard.com/badge/github.com/owulveryck/cortical)](https://goreportcard.com/report/github.com/owulveryck/cortical)

![Picture](https://github.com/owulveryck/cortical/raw/master/doc/cortical.png)

# Cortical

```go
import "github.com/owulveryck/cortical"
```

## What is Cortical?

Cortical is a go ~~framework~~ ~~middleware~~ piece of code that acts as a message dispatcher. The messages are transmitted in full duplex over a websocket.
Cortical is therefore a very convenient way to distribute messages to "processing units" (other go functions) and to get the responses back in a **concurrent** and **asynchronous** way.

The "processing units" are called _Cortexes_ and do not need to be aware of any web mechanism.

## _Cortical_? _Cortex_? is it related to ML?

Actually I have developed this code as a support for my tests with tensorflow, Google Cloud Plafeform and AWS ML services.
I needed a way to capture images from my webcam and to speak out loud. I have used my chrome browser for this purpose.
Every Cortex is a specific ML implementation (eg a _Memory Cortex_ that captures all the images and send them to a cloud storage is needed for training models).

See my [blog post _"Chrome, the eye of the cloud - Computer vision with deep learning and only 2Gb of RAM"_](https://blog.owulveryck.info/2017/05/16/chrome-the-eye-of-the-cloud---computer-vision-with-deep-learning-and-only-2gb-of-ram.html) for more explanation.

### Cortexes

A cortex is any go code that provides two functions:

* A "send" function that returns a channel of `[]byte`. The content of the channel is sent to the websocket once available (cf [`GetInfoFromCortexFunc`](https://godoc.org/github.com/owulveryck/cortical#GetInfoFromCortexFunc))
* A "receive" method that take a pointer of `[]byte`. This function is called each time a message is received (cf [`SendInfoToCortex`](https://godoc.org/github.com/owulveryck/cortical#SendInfoToCortex))

A cortex object must therefore be compatible with the `cortical.Cortex` interface:

ex:
```go
// echo is a dummy type that reads a message, and send back an "ack"
type echo struct{}

func new() *echo {
	return &echo{}
}

// NewCortex is filling the cortical.Cortex interface
func (e *echo) NewCortex(ctx context.Context) (cortical.GetInfoFromCortexFunc, cortical.SendInfoToCortex) {
	c := make(chan []byte)
	return func(ctx context.Context) chan []byte {
			return c
		}, func(ctx context.Context, b *[]byte) {
			c <- *b
		}
}
```

Cortical take care of extracting and sending the `[]byte` to the websocket and dispatches them through all the cortexes.

### Registering the cortexes and creating a [http.HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)

By now, the registration of cortexes is done at the creation of the Cortical object.

*WARNING* This will probably change in the future

```go
brain := &cortical.Cortical{
     Upgrader: websocket.Upgrader{},
     Cortexes:  []cortical.Cortex{
          &echo{}, // See example in the previous paragraph
     }, 
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

# TODO

- [x] Tests
- [x] Doc
- [x] Benchmarks
- [ ] Demo
- [ ] More functional tests
- [ ] Using the context cancel mechanism instead of my `done chan`
- [ ] Evaluating the opportunity of changing the interface for an `io.ReadWriteCloser`
- [ ] More and better doc
