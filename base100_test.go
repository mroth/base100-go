// Package base100 provides a Go implementation of base100
package base100

import (
	"fmt"
	"reflect"
	"testing"
)

var samplecases = []struct {
	data []byte // raw data
	text []byte // encoded version
}{
	{
		[]byte("the quick brown fox jumped over the lazy dog\n"),
		[]byte("👫👟👜🐗👨👬👠👚👢🐗👙👩👦👮👥🐗👝👦👯🐗👡👬👤👧👜👛🐗👦👭👜👩🐗👫👟👜🐗👣👘👱👰🐗👛👦👞🐁"),
	},
}

func TestEncode(t *testing.T) {
	for n, sc := range samplecases {
		t.Run(fmt.Sprintf("sample%02d", n), func(t *testing.T) {
			src := sc.data
			dst := make([]byte, EncodedLen(len(src)))
			Encode(dst, src)
			if got, want := dst, sc.text; !reflect.DeepEqual(got, want) {
				t.Errorf("Encode() = %v, want %v", got, want)
			}
		})
	}
}

func TestEncodeToString(t *testing.T) {
	for n, sc := range samplecases {
		t.Run(fmt.Sprintf("sample%02d", n), func(t *testing.T) {
			want := string(sc.text)
			if got := EncodeToString(sc.data); got != want {
				t.Errorf("EncodeToString() = %v, want %v", got, want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	// handle all "normal" sample cases
	for n, sc := range samplecases {
		t.Run(fmt.Sprintf("sample%02d", n), func(t *testing.T) {
			src := sc.text
			dst := make([]byte, DecodedLen(len(src)))

			n, err := Decode(dst, src)
			var (
				wantN   = len(sc.data)
				wantErr = false
				wantDst = sc.data
			)
			if n != wantN {
				t.Errorf("n = %v, want %v", n, wantN)
			}
			if (err != nil) != wantErr {
				t.Errorf("err = %v, wantError: %v", err, wantErr)
			}
			if !reflect.DeepEqual(dst, wantDst) {
				t.Errorf("dst = %v, want %v", dst, wantDst)
			}
		})
	}

	// TODO: error/invalid cases....
	// TODO: expect a panic when dst too small
}

func TestDecodeString(t *testing.T) {
	// handle all "normal" sample cases
	for n, sc := range samplecases {
		t.Run(fmt.Sprintf("sample%02d", n), func(t *testing.T) {
			src := string(sc.text)
			got, err := DecodeString(src)
			if wantErr := false; (err != nil) != wantErr {
				t.Errorf("err = %v, wantError: %v", err, wantErr)
			}
			if want := sc.data; !reflect.DeepEqual(got, want) {
				t.Errorf("result = %v, want %v", got, want)
			}
		})
	}
}

var (
	benchdata = samplecases[0].data
	benchtext = samplecases[0].text
)

func BenchmarkEncode(b *testing.B) {
	src := benchdata
	dst := make([]byte, EncodedLen(len(src)))
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Encode(dst, src)
	}
}

func BenchmarkEncodeToString(b *testing.B) {
	src := benchdata
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeToString(src)
	}
}

func BenchmarkDecode(b *testing.B) {
	src := benchtext
	dst := make([]byte, DecodedLen(len(src)))
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decode(dst, src)
	}
}

func BenchmarkDecodeString(b *testing.B) {
	src := string(benchtext)
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodeString(src)
	}
}
