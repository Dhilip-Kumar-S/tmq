package Q

import (
	"log"
	"net"
)

const (
	CREATE = iota
	DELETE
	OPEN
	CLOSE
	ENQ
	DEQ
	TS
	TE
	TA
)

/*
 * This is a TCP server and it runs for ever waiting for connection
 * Spwans a go routine per TCP connection.
 * Each Go routine will be capable of handling requests.
 */

func Server(port string) bool {

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

type Protocol struct {
	op  int // The operation you want to perform on the server
	len int // Length of the message to read on this socekt
}

func QWorker(conn net.Conn) bool {

	var msg Protocol
	/* Wait for a message */

	for {

		switch msg.op {
		case CREATE:
			/*
			 * Read until the file name and persistance value
			 * Here {namelen(4bytes), name(namelen bytes), store (1 byte)}
			 * We return a one byte ack 1=created, 0=not created
			 */
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

		}

	}
	return true

}
