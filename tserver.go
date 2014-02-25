package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"runtime"
	"time"
)

type MSG struct {
	otype int32
	mlen  int32
	m     []byte
}

var tps uint64

func printstats() {
	for {
		select {
		case <-time.After(1 * time.Second):
			fmt.Printf("tps = %v\n", tps)
			tps = 0

		}
	}
}

func DispStts(s chan int) {

	go printstats()
	for {
		select {
		case <-s:
			tps++
		}
	}
}

/*

func ReadN (c net.Conn, n int) ([]byte, error, int) {
        msg := make ([]byte, n)

        cnt , err := c.Read(msg)
        if err != nil {
                return msg,err, cnt
        }
        if cnt == n {
                return msg, nil, cnt
        }
        tcnt := cnt
        for tcnt == n {
                cnt, err := c.Read(msg[tcnt:])
                if err != nil {
                        return msg,err, tcnt
                }
                tcnt += cnt
        }
        return msg,nil,tcnt

}
*/

func HandleConnection(c net.Conn, s chan int) {

	for {
		var v MSG

		var m []byte

		m = make([]byte, 40)

		debug := false

		cnt, err := io.ReadAtLeast(c, m, 40)

		if err != nil {
			fmt.Println("ReadAtleast :", err)
			break
		}

		if cnt != 40 {
			fmt.Println("Bad message len:", cnt)
			break
		}

		if debug {
			fmt.Printf("Recived %d bytes %v\n", cnt, err)
		}

		bf := bytes.NewBuffer(m)

		binary.Read(bf, binary.BigEndian, &(v.otype))

		binary.Read(bf, binary.BigEndian, &(v.mlen))

		s <- 1
		if debug {
			fmt.Printf("Message Type=%d Mlen=%d %s\n", v.otype, v.mlen, string(m[8:cnt]))
		}

	}

	c.Close()

}

func main() {

	l, err := net.Listen("tcp", ":6060")
	if err != nil {
		fmt.Println("Listen Failed:", err)
		panic("Server failed to bind:")
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	s := make(chan int)

	go DispStts(s)

	for {

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Accep failed:", err)
		}

		go HandleConnection(conn, s)

	}

}
