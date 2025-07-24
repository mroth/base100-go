# Base💯

A Go implementation of [base100](https://github.com/AdamNiederer/base100) with a
permissive software license.

Base💯 can represent any byte with a unique emoji symbol, therefore it can
represent binary data with zero printable overhead.

## Usage

### Library

The API is nearly identical to other modules from the Go `encoding/*` standard
library. See the [Go Docs](https://pkg.go.dev/github.com/mroth/base100-go) for
more information.

### Command line tool

A CLI tool is also provided for convenience, and ease of cross-compilation for
multiple operating systems.

    base💯 (Go)
    Encodes things into emoji

    USAGE:
        base100 [FLAGS]

    FLAGS:
        -d, --decode     Decodes input
        -i, --input      Input file (default use STDIN)
        -o, --output     Output file (default use STDOUT)
        -h, --help       Prints help information

`base100` will read from stdin unless a file is specified, will write UTF-8 to
stdout, and has a similar API to GNU's base64. Data is encoded by default,
unless `--decode` is specified

## Performance

The implementation is fairly performant, and appears to perform roughly
equivalent to the scalar Rust version on my machine. A future optimization
utilizing SIMD/AVX could be possible with Go assembly code, however the
throughput already far exceeds any known use case I can see for this, so I'll
leave that out unless I get incredibly bored some day.

Library single-cpu benchmarks from my laptop (the throughput values are the
relevant ones):
```
$ go test -bench=. -cpu=1
goos: darwin
goarch: arm64
pkg: github.com/mroth/base100-go
cpu: Apple M2 Pro
BenchmarkEncode                 24076875                49.24 ns/op      913.81 MB/s
BenchmarkEncodeToString         12909928                90.94 ns/op      494.83 MB/s
BenchmarkDecode                 31709400                37.85 ns/op     1189.02 MB/s
BenchmarkDecodeString           25889036                45.45 ns/op      990.11 MB/s
BenchmarkEncoder                   18339             65083 ns/op        1006.96 MB/s
BenchmarkDecoder                   21528             56465 ns/op        1160.65 MB/s
```
