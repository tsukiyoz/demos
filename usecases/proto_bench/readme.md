## 测试结果

```shell
proto_bench git:(main) ✗ make benchmark 
cd goprotobuf && go test -bench=.
goos: darwin
goarch: arm64
pkg: demo/goprotobuf
cpu: Apple M2
BenchmarkMarshal-8                       6380127               189.1 ns/op           384 B/op          1 allocs/op
BenchmarkUnmarshal-8                     4035450               297.2 ns/op           576 B/op          9 allocs/op
BenchmarkMarshalInParallel-8            11464707               130.6 ns/op           384 B/op          1 allocs/op
BenchmarkUnmarshalInParallel-8           7417598               175.2 ns/op           576 B/op          9 allocs/op
PASS
ok      demo/goprotobuf 6.329s
cd gogoprotobuf-fast && go test -bench=.
goos: darwin
goarch: arm64
pkg: demo/gogoprotobuf-fast
cpu: Apple M2
BenchmarkMarshal-8                      14281066                70.23 ns/op          384 B/op          1 allocs/op
BenchmarkUnmarshal-8                     6654720               174.4 ns/op           576 B/op          9 allocs/op
BenchmarkMarshalInParallel-8            19749049                61.21 ns/op          384 B/op          1 allocs/op
BenchmarkUnmarshalInParallel-8          10221326               121.1 ns/op           576 B/op          9 allocs/op
PASS
ok      demo/gogoprotobuf-fast  5.344s
cd gogoprotobuf-faster && go test -bench=.
goos: darwin
goarch: arm64
pkg: demo/gogoprotobuf-faster
cpu: Apple M2
BenchmarkMarshal-8                      16215603                71.50 ns/op          384 B/op          1 allocs/op
BenchmarkUnmarshal-8                     7059052               209.3 ns/op          544 B/op          9 allocs/op
BenchmarkMarshalInParallel-8            18164936                65.35 ns/op         384 B/op          1 allocs/op
BenchmarkUnmarshalInParallel-8          10211559               114.7 ns/op          544 B/op          9 allocs/op
PASS
ok      demo/gogoprotobuf-faster        5.707s
cd gogoprotobuf-slick && go test -bench=.
goos: darwin
goarch: arm64
pkg: demo/gogoprotobuf-slick
cpu: Apple M2
BenchmarkMarshal-8                      15635077                70.82 ns/op         384 B/op          1 allocs/op
BenchmarkUnmarshal-8                     6933984               176.1 ns/op          544 B/op          9 allocs/op
BenchmarkMarshalInParallel-8            18396477                62.21 ns/op         384 B/op          1 allocs/op
BenchmarkUnmarshalInParallel-8          10448790               117.0 ns/op          544 B/op          9 allocs/op
PASS
ok      demo/gogoprotobuf-slick 5.403s
```