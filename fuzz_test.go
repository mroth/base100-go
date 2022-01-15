//go:build go1.18
// +build go1.18

package base100

import (
	"bytes"
	"testing"
	"unicode/utf8"
)

func FuzzEncode(f *testing.F) {
	var fuzzcases = [][]byte{
		[]byte("Hello, world"),
		[]byte("Hello, world\n"),
		[]byte("\nHello, world"),
		[]byte(""),
		[]byte(" "),
		[]byte("\n"),
		[]byte("\r\n"),
		[]byte("!12345"),
		[]byte("제주도"),
	}

	for _, tc := range fuzzcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig []byte) {
		encoded := make([]byte, EncodedLen(len(orig)))
		Encode(encoded, orig)
		if got, want := len(encoded), EncodedLen(len(orig)); got != want {
			t.Errorf("encoded is %d bytes, but EncodedLen predicted %d", got, want)
		}
		if !utf8.Valid(encoded) { // encoded version should always be valid utf8
			t.Errorf("Encode produced invalid UTF-8 %q", encoded)
		}

		decoded := make([]byte, DecodedLen(len(encoded)))
		n, err := Decode(decoded, encoded)
		if err != nil {
			t.Errorf("Decode returned an error: %v", err)
		}
		if got, want := len(decoded), n; got != want {
			t.Errorf("dst is %d bytes, but Decode wrote %d", got, want)
		}
		if got, want := len(decoded), DecodedLen(len(encoded)); got != want {
			t.Errorf("dst is %d bytes, but DecodedLen predicted %d", got, want)
		}

		if !bytes.Equal(orig, decoded) {
			t.Errorf("original: %q, decoded: %q", orig, decoded)
		}
	})
}
