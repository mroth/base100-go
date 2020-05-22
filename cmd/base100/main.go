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

	flag.BoolVar(&opts.decode, "decode", false, "")
	flag.BoolVar(&opts.decode, "d", false, "")
	flag.StringVar(&opts.input, "input", "", "")
	flag.StringVar(&opts.input, "i", "", "")
	flag.StringVar(&opts.output, "output", "", "")
	flag.StringVar(&opts.output, "o", "", "")

	flag.Parse()
	return
}

func main() {
	opts := cliParse()

	if opts.input != "" || opts.output != "" {
		log.Fatal("TODO: not handling files just yet...")
	}
	in, out := os.Stdin, os.Stdout
	defer out.Close()

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	if opts.decode {
		// decoder currently dies due to lack of CRLF filtering
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
