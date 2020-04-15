// Package base100 provides a Go implementation of base100.
//
// Potential differences from the Rust implementation of base64.
// - Like base64,
package base100

/* ENCODE */

// Encode encodes src to its base100 encoding, writing EncodedLen(len(src))
// bytes to dst.
func Encode(dst, src []byte) {
	// we use this alternative loop and reslicing to help the compiler
	// perform bounds check elimination.
	//
	// thanks to the #performance channel on the gophers slack!
	for len(dst) >= 4 && len(src) >= 1 {
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
		dst[2] = byte((uint16(b)+55)/64 + 143)
		dst[3] = (b+55)%64 + 128

		dst = dst[4:]
		src = src[1:]
	}
}

// EncodedLen returns the length in bytes of the base100 encoding of an input
// buffer of length n.
func EncodedLen(n int) int {
	return n * 4
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
	max := len(src) / 4
	for i := 0; i < max; i++ {
		// if checked, verify first and second position?
		// pos1 := src[4*i+0]
		// pos2 := src[4*i+1]
		pos3 := src[4*i+2]
		pos4 := src[4*i+3]
		dst[i] = (pos3-143)*64 + pos4 - 128 - 55
	}

	// our only option before this point is panic, so we dont need to keep track
	// of bytes written
	return max, nil
}

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base100-encoded data.
func DecodedLen(n int) int {
	// max would be no line breaks (e.g. every input is payload), so easy to
	// calculate.
	//
	// if input is malformed length, doesnt matter here? since wont result in
	// needing that space....
	return n / 4
}

// DecodeString returns the bytes represented by the base100 string s.
func DecodeString(s string) ([]byte, error) {
	src := []byte(s)
	buf := make([]byte, DecodedLen(len(src)))
	_, err := Decode(buf, src)
	return buf, err
}

// func NewEncoder(w io.Writer) io.WriteCloser
