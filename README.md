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
BenchmarkIsJson/LongString/IsJson-8                                                        89767             13210 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongString/EncodingJson-8                                                  28606             42841 ns/op           10416 B/op          4 allocs/op
BenchmarkIsJson/LongNumber/IsJson-8                                                       122943              9699 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/IsJson-8                                                         142365              8434 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongName/EncodingJson-8                                                    28135             42540 ns/op           10752 B/op          6 allocs/op
BenchmarkIsJson/LongWhitespace/IsJson-8                                                   134277              8928 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongWhitespace/EncodingJson-8                                              25760             46573 ns/op             176 B/op          3 allocs/op
BenchmarkIsJson/longArray/IsJson-8                                                         24716             48505 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/longArray/EncodingJson-8                                                    4634            250700 ns/op          281131 B/op       5138 allocs/op
BenchmarkIsJson/LongObject/IsJson-8                                                        56740             21179 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/LongObject/EncodingJson-8                                                   6640            172589 ns/op          171248 B/op       1287 allocs/op
BenchmarkIsJson/ManyArrays/IsJson-8                                                        38317             31371 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyArrays/EncodingJson-8                                                   5436            207378 ns/op          240193 B/op       3432 allocs/op
BenchmarkIsJson/ManyObjects/IsJson-8                                                       37048             32436 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyObjects/EncodingJson-8                                                  3840            296951 ns/op          563246 B/op       4406 allocs/op
BenchmarkIsJson/ManyTrues/IsJson-8                                                         45331             26472 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyTrues/EncodingJson-8                                                   18610             64220 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyFalses/IsJson-8                                                        55848             21343 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyFalses/EncodingJson-8                                                  20574             58339 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/ManyNulls/IsJson-8                                                         45309             26479 ns/op               0 B/op          0 allocs/op
BenchmarkIsJson/ManyNulls/EncodingJson-8                                                   18403             66726 ns/op          100920 B/op         17 allocs/op
BenchmarkIsJson/PackageLockAxios/IsJson-8                                                    511           2331683 ns/op              16 B/op          0 allocs/op
BenchmarkIsJson/PackageLockAxios/EncodingJson-8                                               97          11834472 ns/op         5255954 B/op      97406 allocs/op
BenchmarkRedactAllValues/LongString/RedactAllValues-8                                      91222             13151 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongString/EncodingJson-8                                         28405             42304 ns/op           10440 B/op          5 allocs/op
BenchmarkRedactAllValues/LongNumber/RedactAllValues-8                                     123540              9704 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/RedactAllValues-8                                        91056             13176 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongName/EncodingJson-8                                           22549             52959 ns/op           21237 B/op         10 allocs/op
BenchmarkRedactAllValues/LongWhitespace/RedactAllValues-8                                 133742              8929 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongWhitespace/EncodingJson-8                                     25627             47297 ns/op             180 B/op          4 allocs/op
BenchmarkRedactAllValues/longArray/RedactAllValues-8                                       20440             58710 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/longArray/EncodingJson-8                                           2785            423017 ns/op          292788 B/op       5140 allocs/op
BenchmarkRedactAllValues/LongObject/RedactAllValues-8                                      37698             30898 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/LongObject/EncodingJson-8                                          2361            497345 ns/op          280723 B/op       3813 allocs/op
BenchmarkRedactAllValues/ManyArrays/RedactAllValues-8                                      30164             39730 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyArrays/EncodingJson-8                                          2931            395539 ns/op          334008 B/op       6847 allocs/op
BenchmarkRedactAllValues/ManyObjects/RedactAllValues-8                                     28626             41896 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyObjects/EncodingJson-8                                         1812            645269 ns/op          694352 B/op       8798 allocs/op
BenchmarkRedactAllValues/ManyTrues/RedactAllValues-8                                       38942             30758 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyTrues/EncodingJson-8                                           9067            123876 ns/op          112605 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyFalses/RedactAllValues-8                                      48184             24850 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyFalses/EncodingJson-8                                         10000            108082 ns/op          111229 B/op         19 allocs/op
BenchmarkRedactAllValues/ManyNulls/RedactAllValues-8                                       38722             30701 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValues/ManyNulls/EncodingJson-8                                          13690             87782 ns/op          112621 B/op         19 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/RedactAllValues-8                                  428           2785106 ns/op              20 B/op          0 allocs/op
BenchmarkRedactAllValues/PackageLockAxios/EncodingJson-8                                      58          19400187 ns/op         9068179 B/op     177202 allocs/op
BenchmarkRedactAllValuesSamplePackageLockAxiosReferenceImplementationPremarshalled-8         840           1445712 ns/op             631 B/op         26 allocs/op
PASS
ok      github.com/theteacat/jsonbytes  158.470s
```
