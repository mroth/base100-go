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

	encoder := base100.NewEncoder(writer)
	_, err := io.Copy(encoder, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
		os.Exit(1)
	}
}
