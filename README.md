# Base💯

A Go implementation of [base100](https://github.com/AdamNiederer/base100).

Base💯 can represent any byte with a unique emoji symbol, therefore it can
represent binary data with zero printable overhead.

## Usage

### Library

The API is nearly identical to other modules from the Go `encoding/*` standard
library. See the [Go Docs](https://godoc.org/github.com/mroth/base100-go) for
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
equivalent to the standard Rust version (e.g. non-AVX) on my machine. A future
optimization utilizing SIMD/AVX could be possible with Go assembly code, however
the throughput already far exceeds any known use case I can see for this, so
I'll leave that out unless I get incredibly bored some day.

Library benchmarks from my laptop (the throughput values as the relevant ones):
```
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/mroth/base100-go
BenchmarkEncode-8               15173618                72.9 ns/op       617.19 MB/s
BenchmarkEncodeToString-8        7025580               167 ns/op         269.30 MB/s
BenchmarkDecode-8               25826598                41.5 ns/op      1084.97 MB/s
BenchmarkDecodeString-8          9466012               121 ns/op         371.00 MB/s
BenchmarkEncoder-8                  6042            178455 ns/op         560.37 MB/s
BenchmarkDecoder-8                 11887             97977 ns/op        1020.65 MB/s
```

## License

TODO
