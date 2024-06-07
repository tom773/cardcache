package pkg

import (
	"bytes"
	"fmt"
	"reflect"
)

func lex(raw []byte) {
	if raw[0] == '[' {
		arr := bytes.Split(raw[1:len(raw)-1], []byte(","))
		fmt.Println(reflect.TypeOf(arr), arr)
	} else if raw[0] == '{' {
		fmt.Println("object")
	} else {
		fmt.Println("string")
	}
}
