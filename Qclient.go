package main

import (
//	"bytes"
//	"encoding/binary"
	"fmt"
	"net"
	"os"
	"container/list"
	"strconv"
	"tcp"
	"log"
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
  COUNT int
  server string
  idlist []string
  qbuffer map[string]list.List
  verbose bool
  sconn	net.Conn
)

func vprintf (format string, v ...interface{}) {

	if verbose {
		log.Printf(format, v...)
	}
}

func testCREATE() bool {
	
	var tmp int
	vprintf ("Creating Q\n")
	for i:=0; i<COUNT; i++ {
		cnt , err := tcp.WriteBYTE(sconn, CREATE)
		if err != nil {
			vprintf ("Writebyte error wrote %dbytes error %v\n", cnt, err)
		}
		vprintf ("Writebyte wrote %dbytes error %v\n", cnt, err)
		//fmt.Scanf("%d", &tmp)
		qname := []byte("TEMPQ"+strconv.Itoa(i))
		cnt , err = tcp.WriteINT32 (sconn, int32(len(qname)))
		if err != nil {
			vprintf ("Writebyte error wrote %dbytes error %v\n", cnt, err)
		}
		vprintf ("qnamelen wrote value=%d %dbytes error %v\n", int32(len(qname)), cnt, err)
		//fmt.Scanf("%d", &tmp)
		cnt , err = tcp.WriteBytes (sconn, qname)
		if err != nil {
			vprintf ("Writebyte error wrote %dbytes error %v\n", cnt, err)
		}
		vprintf ("Qname wrote %dbytes error %v\n", cnt, err)
		//fmt.Scanf("%d", &tmp)
		ack, ok := tcp.ReadBYTE(sconn)
		if ok != nil {
			vprintf("Error Creating Q %v error %v\n", qname, ok)
			return false
		}
		//fmt.Scanf("%d", &tmp)
		if ack == 0x01 {
			vprintf ("Q Created\n")
		} else {
			vprintf ("Q Creation failed ack=%v\n", ack)
		}
	}
	fmt.Scanf("%d", &tmp)
	return true
}

func testOPEN () bool {

	return true
}

func testENQ () bool {
	return true
}

func testDQ () bool {
	return true
}

func runTest (test func()bool, val string) {
	if test() {
		log.Printf("[%s]...\t...\t...\t...[OK]", val)
	} else {
		log.Printf("[%s]...\t...\t...\t...[FAIL]", val)
	}
}

func testAll() bool {
	
	var err error
/* Connect to the server */
	sconn, err = net.Dial("tcp", server)
	
	if err != nil {
		vprintf ("Error Connecting %v\n", err)
		return false
	} else {
		vprintf ("Connected to Server %v\n", sconn)
	}
		
/* Run the tests */
	runTest(testCREATE, "testCREATE")
	runTest(testOPEN, "testOPEN")
	runTest(testENQ, "testENQ")
	runTest(testDQ, "testDQ")

	return true
}



func main() {
	if len(os.Args) < 2 {
		log.Fatal("In suffcient Arguments\nShould be:\n$", os.Args[0], "<[ip]:port> <count>")
	}
	
	
	COUNT , _ = strconv.Atoi(os.Args[2])
	server = os.Args[1]
	idlist = make([]string, COUNT)
	qbuffer = make(map[string]list.List)
	verbose = true
	
	
	if testAll() {
		log.Println ("[All]...................................................Passed")
	} else {
		log.Println ("[All]...................................................Failed")
	}
		
	
}
