# etime [![Build Status](https://travis-ci.org/go-eden/etime.svg?branch=master)](https://travis-ci.org/go-eden/etime)

`etime` extend golang's `time` package, to provide more and better features.

# Install

```bash
go get github.com/go-eden/common/etime
```

# Usage
 
This section shows all features, and their usage.
 
## Current Timestamp

This feature works like Java's `System.currentTimeMillis()`, it will return `int64` value directly:

+ `NowSecond`: obtains the current second, use syscall for better performance
+ `NowMillisecond`: obtains the current microsecond, use syscall for better performance
+ `NowMicrosecond`: obtains the current millisecond, use syscall for better performance

For better performance, `Now*` didn't use `time.Now()`, because it's a bit slow. 

### Demo

```go
package main

import (
	"github.com/go-eden/common/etime"
	"time"
)

func main() {
	println(etime.NowSecond())
	println(etime.NowMillisecond())
	println(etime.NowMicrosecond())

	println(time.Now().Unix())
	println(time.Now().UnixNano() / 1e6)
	println(time.Now().UnixNano() / 1e3)
}
```

### Performance

In my benchmark, the performance of `etime.NowSecond` was about `40 ns/op`, the performance of `time.Now()` was about `68 ns/op`:

```
BenchmarkNowSecond-12         	29296284	        39.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkNowMillisecond-12    	29361312	        40.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkNowMicrosecond-12    	29742286	        40.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkTimeNow-12           	15010953	        70.0 ns/op	       0 B/op	       0 allocs/op
```

Under the same hardware environment, Java's `System.currentTimeMillis()` was like this:

```
Benchmark               Mode  Cnt   Score   Error  Units
TimestampBenchmark.now  avgt    9  25.697 ± 0.139  ns/op
```

Some library may be sensitive to this `28ns` optimization, like [slf4go](https://github.com/go-eden/slf4go). 

By the way, `System.currentTimeMillis()`'s implementation was similar with `etime.NowSecond`:

```c++
jlong os::javaTimeMillis() {
  timeval time;
  int status = gettimeofday(&time, NULL);
  assert(status != -1, "bsd error");
  return jlong(time.tv_sec) * 1000  +  jlong(time.tv_usec / 1000);
}
```

So, there should have room for improvement.

# License

MIT
