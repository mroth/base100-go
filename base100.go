// Package base100 provides a Go implementation of Base💯.
//
// For the original Rust version see https://github.com/AdamNiederer/base100.
package base100

import (
	"errors"
	"io"
)

const (
	fixedByte1      = 0xf0 // 1st byte of an encoded base100 rune
	fixedByte2      = 0x9f // 2nd byte of an encoded base100 rune
	encodedByteSize = 4    // size of single "raw" byte after base100 encoding
)

/* ENCODE */

// Encode encodes src to its base100 encoding, writing EncodedLen(len(src))
// bytes to dst.
func Encode(dst, src []byte) {
	// We use this alternative loop and reslicing to help the compiler
	// perform bounds check elimination.
	//
	// Thanks to the #performance channel on the Gophers Slack!
	for len(dst) >= encodedByteSize && len(src) >= 1 {
		b := src[0]

		/* Rust version:
		out[4 * i + 0] = 0xf0;
		out[4 * i + 1] = 0x9f;
		// (ch + 55) >> 6 approximates (ch + 55) / 64
		out[4 * i + 2] = ((((*ch as u16).wrapping_add(55)) >> 6) + 143) as u8;
		// (ch + 55) & 0x3f approximates (ch + 55) % 64
		out[4 * i + 3] = (ch.wrapping_add(55) & 0x3f).wrapping_add(128);
		*/
		dst[0] = fixedByte1
		dst[1] = fixedByte2
		dst[2] = byte((uint16(b)+55)/64 + 143)
		dst[3] = (b+55)%64 + 128
		// dst[2] = byte(((uint16(b) + 55) >> 6) + 143)
		// dst[3] = (b+55)&0x3f + 128
		//
		// ^^ These bit shifting variations from Rust version don't appear to
		// impact speed at all when benchmarking here (in go1.14), let's assume
		// the Go compiler is smart enough to apply them already automatically?
		// Thus omit for readability.

		dst = dst[encodedByteSize:]
		src = src[1:]
	}
}

// EncodedLen returns the length in bytes of the base100 encoding of an input
// buffer of length n.
func EncodedLen(n int) int {
	return n * encodedByteSize
}

// EncodeToString returns the base100 encoding of src.
func EncodeToString(src []byte) string {
	buf := make([]byte, EncodedLen(len(src)))
	Encode(buf, src)
	return string(buf)
}

/* DECODE */

// Decode decodes src using base100. It writes at most DecodedLen(len(src))
// bytes to dst and returns the number of bytes written.
//
// New line characters (\r and \n) should be stripped beforehand.
func Decode(dst, src []byte) (n int, err error) {
	// if len(src)%4 != 0 {
	// 	return 0, errors.New("invalid length")
	// }
	/* we no longer need above check, as any trailing bytes will be sliced off
	anyhow during BCE hinting */

	/* Rust version:
	for (i, chunk) in buf.chunks(4).enumerate() {
	    out[i] = ((chunk[2].wrapping_sub(143)).wrapping_mul(64))
	        .wrapping_add(chunk[3].wrapping_sub(128)).wrapping_sub(55)
	}
	*/

	// Wacky compiler shenanigans for "undetected" bounds check elimination.
	//
	// When building with `go build -gcflags="-d=ssa/check_bce/debug=1"` this
	// reports BCE has not occured for src, e.g. it sees `IsInBounds` checks
	// occuring in the hoot loop, BUT the benchmarks nearly double from 550MB/s
	// to 1GB/s, so this is clearly doing BCE despite gcflags not reporting it?
	// Hence the "undetected" comment above. This lack of reporting may actually
	// be a bug in the Go compiler toolchain, we should check in next patch
	// version and file a bug if so.
	max := len(src) / encodedByteSize
	const employBCE = true
	if employBCE { // ^^ hard coded enabled above
		if len(dst) >= max && len(src) >= max*encodedByteSize {
			dst = dst[:max]                 // BCE hint!
			src = src[:max*encodedByteSize] // BCE hint!
		} else {
			// In the standard library implmentation of base64, if len(dst) <
			// DecodedLen(len(dst))), the method will panic. However, it seems
			// like we can be a bit more graceful here.
			return n, errors.New("insufficient slice size")
		}
	}

	for i := 0; i < max; i++ {
		offset := encodedByteSize * i
		// Optionally validate the prefix bytes are as expected on every rune.
		// This is quite wasteful for performance, and the Rust version of the
		// algorithm does not perform it, so leave it disabled here for now.
		const validate = false
		if validate { // ^^ hard coded disabled above
			pos1 := src[offset+0]
			pos2 := src[offset+1]
			if pos1 != fixedByte1 || pos2 != fixedByte2 {
				return n, errors.New("invalid encoding")
			}
		}

		pos3 := src[offset+2]
		pos4 := src[offset+3]
		dst[i] = (pos3-143)*64 + pos4 - 128 - 55
		n++
	}

	// Alternate version of the loop that employs via a different BCE technique
	// similar to Encode(), however it ends up being slower when benchmarked:
	//
	// 		for len(dst) >= 1 && len(src) >= 4 {
	// 			chunk := src[:4]
	// 			dst[0] = (chunk[2]-143)*64 + chunk[3] - 128 - 55
	// 			dst = dst[1:]
	// 			src = src[4:]
	// 		}

	return n, nil
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base100-encoded data.
func DecodedLen(n int) int {
	// max is no line breaks (e.g. every input is payload), easy to calculate
	return n / encodedByteSize
}

// DecodeString returns the bytes represented by the base100 string s.
func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	buf := make([]byte, DecodedLen(len(src)))
	_, err := Decode(buf, src)
	return buf, err
}

/* ENCODER */

// NewEncoder returns a new base100 stream encoder. Data written to the returned
// writer will be encoded using base100 and then written to w.
func NewEncoder(w io.Writer) io.Writer {
	return &encoder{w: w}
}

const bufferSize = 1024

type encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

func (e *encoder) Write(p []byte) (n int, err error) {
	/* (io.Writer).Write() notes:

	Write writes len(p) bytes from p to the underlying data stream. It returns
	the number of bytes written from p (0 <= n <= len(p)) and any error
	encountered that caused the write to stop early. Write must return a non-nil
	error if it returns n < len(p). Write must not modify the slice data, even
	temporarily.

	Implementations must not retain p.
	*/
	for len(p) > 0 && e.err == nil {
		chunkSize := bufferSize / encodedByteSize
		if len(p) < chunkSize {
			chunkSize = len(p)
		}

		chunk := p[:chunkSize]
		Encode(e.out[:], chunk)
		numBytesEncoded := EncodedLen(len(chunk))

		var written int
		written, e.err = e.w.Write(e.out[:numBytesEncoded])

		n += written / encodedByteSize
		p = p[chunkSize:]
	}
	return n, e.err
}

/* DECODER */

// NewDecoder constructs a new base100 stream decoder.
func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

type decoder struct {
	r   io.Reader
	err error
	in  []byte           // input buffer (encoded form)
	arr [bufferSize]byte // backing array for in
}

func (d *decoder) Read(p []byte) (n int, err error) {
	/* (io.Reader).Read() notes:

	Reader is the interface that wraps the basic Read method.

	Read reads up to len(p) bytes into p. It returns the number of bytes read (0
	<= n <= len(p)) and any error encountered. Even if Read returns n < len(p),
	it may use all of p as scratch space during the call. If some data is
	available but not len(p) bytes, Read conventionally returns what is
	available instead of waiting for more.

	When Read encounters an error or end-of-file condition after successfully
	reading n > 0 bytes, it returns the number of bytes read. It may return the
	(non-nil) error from the same call or return the error (and n == 0) from a
	subsequent call. An instance of this general case is that a Reader returning
	a non-zero number of bytes at the end of the input stream may return either
	err == EOF or err == nil. The next Read should return 0, EOF.

	Callers should always process the n > 0 bytes returned before considering
	the error err. Doing so correctly handles I/O errors that happen after
	reading some bytes and also both of the allowed EOF behaviors.

	Implementations of Read are discouraged from returning a zero byte count
	with a nil error, except when len(p) == 0. Callers should treat a return of
	0 and nil as indicating that nothing happened; in particular it does not
	indicate EOF.

	Implementations must not retain p.
	*/

	// Some references for comparison and ideas:
	// 	- https://golang.org/src/encoding/base64/base64.go
	// 	- https://golang.org/src/encoding/hex/hex.go

	// assume \r \n stripped via stripper...

	// Fill internal buffer with sufficient bytes to decode
	if len(d.in) < encodedByteSize && d.err == nil {
		var numCopy, numRead int
		numCopy = copy(d.arr[:], d.in)             // Copies any remainder bytes [0-3] from before into beginning of backing array
		numRead, d.err = d.r.Read(d.arr[numCopy:]) // read from internal reader with slice of rest of remaining internal buffer
		d.in = d.arr[:numCopy+numRead]             // reset in to resliced arr containing all data

		// handle case: we got an EOF but the bytes we have in our internal
		// buffer are not a proper multiple of the encodedByteSize.
		if d.err == io.EOF && len(d.in)%encodedByteSize != 0 {
			// TODO: hex.go actually checks to see if there is an invalid
			// encoding byte first, should we do something similar?
			d.err = io.ErrUnexpectedEOF
		}
	}

	if numDecodedBytesAvail := len(d.in) / encodedByteSize; len(p) > numDecodedBytesAvail {
		p = p[:numDecodedBytesAvail] // reslice p to be only size needed...(why?)
	}

	numDecodedBytes, err := Decode(p, d.in[:len(p)*encodedByteSize]) // decode into p
	d.in = d.in[encodedByteSize*numDecodedBytes:]                    // reslice in to remainder

	// if decode error; discard input remainder & bubble up error
	if err != nil {
		d.in, d.err = nil, err
	}

	// only expose errors when buffer fully consumed
	if len(d.in) < encodedByteSize {
		return numDecodedBytes, d.err
	}

	return numDecodedBytes, nil
}

// NewlineFilteringReader wraps an io.Reader and strips CRLF from the stream.
//
// This implementation is taken verbatim from the base64 module of the Go
// standard library (where it is not exported, hence the copypasta).
//
// Currently this is not utilized the decoder by default, since it hits
// performance very hard. Leaving here for now as a hint for the future.

/*
type NewlineFilteringReader struct {
	wrapped io.Reader
}

func (r *NewlineFilteringReader) Read(p []byte) (int, error) {
	n, err := r.wrapped.Read(p)
	for n > 0 {
		offset := 0
		for i, b := range p[:n] {
			if b != '\r' && b != '\n' {
				if i != offset {
					p[offset] = b
				}
				offset++
			}
		}
		if offset > 0 {
			return offset, err
		}
		// Previous buffer entirely whitespace, read again
		n, err = r.wrapped.Read(p)
	}
	return n, err
}
*/
