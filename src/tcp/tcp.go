package tcp

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

func ReadINT32(conn net.Conn) (int32, error) {

	var rval int32
	tbyte := make([]byte, 4)

	cnt, err := io.ReadAtLeast(conn, tbyte, 4)
	if err != nil {
		log.Fatal("ReadINT32() Could not read 4 bytes")
		return int32(cnt), err
	}
	err = binary.Read(bytes.NewBuffer(tbyte), binary.BigEndian, &rval)

	if err != nil {
		return int32(cnt), err
	}
	return rval, err

}

func ReadBYTE(conn net.Conn) (byte, error) {

	rbyte := make([]byte, 1)

	cnt, err := io.ReadAtLeast(conn, rbyte, 1)

	if err != nil {
		log.Println("ReadBYTE() failed:", err," read ",cnt," bytes rbyte=", rbyte)
	}
	return rbyte[0], err

}

func ReadNBytes(conn net.Conn, N int32) ([]byte, error) {

	rbytes := make([]byte, N)

	cnt, err := io.ReadAtLeast(conn, rbytes, int(N))

	if err != nil {
		log.Fatal("ReadNBytes() Unable to read ", N, " Bytes read only ", cnt, " bytes")

	}

	return rbytes, err
}

func ReadMQID(conn net.Conn) (string, error) {
	rbytes := make([]byte, 32)

	cnt, err := io.ReadAtLeast(conn, rbytes, 32)

	if err != nil {
		log.Fatal("ReadMQID() Unable to read ", 32, " Bytes read only", cnt, " bytes")

	}

	return string(rbytes), err

}

func WriteBYTE(conn net.Conn, b byte) (int, error) {

	Buff := make([]byte, 1)
	Buff[0] = b
	return conn.Write(Buff)
}

func WriteBytes(conn net.Conn, Buff []byte) (int, error) {

	return conn.Write(Buff)

}

func WriteINT32(conn net.Conn, v int32) (int, error) {

	tbuff := new(bytes.Buffer)
	err := binary.Write(tbuff, binary.BigEndian, v)

	if err != nil {
		log.Fatal("WriteINT32() failed:", err)
		return 0, err
	}

	return conn.Write(tbuff.Bytes())
}

func WriteMQID(conn net.Conn, id string) (int, error) {

	return conn.Write([]byte(id))

}
