# jsonbytes [![Go Reference](https://pkg.go.dev/badge/github.com/theteacat/jsonbytes.svg)](https://pkg.go.dev/github.com/theteacat/jsonbytes)

Package jsonbytes provides utilities for operating on JSON values expressed as `[]byte`. There are various operations you may want to perform on a JSON value that may be a bit quicker or more memory efficient to perform without unmarshalling it, such as:

- [`IsJson(maybeJson []byte) error`](https://pkg.go.dev/github.com/theteacat/jsonbytes#IsJson): returns `nil` if `maybeJson` is valid JSON, else an error detailing why.
- [`RedactAllValues(inputJson []byte) ([]byte, error)`](https://pkg.go.dev/github.com/theteacat/jsonbytes#RedactAllValues): returns a new `[]byte` equivalent to `inputJson`, but with all the strings replaced with `""`, numbers replaced with `0` and booleans replaced with `true`; this may be useful if you want to log API request and response payloads that contain sensitive values.

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
BenchmarkIsJson/LongString/JsonBytes-8                                     89898             13166 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongString/EncodingJson-8                                  28560             41895 ns/op           10416 B/op          4 allocs/op
BenchmarkIsJson/LongNumber/JsonBytes-8                                    120060              9702 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/JsonBytes-8                                      142273              8411 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/EncodingJson-8                                    27973             42833 ns/op           10752 B/op          6 allocs/op
BenchmarkIsJson/LongWhitespace/JsonBytes-8                                134550              8923 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongWhitespace/EncodingJson-8                              25798             46500 ns/op             176 B/op          3 allocs/op
BenchmarkIsJson/longArray/JsonBytes-8                                      24693             48494 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/longArray/EncodingJson-8                                    4657            246836 ns/op          281145 B/op       5139 allocs/op
BenchmarkIsJson/LongObject/JsonBytes-8                                     56748             21502 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongObject/EncodingJson-8                                   6543            172730 ns/op          171244 B/op       1286 allocs/op
BenchmarkIsJson/ManyArrays/JsonBytes-8                                     38199             31375 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyArrays/EncodingJson-8                                   5499            207096 ns/op          240193 B/op       3432 allocs/op
BenchmarkIsJson/ManyObjects/JsonBytes-8                                    36949             32768 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyObjects/EncodingJson-8                                  3860            295325 ns/op          563245 B/op       4406 allocs/op
BenchmarkIsJson/ManyTrues/JsonBytes-8                                      45362             26488 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyTrues/EncodingJson-8                                   18580             64602 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyFalses/JsonBytes-8                                     56109             21386 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyFalses/EncodingJson-8                                  20599             58479 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyNulls/JsonBytes-8                                      45127             26489 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyNulls/EncodingJson-8                                   18494             64842 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/PackageLockAxios/JsonBytes-8                                 511           2334075 ns/op               1 B/op          0 allocs/op
BenchmarkIsJson/PackageLockAxios/EncodingJson-8                               99          11765450 ns/op         5255733 B/op      97405 allocs/op
BenchmarkRedactAllValues/LongString/JsonBytes-8                            91062             13152 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongString/EncodingJson-8                         28478             42121 ns/op           10440 B/op          5 allocs/op
BenchmarkRedactAllValues/LongNumber/JsonBytes-8                           123180              9706 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/JsonBytes-8                              90950             13186 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/EncodingJson-8                           22634             52932 ns/op           21233 B/op         10 allocs/op
BenchmarkRedactAllValues/LongWhitespace/JsonBytes-8                       134214              8938 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongWhitespace/EncodingJson-8                     25609             46737 ns/op             180 B/op          4 allocs/op
BenchmarkRedactAllValues/longArray/JsonBytes-8                             20401             59251 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/longArray/EncodingJson-8                           2866            410394 ns/op          293299 B/op       5141 allocs/op
BenchmarkRedactAllValues/LongObject/JsonBytes-8                            38293             30963 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongObject/EncodingJson-8                          2318            495435 ns/op          280236 B/op       3810 allocs/op
BenchmarkRedactAllValues/ManyArrays/JsonBytes-8                            30099             39970 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyArrays/EncodingJson-8                          2883            399193 ns/op          334166 B/op       6847 allocs/op
BenchmarkRedactAllValues/ManyObjects/JsonBytes-8                           28489             41835 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyObjects/EncodingJson-8                         1818            636559 ns/op          694086 B/op       8798 allocs/op
BenchmarkRedactAllValues/ManyTrues/JsonBytes-8                             38758             30914 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyTrues/EncodingJson-8                           9290            122919 ns/op          112672 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyFalses/JsonBytes-8                            48109             24886 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyFalses/EncodingJson-8                         10000            107101 ns/op          111297 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyNulls/JsonBytes-8                             38863             30859 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyNulls/EncodingJson-8                          13708             87966 ns/op          112667 B/op         19 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/JsonBytes-8                        426           2788767 ns/op               1 B/op          0 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/EncodingJson-8                      60          19153912 ns/op         8908000 B/op     177200 allocs/op
BenchmarkRedactAllValuesPackageLockAxiosEncodingJsonPremarshalled-8          844           1406587 ns/op             634 B/op         26 allocs/op
PASS
ok      github.com/theteacat/jsonbytes  158.705s
```
