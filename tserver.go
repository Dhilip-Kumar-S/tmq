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


func HandleConnection (c net.Conn) {
	
	var v MSG
	
	var m []byte
	
	m = make([]byte, 1024)
	
	cnt, err := c.Read(m)
	
	fmt.Printf("Recived %d bytes %v\n", cnt, err)
	
	bf := bytes.NewBuffer(m)
	
	binary.Read(bf, binary.BigEndian, &(v.otype))
	
	binary.Read(bf, binary.BigEndian, &(v.mlen))
	
	fmt.Printf("Message Type=%d Mlen=%d %s\n", v.otype, v.mlen, string(m[8:cnt]))
	
	c.Close ()
	
}

func main () {
	
	l, err := net.Listen("tcp", "127.0.0.1:6060")
	if err != nil {
		fmt.Println ("Listen Failed:", err)
		panic ("Server failed to bind:")
	}
	fmt.Printf("Listen sucessful:", l)
	
	for {
		
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accep failed:", err)
		}
		
		go HandleConnection (conn)
		
	}
	
}



