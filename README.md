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
BenchmarkIsJson/LongString/JsonBytes-8                                     90008             13185 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongString/EncodingJson-8                                  28626             41906 ns/op           10416 B/op          4 allocs/op
BenchmarkIsJson/LongNumber/JsonBytes-8                                    123225              9728 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/JsonBytes-8                                      142057              8456 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/EncodingJson-8                                    28161             42585 ns/op           10752 B/op          6 allocs/op
BenchmarkIsJson/LongWhitespace/JsonBytes-8                                133995              8952 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongWhitespace/EncodingJson-8                              25714             46581 ns/op             176 B/op          3 allocs/op
BenchmarkIsJson/LongArray/JsonBytes-8                                      24678             48553 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongArray/EncodingJson-8                                    4728            248547 ns/op          281145 B/op       5139 allocs/op
BenchmarkIsJson/LongObject/JsonBytes-8                                     55938             21186 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongObject/EncodingJson-8                                   6655            174718 ns/op          171259 B/op       1286 allocs/op
BenchmarkIsJson/ManyArrays/JsonBytes-8                                     38143             31444 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyArrays/EncodingJson-8                                   5643            207253 ns/op          240193 B/op       3432 allocs/op
BenchmarkIsJson/ManyObjects/JsonBytes-8                                    36948             33444 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyObjects/EncodingJson-8                                  3886            295625 ns/op          563244 B/op       4406 allocs/op
BenchmarkIsJson/ManyTrues/JsonBytes-8                                      45248             26517 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyTrues/EncodingJson-8                                   18601             64407 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyFalses/JsonBytes-8                                     56048             21394 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyFalses/EncodingJson-8                                  20580             58530 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyNulls/JsonBytes-8                                      45226             26568 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyNulls/EncodingJson-8                                   18403             65189 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/DeeplyNestedArray/JsonBytes-8                              10000            100032 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/DeeplyNestedArray/EncodingJson-8                            2318            496566 ns/op          333203 B/op      10257 allocs/op
BenchmarkIsJson/DeeplyNestedObject/JsonBytes-8                             22525             53221 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/DeeplyNestedObject/EncodingJson-8                           3486            319813 ns/op          748369 B/op       4113 allocs/op
BenchmarkIsJson/PackageLockAxios/JsonBytes-8                                 510           2337905 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/PackageLockAxios/EncodingJson-8                               93          11872456 ns/op         5256640 B/op      97408 allocs/op
BenchmarkRedactAllValues/LongString/JsonBytes-8                            90871             13178 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongString/EncodingJson-8                         28302             42234 ns/op           10440 B/op          5 allocs/op
BenchmarkRedactAllValues/LongNumber/JsonBytes-8                           123178              9710 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/JsonBytes-8                              90960             13202 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/EncodingJson-8                           22594             53156 ns/op           21239 B/op         10 allocs/op
BenchmarkRedactAllValues/LongWhitespace/JsonBytes-8                       134090              8963 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongWhitespace/EncodingJson-8                     25446             46891 ns/op             180 B/op          4 allocs/op
BenchmarkRedactAllValues/LongArray/JsonBytes-8                             20289             59016 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongArray/EncodingJson-8                           2880            408689 ns/op          293343 B/op       5141 allocs/op
BenchmarkRedactAllValues/LongObject/JsonBytes-8                            38708             30947 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongObject/EncodingJson-8                          2335            495433 ns/op          280150 B/op       3811 allocs/op
BenchmarkRedactAllValues/ManyArrays/JsonBytes-8                            30004             39816 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyArrays/EncodingJson-8                          2944            395864 ns/op          333963 B/op       6847 allocs/op
BenchmarkRedactAllValues/ManyObjects/JsonBytes-8                           28574             41908 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyObjects/EncodingJson-8                         1827            632771 ns/op          694482 B/op       8798 allocs/op
BenchmarkRedactAllValues/ManyTrues/JsonBytes-8                             38805             30902 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyTrues/EncodingJson-8                           9538            122595 ns/op          112544 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyFalses/JsonBytes-8                            47998             25237 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyFalses/EncodingJson-8                         10000            106935 ns/op          111253 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyNulls/JsonBytes-8                             38931             30823 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyNulls/EncodingJson-8                          13684             87774 ns/op          112781 B/op         19 allocs/op
BenchmarkRedactAllValues/DeeplyNestedArray/JsonBytes-8                     10000            109746 ns/op               1 B/op          0 allocs/op
BenchmarkRedactAllValues/DeeplyNestedArray/EncodingJson-8                    652           1783565 ns/op          585871 B/op      19515 allocs/op
BenchmarkRedactAllValues/DeeplyNestedObject/JsonBytes-8                    18454             65018 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/DeeplyNestedObject/EncodingJson-8                  1213            950890 ns/op          934237 B/op      10263 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/JsonBytes-8                        427           2797300 ns/op               1 B/op          0 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/EncodingJson-8                      57          19354903 ns/op         8968428 B/op     177200 allocs/op
BenchmarkRedactAllValuesPackageLockAxiosEncodingJsonPremarshalled-8          812           1447568 ns/op             635 B/op         26 allocs/op
PASS
ok      github.com/theteacat/jsonbytes  172.681s
```
