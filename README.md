# Base💯

A Go implementation of [base100](https://github.com/AdamNiederer/base100).

Base💯 can represent any byte with a unique emoji symbol, therefore it can
represent binary data with zero printable overhead.

## Usage

### Library

See the Go Docs for more information.

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

TODO

## License

TODO
