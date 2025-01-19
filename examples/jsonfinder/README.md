# jsonfinder

jsonfinder is a CLI tool that searches your filesystem for valid JSON files. Its purpose here is simply to demonstrate usage of the jsonbytes package.

If anyone ever finds this genuinely useful, please let me know. I'd be amazed!



## Quickstart

```bash
git clone git@github.com:TheTeaCat/jsonbytes.git
cd jsonbytes/examples/jsonfinder
go build
./jsonfinder -help
```



## Example Usage

```bash
./jsonfinder -dir=../../testdata
```

```
2025/01/19 00:36:44 ✅ ../../testdata/package-lock-axios.json is JSON!
2025/01/19 00:36:44 Every file checked was valid JSON! 🥳
```



