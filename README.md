# jsonbytes [![Go Reference](https://pkg.go.dev/badge/github.com/theteacat/jsonbytes.svg)](https://pkg.go.dev/github.com/theteacat/jsonbytes)

Package jsonbytes provides utilities for operating on JSON values expressed as `[]byte`. There are various operations you may want to perform on a JSON value that may be a bit quicker or more memory efficient to perform without unmarshalling it.

Note that this package is niche; if the JSON you want to operate on has to be unmarshalled at some stage anyway, it will probably be more efficient to operate on it after it has been unmarshalled.



## Example Uses

[The `examples` directory](./examples) contains various example use cases for the `jsonbytes` package.

- [`jsonfinder`](./examples/jsonfinder) is A CLI tool for checking if all the files in a directory are valid JSON.




## Tests

Use `go test` and `go tool cover` to generate a coverage report and open it in your browser:

```bash
go test -coverprofile coverage.out .
go tool cover -html coverage.out
```



## Benchmarks

You can run the benchmarks using `go test`:

```bash
go test -bench=.
```

You can also profile the benchmarks and open the results in your browser using `go tool pprof`:

```bash
go test -bench=. -benchmem -cpuprofile profile.out
go tool pprof -http localhost:8080 profile.out
```

Most of the benchmarks use JSON values that are at least 10KiB in size, except a few which use a sample JSON object that is included in [the `testdata` directory](./testdata); it took me a while to decide what to use here before I remembered that `package-lock.json`s are famously large, and [axios's GitHub repository includes a whopping 1.6MB one](https://github.com/axios/axios/blob/v1.x/package-lock.json).

As you can see in the example run below, when benchmarked on validating axios' package-lock.json, `jsonbytes.IsJson` is currently ~5x faster and uses ~5,000,000x less memory than `json.Unmarshal`! 😱

```
goos: darwin
goarch: arm64
pkg: github.com/theteacat/jsonbytes
cpu: Apple M1 Pro
BenchmarkIsJson/LongString/JsonBytes-8                                     89312             13149 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongString/EncodingJson-8                                  28381             42099 ns/op           10416 B/op          4 allocs/op
BenchmarkIsJson/LongNumber/JsonBytes-8                                    123584              9699 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/JsonBytes-8                                      142039              8450 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/EncodingJson-8                                    28075             42747 ns/op           10752 B/op          6 allocs/op
BenchmarkIsJson/LongWhitespace/JsonBytes-8                                134355              8924 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongWhitespace/EncodingJson-8                              25742             46590 ns/op             176 B/op          3 allocs/op
BenchmarkIsJson/longArray/JsonBytes-8                                      24752             48459 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/longArray/EncodingJson-8                                    4674            246783 ns/op          281134 B/op       5138 allocs/op
BenchmarkIsJson/LongObject/JsonBytes-8                                     56427             21327 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongObject/EncodingJson-8                                   6472            173477 ns/op          171253 B/op       1287 allocs/op
BenchmarkIsJson/ManyArrays/JsonBytes-8                                     38258             31361 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyArrays/EncodingJson-8                                   5414            207676 ns/op          240192 B/op       3432 allocs/op
BenchmarkIsJson/ManyObjects/JsonBytes-8                                    37081             32362 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyObjects/EncodingJson-8                                  3830            296207 ns/op          563244 B/op       4406 allocs/op
BenchmarkIsJson/ManyTrues/JsonBytes-8                                      45247             26463 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyTrues/EncodingJson-8                                   18552             64158 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyFalses/JsonBytes-8                                     56166             21354 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyFalses/EncodingJson-8                                  20533             58302 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyNulls/JsonBytes-8                                      45142             26563 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyNulls/EncodingJson-8                                   18396             65027 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/PackageLockAxios/JsonBytes-8                                 511           2340267 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/PackageLockAxios/EncodingJson-8                               99          11803546 ns/op         5256552 B/op      97408 allocs/op
BenchmarkRedactAllValues/LongString/JsonBytes-8                            91141             13162 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongString/EncodingJson-8                         28399             42311 ns/op           10440 B/op          5 allocs/op
BenchmarkRedactAllValues/LongNumber/JsonBytes-8                           123586              9707 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/JsonBytes-8                              90872             13174 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/EncodingJson-8                           22627             53134 ns/op           21242 B/op         10 allocs/op
BenchmarkRedactAllValues/LongWhitespace/JsonBytes-8                       134295              8979 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongWhitespace/EncodingJson-8                     25626             46841 ns/op             180 B/op          4 allocs/op
BenchmarkRedactAllValues/longArray/JsonBytes-8                             20431             58698 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/longArray/EncodingJson-8                           2864            411853 ns/op          292615 B/op       5140 allocs/op
BenchmarkRedactAllValues/LongObject/JsonBytes-8                            38770             30964 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongObject/EncodingJson-8                          2332            497395 ns/op          280932 B/op       3814 allocs/op
BenchmarkRedactAllValues/ManyArrays/JsonBytes-8                            30108             39832 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyArrays/EncodingJson-8                          2926            398583 ns/op          333932 B/op       6847 allocs/op
BenchmarkRedactAllValues/ManyObjects/JsonBytes-8                           28532             41870 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyObjects/EncodingJson-8                         1824            635670 ns/op          693860 B/op       8798 allocs/op
BenchmarkRedactAllValues/ManyTrues/JsonBytes-8                             38852             30871 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyTrues/EncodingJson-8                           9408            123684 ns/op          112655 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyFalses/JsonBytes-8                            48108             24895 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyFalses/EncodingJson-8                         10000            107630 ns/op          111191 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyNulls/JsonBytes-8                             38965             30768 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyNulls/EncodingJson-8                          13689             87844 ns/op          112671 B/op         19 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/JsonBytes-8                        423           2810718 ns/op               2 B/op          0 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/EncodingJson-8                      58          19496053 ns/op         9176881 B/op     177204 allocs/op
BenchmarkRedactAllValuesPackageLockAxiosEncodingJsonPremarshalled-8          805           1444854 ns/op             643 B/op         26 allocs/op
PASS
ok      github.com/theteacat/jsonbytes  158.463s
```
