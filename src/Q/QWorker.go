package Q

import (
	"log"
	"net"
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

/*
 * This is a TCP server and it runs for ever waiting for connection
 * Spwans a go routine per TCP connection.
 * Each Go routine will be capable of handling requests.
 */

func Start(port string) {

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
		panic("Server Failed:")
	}

	for {

		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("Accept() Error:")
		}
		go QWorker(conn)
	}

}

type TRANSENTRY struct {
	mtype byte
	mlen  int32
	msg   []byte
}

func QWorker(conn net.Conn) bool {

	var mtype byte
	var err error

	differ conn.Close()
	
	var inTransaction bool

	inTransaction = false
	/* Wait for a message */

	transBuffer := make(map[string]TRANSENTRY)

	for {

		mtype, err = ReadBYTE(conn)
		
		if err != nil {
			log.Fatal ("QWorker: Reading type failed. Closing Client", err)
			break
		}
		
		switch mtype {
		case CREATE:
			/*
			 * Read until the file name and persistance value
			 * Here {namelen(4bytes), name(namelen bytes), store (1 byte)}
			 * We return a one byte ack 1=created, 0=not created
			 */
			var (
			QNameLen int32
			QName string
			)
			
			QNameLen, err = tcp.ReadINT32(conn) 
			QName, err = tcp.ReadNBytes(QNameLen)
			ack := Create(string(QName), false)
			tcp.WriteBYTE(conn, ack)
			break
		case DELETE: // Delete a Queue
			/*
			 *	Read until the file name
			 * Here {namelen (4bytes), name (namelen bytes)}
			 * We return a one byte ack 1=create, 0= not created
			 */
			break
		case OPEN:
			/*
			 *	INPUT {namelen (4bytes), name (namelen bytes)}
			 *  OUTPUT {id (32bytes)}
			 * if open failed all 32 bytes will be zero
			 */
			
			break
		case ENQ:
			/*
			 *	INPUT {id (32bytes), msglen (4bytes), msg (msglen bytes)}
			 *  OUTPUT {ack (1 byte) 1=Enquued , 0= failed}
			 */
			break
		case DEQ:
			/*
			 *	INPUT {id (32bytes)}
			 *  OUTPUT { msglen(4 bytes), msg (msglen bytes)}
			 */
			break

		case TS:
			break

		case TE:
			break

		case TA:
			break

		case SELECT:
			break

		}

	}
	return true

}
