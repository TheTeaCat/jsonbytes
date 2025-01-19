# jsonfinder

jsonfinder is a CLI tool that searches your filesystem for valid JSON files. Its purpose here is simply to demonstrate usage of the jsonbytes package.

If anyone ever finds this genuinely useful, please let me know. I'd be amazed!



## Quickstart

```bash
git clone git@github.com/theteacat/jsonbytes
cd jsonbytes/examples/jsonfinder
go build
./jsonfinder -help
```



## Example Usage

```bash
./jsonfinder -dir=../../testdata
```

```
✅ ../../testdata/package-lock-axios.json is JSON!
```



