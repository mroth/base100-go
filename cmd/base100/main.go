package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mroth/base100-go"
)

const (
	productShortName = "base💯"
	productFullName  = "base💯 (Go)"
)

type options struct {
	decode        bool   // decode input instead of encode
	input, output string // optional file paths
}

func cliParse() (opts options) {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `%s
Encodes things into emoji

USAGE:
    base100 [FLAGS]

FLAGS:
    -d, --decode     Decodes input
    -i, --input      Input file (default use STDIN)
    -o, --output     Output file (default use STDOUT)
    -h, --help       Prints help information
`, productFullName)
	}

	const nodesc = "" // descriptions not shown since we override flag.Usage
	flag.BoolVar(&opts.decode, "decode", false, nodesc)
	flag.BoolVar(&opts.decode, "d", false, nodesc)
	flag.StringVar(&opts.input, "input", "", nodesc)
	flag.StringVar(&opts.input, "i", "", nodesc)
	flag.StringVar(&opts.output, "output", "", nodesc)
	flag.StringVar(&opts.output, "o", "", nodesc)

	flag.Parse()
	return
}

func main() {
	opts := cliParse()

	in := os.Stdin
	if opts.input != "" {
		var err error
		in, err = os.Open(opts.input)
		if err != nil {
			log.Fatal(err)
		}
		defer in.Close()
	}

	out := os.Stdout
	if opts.output != "" {
		var err error
		out, err = os.Create(opts.output)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	bufsize := 1024 * 64 // 64 KiB io buffers, same as pipe buffer in linux >=2.6.11
	reader := bufio.NewReaderSize(in, bufsize)
	writer := bufio.NewWriterSize(out, bufsize)
	defer writer.Flush()

	if opts.decode {
		// decoder currently can die due to lack of CRLF filtering
		decoder := base100.NewDecoder(reader)
		_, err := io.Copy(writer, decoder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
			os.Exit(1)
		}
	} else {
		encoder := base100.NewEncoder(writer)
		_, err := io.Copy(encoder, reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
			os.Exit(1)
		}
	}
}
