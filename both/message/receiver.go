package message

import (
	"net"
	"strconv"
	"strings"
)

func ReadFromSocket(conn net.Conn) (string, error) {
	msg, err := readFromSocketWithLength(conn, BufferLength)
	if err != nil {
		return "", err
	}

	bytesLen, err := strconv.Atoi(strings.TrimSpace(msg))
	if err != nil {
		return "", err
	}

	if bytesLen <= 0 {
		return "", nil
	}

	return readFromSocketWithLength(conn, bytesLen)
}

func readFromSocketWithLength(conn net.Conn, length int) (string, error) {
	buffer := make([]byte, length)
	length, err := conn.Read(buffer)

	if err != nil {
		return "", err
	}

	if length <= 0 {
		return "", nil
	}

	return string(buffer[:length]), nil
}
