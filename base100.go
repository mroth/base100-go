// Package base100 provides a Go implementation of base100.
//
// Potential differences from the Rust implementation of base64.
// - Like base64,
package base100

import (
	"io"
)

const (
	encodedByteSize = 4 // size of single "raw" byte after base100 encoding
)

/* ENCODE */

// Encode encodes src to its base100 encoding, writing EncodedLen(len(src))
// bytes to dst.
func Encode(dst, src []byte) {
	// we use this alternative loop and reslicing to help the compiler
	// perform bounds check elimination.
	//
	// thanks to the #performance channel on the gophers slack!
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
		dst[0] = 0xf0
		dst[1] = 0x9f
		// These bit shifting variations from Rust version don't appear to
		// impact speed at all when benchmarking here, likely the Go compiler
		// is smart enough to apply them already automatically?
		//
		// dst[2] = byte(((uint16(b) + 55) >> 6) + 143)
		// dst[3] = (b+55)&0x3f + 128
		dst[2] = byte((uint16(b)+55)/64 + 143)
		dst[3] = (b+55)%64 + 128

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
// bytes to dst and returns the number of bytes written. If src contains invalid
// base100 data, it will return the number of bytes successfully written and
// CorruptInputError. New line characters (\r and \n) are ignored.
//
// TODO: dont bother to check for invalid data? or have different versions?
// TODO: strip new lines (optional?)
//
// TODO: what is expected behavior if len(dst) < DecodedLen(len(dst)))? base64 will panic.
func Decode(dst, src []byte) (n int, err error) {
	// if len(data)%4 != 0 {
	// 	return nil, ErrInvalidLength
	// }

	/* Rust version:
	for (i, chunk) in buf.chunks(4).enumerate() {
	    out[i] = ((chunk[2].wrapping_sub(143)).wrapping_mul(64))
	        .wrapping_add(chunk[3].wrapping_sub(128)).wrapping_sub(55)
	}
	*/
	max := len(src) / encodedByteSize
	for i := 0; i < max; i++ {
		// if checked, verify first and second position?
		// pos1 := src[4*i+0]
		// pos2 := src[4*i+1]
		pos3 := src[encodedByteSize*i+2]
		pos4 := src[encodedByteSize*i+3]
		dst[i] = (pos3-143)*64 + pos4 - 128 - 55
	}

	/*
		Version of  the loop that employs BCE similar to Encode(), however it
		actually ends up being slower when benchmarked.

			for len(dst) >= 1 && len(src) >= 4 {
				chunk := src[:4]
				dst[0] = (chunk[2]-143)*64 + chunk[3] - 128 - 55
				n++
				dst = dst[1:]
				src = src[4:]
			}
			return n, nil
	*/

	// our only option before this point is panic, so we dont need to keep track
	// of bytes written
	return max, nil
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

// NewEncoder returns a new base100 stream encoder. Data written to
// the returned writer will be encoded using base100 and then written to w.
//
// TODO: figure out the below in our world... NOT NEEDED
// Base64 encodings operate in 4-byte blocks; when finished
// writing, the caller must Close the returned encoder to flush any
// partially written blocks.
func NewEncoder(w io.Writer) io.Writer {
	return &encoder{w: w}
}

const bufferSize = 1024

type encoder struct {
	w   io.Writer
	err error
	out [bufferSize]byte // output buffer
}

// Write implements io.Writer.
func (e *encoder) Write(p []byte) (n int, err error) {
	// Write writes len(p) bytes from p to the underlying data stream. It
	// returns the number of bytes written from p (0 <= n <= len(p)) and any
	// error encountered that caused the write to stop early. Write must return
	// a non-nil error if it returns n < len(p). Write must not modify the slice
	// data, even temporarily.
	//
	// Implementations must not retain p.

	// A good reference here is https://golang.org/src/encoding/hex/hex.go
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
type decoder struct {
	r   io.Reader
	err error
	in  []byte           // input buffer (encoded form)
	arr [bufferSize]byte // backing array for in
}

func (d *decoder) Read(p []byte) (n int, err error) {
	/* Reader is the interface that wraps the basic Read method.

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

	// Modeling this one after https://golang.org/src/encoding/base64/base64.go
	// Hmm https://golang.org/src/encoding/hex/hex.go is better actually?

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
			// TODO: hex.go actually checks to see if there is an invalid encoding byte first, do something similar?
			d.err = io.ErrUnexpectedEOF
		}
	}

	if numDecodedBytesAvail := len(d.in) / encodedByteSize; len(p) > numDecodedBytesAvail {
		p = p[:numDecodedBytesAvail] // reslice p to be only size needed...(why?)
	}

	numDecodedBytes, err := Decode(p, d.in[:len(p)*encodedByteSize]) // decode into p
	d.in = d.in[encodedByteSize*numDecodedBytes:]                    // reslice in to remainder
	// decode error; discard input remainder & bubble up error
	if err != nil {
		d.in, d.err = nil, err
	}
	// only expose errors when buffer fully consumed
	if len(d.in) < encodedByteSize {
		return numDecodedBytes, d.err
	}

	return numDecodedBytes, nil
}

func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

// func NewDecoder(r io.Reader, validated bool) io.Reader
