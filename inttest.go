package main

import (
	"fmt"
	"bytes"
	"encoding/binary"
)

func main () {
	var val int32
	val = 8553090
	buf := new (bytes.Buffer)
	
	err := binary.Write (buf, binary.BigEndian, val)
	
	if err != nil {
		fmt.Println ("Binary encoding failed\n", err)
	}
	
	fmt.Println (buf.Bytes())
}
