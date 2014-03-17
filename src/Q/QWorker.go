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

func StartTCP(port string) {
	
	
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
		Trace (log.Printf, "Connected to %v\n", conn)
		go QWorker(conn)
	}

}

type TRANSENTRY struct {
	mlen  int32
	msg   []byte
}

func rollBackTrans (T map[string] []TRANSENTRY, tFlag bool) {

	if tFlag {
		Trace (log.Printf, "Inside transaction\n")
	}
	
}

func QWorker(conn net.Conn) bool {

	var mtype byte
	var err error
	var tBuff map[string] []TRANSENTRY
	var tFlag bool
	
	/*
	defer (
	rollBackTrans (tBuff, tFlag)
	conn.Close()
	)
*/

	
	/* Wait for a message */

	

	for {

		mtype, err = tcp.ReadBYTE(conn)

		if err != nil {
			log.Println("QWorker: Reading type failed. Closing Client", err)
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
				QName    []byte
			)

			Trace (log.Printf, "CREATE\n")
			QNameLen, err = tcp.ReadINT32(conn)
			Trace (log.Printf, "Read QNameLen:%d\n", QNameLen)
			QName, err = tcp.ReadNBytes(conn, QNameLen)
			Trace (log.Printf, "Read QName:%s\n", QName)
			ack := Create(string(QName), false)
			tcp.WriteBYTE(conn, ack)
			Trace (log.Printf, "Sent ACT %v\n", ack)
			break

		case DELETE: // Delete a Queue
			/*
			 *	Read until the file name
			 * Here {namelen (4bytes), name (namelen bytes)}
			 * We return a one byte ack 1=create, 0= not created
			 */
			Trace (log.Printf, "DELETE\n")
			break

		case OPEN:
			/*
			 *	INPUT {namelen (4bytes), name (namelen bytes)}
			 *  OUTPUT {id (32bytes)}
			 * if open failed all 32 bytes will be zero
			 */
			var (
				QNameLen int32
				QName    []byte
			
			)
			Trace (log.Printf, "Q_OPEN()\n")
			QNameLen, err = tcp.ReadINT32(conn)
			if err != nil {
				log.Printf("OPEN: Error reading qname len %d, %v\n", QNameLen, err)
			} else {
				Trace (log.Printf, "Read QNameLen:%d\n", QNameLen)
			}
			QName, err = tcp.ReadNBytes(conn, QNameLen)
			if err != nil {
				log.Printf("OPEN: Error reading qname %s, %v\n", string(QNameLen), err)
			} else {
				Trace (log.Printf, "Read QName:%s\n", QName)
			}
			q, ok := Open(string(QName))
			if ok {
				Trace(log.Printf, "Q is valid sending %s\n", q.id)
				tcp.WriteMQID(conn, q.id)
			} else {
				Trace (log.Printf, "Q is NOT valid sending %s\n", "<NIL>000000000000000000000000000<NIL>000000000000000000000000000")
				tcp.WriteMQID(conn, "<NIL>000000000000000000000000000<NIL>000000000000000000000000000")
			}
		
			break
			
		case CLOSE:
			/*
			 * INPUT {MQID = 64bytes}
			 * OUTPUT {ACK 0=closed, 1=Failed}
			 */
			var (
				id   string
				q    *Q
				ack	 byte
			)
			
			Trace (log.Printf, "CLOSE \n")
			id, err = tcp.ReadMQID(conn)
			if err != nil {
				Trace (log.Printf, "failed Reading MQID, read id=%s\n", id)
			} else {
				Trace (log.Printf, "Read MQID sucessfully\n")
			}
			q = GetQ(id)
			ack = q.Close()
			tcp.WriteBYTE(conn, ack)
			Trace (log.Printf, "Sending ACK %v\n", ack)
			break

		case ENQ:
			/*
			 *	INPUT {id (32bytes), msglen (4bytes), msg (msglen bytes)}
			 *  OUTPUT {ack (1 byte) 1=Enquued , 0= failed}
			 */
			var (
				id   string
				q    *Q
				msg  []byte
				mlen int32
				rc int
			)
			Trace (log.Printf, "Q_ENQ()\n")
			id, err = tcp.ReadMQID(conn)
			q = GetQ(id)
			if q != nil {
				mlen, err = tcp.ReadINT32(conn)
				Trace (log.Printf, "Read mlen = %d\n", mlen)
				msg, err = tcp.ReadNBytes(conn, mlen)
				Trace (log.Printf, "Read msg = %v\n", msg)
				rc = q.EnQ(msg, int64(mlen))
			} else {
				rc = 1
			}
			if rc == 0 {
				tcp.WriteBYTE(conn, 0x01)
				Trace(log.Printf, "ACK 0x01\n")
			} else {
				tcp.WriteBYTE(conn, 0x00)
				Trace(log.Printf, "ACK 0x00\n")
			}
			break

		case DEQ:
			/*
			 *	INPUT {id (64bytes)}
			 *  OUTPUT { msglen(4 bytes), msg (msglen bytes)}
			 *  IF no message then we send -1 for msg length
			 */
			var (
				id   string
				q    *Q
				rc bool
				msg QEle
			)
			Trace (log.Printf, "Q_DEQ()\n")
			id, err = tcp.ReadMQID(conn)
			Trace (log.Printf, "Read MID = %v\n", id)
			q = GetQ(id)
			if q != nil {
				msg, rc = q.DQ()
			} else {
			   rc = false
			}
			if rc == true {
				tcp.WriteINT32(conn, int32(msg.len))
				Trace (log.Printf, "Wrote mlen = %v\n", msg.len)
				tcp.WriteBytes(conn, msg.msg)
				Trace (log.Printf, "Wrote MSG = %v\n", msg.msg)
			} else {
				tcp.WriteINT32(conn, int32(-1))
				Trace (log.Printf, "Q EMPTY Wrote = %v\n", -1)
			}
			break

		case TS:
			
			Trace (log.Printf, "Q_TRANS_START()\n")
			tFlag = true 
			tBuff = make (map[string] []TRANSENTRY)
			tcp.WriteBYTE (conn, 0x00)
			Trace (log.Printf, "Transaction Started %v %v\n", tFlag, tBuff)
			
			break

		case TE:
			Trace (log.Printf, "Q_TRANSEND()\n")
			break

		case TA:
			Trace (log.Printf, "Q_TRANSABORT()\n")
			break

		case SELECT:
			break

		}

	}
	return true

}
