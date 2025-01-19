# jsonbytes

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
pkg: jsonbytes
cpu: Apple M1 Pro
BenchmarkIsJsonLongString10KiB-8                                                           90062             13198 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonLongNumber10KiB-8                                                          123435              9733 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonLongName10KiB-8                                                            141166              8450 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonLongWhitespace10KiB-8                                                      134388              8929 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonLongArray10KiB-8                                                            24696             48613 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonLongObject10KiB-8                                                           73126             16373 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonManyArrays10KiB-8                                                           38200             31440 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonManyObjects10KiB-8                                                          39963             30163 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonManyTrues10KiB-8                                                            45212             26558 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonManyFalses10KiB-8                                                           55846             21438 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonManyNulls10KiB-8                                                            45046             26610 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesLongString10KiB-8                                                  90817             13186 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesLongNumber10KiB-8                                                 123235              9724 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesLongName10KiB-8                                                    90542             13223 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesLongWhitespace10KiB-8                                             132949              8967 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesLongArray10KiB-8                                                   20362             58941 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesLongObject10KiB-8                                                  48667             24643 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesManyArrays10KiB-8                                                  29868             40145 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesManyObjects10KiB-8                                                 29346             40874 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesManyTrues10KiB-8                                                   38766             30999 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesManyFalses10KiB-8                                                  48090             25010 ns/op               0 B/op          0 allocs/op
BenchmarkRedactAllValuesManyNulls10KiB-8                                                   38714             30960 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonSamplePackageLockAxios-8                                                      511           2348059 ns/op               0 B/op          0 allocs/op
BenchmarkIsJsonSamplePackageLockAxiosReferenceImplementation-8                                97          11796296 ns/op         5257267 B/op      97409 allocs/op
BenchmarkRedactAllValuesSamplePackageLockAxios-8                                             424           2817082 ns/op              20 B/op          0 allocs/op
BenchmarkRedactAllValuesSamplePackageLockAxiosReferenceImplementation-8                       57          19392148 ns/op         8894840 B/op     177199 allocs/op
BenchmarkRedactAllValuesSamplePackageLockAxiosReferenceImplementationPremarshalled-8         840           1418088 ns/op             634 B/op         26 allocs/op
PASS
ok      jsonbytes       116.554s
```
