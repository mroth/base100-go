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

var (
	decode = flag.Bool("decode", false, "decode data")
)

func main() {
	flag.Parse()

	in, out := os.Stdin, os.Stdout
	if flag.Arg(0) != "" {
		log.Fatal("not handling <INPUT> yet...")
	}
	defer out.Close()

	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	if *decode {
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
