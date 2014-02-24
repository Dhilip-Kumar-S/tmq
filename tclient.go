package main

import (
"fmt"
"net"
"encoding/binary"
"bytes"
)

type MSG struct {
	otype	int32
	mlen	int32
	m		[]byte	
}

func main () {

	var err error
	c, err := net.Dial ("tcp", "[127.0.0.1]:6060")
	
	if err != nil {
		fmt.Println("Unable to Dial:", err)
		panic ("Unable to connect")
	}
	
	var m MSG
	
	m.otype = 5
	m.mlen = 32
	m.m = make ([]byte, 32)
	
	for i:=0; i<32; i++ {
		m.m[i] = 0x41
	}
	
	buf := new (bytes.Buffer)
	err = binary.Write (buf, binary.BigEndian, m.otype)
	err = binary.Write (buf, binary.BigEndian, m.mlen)
	m.m = append (buf.Bytes(), m.m...)
	c.Write (m.m)
	c.Close()
}