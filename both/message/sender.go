package message

import (
	"fmt"
	"math"
	"net"
)

func Send(conn net.Conn, msg []byte) bool {
	bytesLen := len(msg)
	if float64(bytesLen) > math.Pow(10, float64(BufferLength-1))-1 {
		return false
	}

	format := fmt.Sprintf("%%%dv", BufferLength)
	bufferMsg := []byte(fmt.Sprintf(format, bytesLen))

	_, err := conn.Write(bufferMsg)
	if err != nil {
		return false
	}

	_, err = conn.Write(msg)
	if err != nil {
		return false
	}

	return true
}
