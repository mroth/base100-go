package base100_test

import (
	"fmt"
	"log"

	"github.com/mroth/base100-go"
)

func ExampleEncode() {
	src := []byte("the quick brown fox jumped over the lazy dog\n")
	dst := make([]byte, base100.EncodedLen(len(src)))
	base100.Encode(dst, src)
	fmt.Printf("%s", dst)
	// Output: 👫👟👜🐗👨👬👠👚👢🐗👙👩👦👮👥🐗👝👦👯🐗👡👬👤👧👜👛🐗👦👭👜👩🐗👫👟👜🐗👣👘👱👰🐗👛👦👞🐁
}

func ExampleEncodeToString() {
	src := []byte("the quick brown fox jumped over the lazy dog\n")
	fmt.Println(base100.EncodeToString(src))
	// Output: 👫👟👜🐗👨👬👠👚👢🐗👙👩👦👮👥🐗👝👦👯🐗👡👬👤👧👜👛🐗👦👭👜👩🐗👫👟👜🐗👣👘👱👰🐗👛👦👞🐁
}

func ExampleDecode() {
	src := []byte("👫👟👜🐗👨👬👠👚👢🐗👙👩👦👮👥🐗👝👦👯🐗👡👬👤👧👜👛🐗👦👭👜👩🐗👫👟👜🐗👣👘👱👰🐗👛👦👞🐁")
	dst := make([]byte, base100.DecodedLen(len(src)))
	_, err := base100.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", dst)
	// Output: the quick brown fox jumped over the lazy dog
}

func ExampleDecodeString() {
	src := "👫👟👜🐗👨👬👠👚👢🐗👙👩👦👮👥🐗👝👦👯🐗👡👬👤👧👜👛🐗👦👭👜👩🐗👫👟👜🐗👣👘👱👰🐗👛👦👞🐁"
	result, _ := base100.DecodeString(src)
	fmt.Printf("%s", result)
	// Output: the quick brown fox jumped over the lazy dog
}
