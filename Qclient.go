package main

import (
	//	"bytes"
	//	"encoding/binary"
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"tcp"
)

const (
	CREATE = iota // 0
	DELETE        // 1
	OPEN          // 2
	CLOSE         // 3
	ENQ           // 4
	DEQ           // 5
	TS            // 6
	TE            // 7
	TA            // 8
	SELECT        // 9
)

var (
	COUNT, MCOUNT   int
	server  string
	idlist  []string
	qbuffer map[string]list.List
	verbose bool
	sconn   net.Conn
)

func vprintf(format string, v ...interface{}) {

	if verbose {
		log.Printf(format, v...)
	}
}

func testCREATE() bool {

	vprintf("Creating Q\n")
	for i := 0; i < COUNT; i++ {
		cnt, err := tcp.WriteBYTE(sconn, CREATE)
		if err != nil {
			vprintf("Writebyte error wrote %dbytes error %v\n", cnt, err)
		}

		qname_str := "TEMPQ" + strconv.Itoa(i)
		qname := []byte(qname_str)

		cnt, err = tcp.WriteINT32(sconn, int32(len(qname)))
		if err != nil {
			vprintf("Writebyte error wrote %dbytes error %v\n", cnt, err)
		}

		cnt, err = tcp.WriteBytes(sconn, qname)
		if err != nil {
			vprintf("Writebyte error wrote %dbytes error %v\n", cnt, err)
		}

		ack, ok := tcp.ReadBYTE(sconn)
		if ok != nil {
			vprintf("Error Creating Q %v error %v\n", qname, ok)
			return false
		}

		if ack == 0x01 {
			vprintf("Q %s Created\n", qname_str)
		} else {
			vprintf("Q %s Already Exisist ACK=%v\n", qname_str, ack)
		}
	}

	return true
}

func testOPEN() bool {

	for i := 0; i < COUNT; i++ {
		cnt, err := tcp.WriteBYTE(sconn, OPEN)
		if err != nil {
			vprintf("WriteBYTE error wrote %dbytes error %v\n", cnt, err)
		}

		qname_str := "TEMPQ" + strconv.Itoa(i)
		qname := []byte(qname_str)

		cnt, err = tcp.WriteINT32(sconn, int32(len(qname)))
		if err != nil {
			vprintf("WriteINT32 error wrote %dbytes error %v\n", cnt, err)
		}

		cnt, err = tcp.WriteBytes(sconn, qname)
		if err != nil {
			vprintf("Writebyte error wrote %dbytes error %v\n", cnt, err)
		}

		id, err := tcp.ReadMQID(sconn)
		if err != nil {
			vprintf("ReadMQID error wrote %dbytes error %v\n", cnt, err)
		} else {
			if strings.Contains(id, "<NIL>") {
				vprintf("OPEN Q:%s FAILED recived %s\n", qname_str, id)
			} else {
				vprintf("OPEN Q: %s SUCESS\n", qname_str)
				idlist[i] = id
			}

		}

	}
	return true
}

func testENQ() bool {

	for i, id := range idlist {
		
		for j := 0; j < MCOUNT; j++ {
		
			strmsg := "QMESSAGEQ=" + strconv.Itoa(i) + ":" + strconv.Itoa(j)
			bytemsg := []byte(strmsg)
			lenmsg := len(bytemsg)
			tcp.WriteBYTE(sconn, ENQ)
			tcp.WriteMQID(sconn, id)
			tcp.WriteINT32(sconn, int32(lenmsg))
			tcp.WriteBytes(sconn, bytemsg)
			b, err := tcp.ReadBYTE(sconn)
			if err != nil {
				vprintf("READ ERROR ENQUE FAILED\n")
				return false
			}
			if b == 0x00 {
				vprintf("ENQ() sucess sent=%s\n", strmsg )
			} else {
				vprintf("ENQ () FAILED\n")
			}
		}
	}
	return true
}

func testDQ() bool {

	for _, id := range idlist {
		for {
			
			tcp.WriteBYTE(sconn, DEQ)
			tcp.WriteMQID(sconn, id)
			mlen, err := tcp.ReadINT32(sconn)
			if err != nil {
				vprintf("DeQ() error reading len %v\n", err)
				return false
			}
			if mlen == -1 {
				vprintf("Queue is empty\n")
				break
			} else {
				msg, err := tcp.ReadNBytes(sconn, mlen)
				if err != nil {
					vprintf("DeQ() error reading msg %v\n", err)
					return false
				} else {
					vprintf("DeQ() recived %s\n", string(msg))
				}
			}
		}

	}

	return true
}

func runTest(test func() bool, val string) {
	if test() {
		log.Printf("[%s]\t...\t...\t...\t...[OK]", val)
	} else {
		log.Printf("[%s]\t...\t...\t...\t...[FAIL]", val)
	}
}

func testAll() bool {

	var err error
	/* Connect to the server */
	sconn, err = net.Dial("tcp", server)

	if err != nil {
		vprintf("Error Connecting %v\n", err)
		return false
	} else {
		vprintf("Connected to Server %v\n", sconn)
	}

	/* Run the tests */
	runTest(testCREATE, "testCREATE")
	runTest(testOPEN, "testOPEN")
	runTest(testENQ, "testENQ")
	runTest(testDQ, "testDQ")

	return true
}

func main() {
	if len(os.Args) <= 3 {
		log.Fatal("In suffcient Arguments\nShould be:\n$", os.Args[0], " <[ip]:port> <count> <msg_count>")
	}

	COUNT, _ = strconv.Atoi(os.Args[2])
	MCOUNT, _ = strconv.Atoi(os.Args[3])
	server = os.Args[1]
	idlist = make([]string, COUNT)
	qbuffer = make(map[string]list.List)
	verbose = true

	if testAll() {
		log.Println("[All]...................................................Passed")
	} else {
		log.Println("[All]...................................................Failed")
	}

}
