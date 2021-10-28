# esync

# `ReMutex`

An implementation of `reentrant lock`, based on `sync.Mutex`:

demo1:

```go
var m ReMutex

m.Lock()
m.Lock()
m.Unlock()
m.Unlock()
```

more demo in [`sync_remutex_test.go`](./sync_remutex_test.go)