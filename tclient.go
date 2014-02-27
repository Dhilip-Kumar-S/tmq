package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
//"time"
)

type MSG struct {
	otype int32
	mlen  int32
	m     []byte
}

func main() {

	var ip string
	if len(os.Args) > 1 {
		ip = os.Args[1]
	} else {
		ip = "127.0.0.1"
	}
	c, err := net.Dial("tcp", "["+ip+"]:6060")

	if err != nil {
		fmt.Println("Unable to Dial:", err)
		panic("Unable to connect")
	}

	var m MSG

	m.otype = 5000
	m.mlen = 32
	m.m = make([]byte, 32)

	for i := 0; i < 32; i++ {
		m.m[i] = 0x41
	}
	for {

		buf := new(bytes.Buffer)
		err = binary.Write(buf, binary.BigEndian, m.otype)
		err = binary.Write(buf, binary.BigEndian, m.mlen)
		smg := append(buf.Bytes(), m.m...)
		c.Write(smg)
		//time.Sleep (time.Millisecond)
	}
	c.Close()
}
