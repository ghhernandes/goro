# Goro

Goro is a simple Go package that helps with working with Concurrency.

## FanOut

Fan-out receives a channel and multiplexes it into N channels. N
goroutines will be started to read from the input and write to the output
channels.

```go
// input = <-chan int
input := generator()

// outputs = []<-chan int
outputs := goro.FanOut(ctx, input, 5)
```

## FanIn

FanIn receives multiple channels and returns a single channel. N+1 goroutines
will be started. Where N is equal to the length of input channels and another
channel to wait goroutines to finish.

```go
// inputs = []<-chan int
// output = <-chan int
output := goro.FanIn(ctx, inputs)
```

## Map

Map is helpful when applying functions to channels is necessary. A new
goroutine will be created to receive input data, apply the map function and
send to the output channel.
```go

byteCh := ByteStream()

toStr := func(input []byte) string {
    return string(input)
}

strCh := goro.Map(ctx, byteCh, toStr)

for value := range strCh {
    fmt.Println(value)
}
```

## Filter

Filter is helpful to filter data and send to a channel. This function returns a
channel that only returns data thats the filter function returns `true`.

```go
responsesCh := Stream()

badRequest := func(input <-chan http.Response) bool {
    return input.Request.StatusCode == http.StatusBadRequest 
}

// filteredCh only receives data thats filter function 
filteredCh := goro.Filter(ctx, badRequest, responsesCh)
```
