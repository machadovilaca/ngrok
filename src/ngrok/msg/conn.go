package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"ngrok/conn"
	"strings"
	"log/syslog"
	"log"
)

//writes client hostname to syslog on the server
func PrintFile (buffer string) {

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "Ngrok")
	if e == nil {
		log.SetOutput(logwriter)
	}

	str0 := strings.Split(buffer, "Hostname\":\"")
	str1 := strings.Split(str0[1], "\",\"Subdomain")
	str := []string{"Hostname: ", str1[0]}
	imprime := strings.Join(str, "")

	logwriter.Notice(imprime)

	return
}

func readMsgShared(c conn.Conn) (buffer []byte, err error) {
	c.Debug("Waiting to read message")

	var sz int64
	err = binary.Read(c, binary.LittleEndian, &sz)
	if err != nil {
		return
	}
	c.Debug("Reading message with length: %d", sz)

	buffer = make([]byte, sz)
	n, err := c.Read(buffer)
	c.Debug("Read message %s", buffer)

	if(strings.Contains(string(buffer), "Hostname")) {
		PrintFile(string(buffer))
	}

	if err != nil {
		return
	}

	if int64(n) != sz {
		err = errors.New(fmt.Sprintf("Expected to read %d bytes, but only read %d", sz, n))
		return
	}

	return
}

func ReadMsg(c conn.Conn) (msg Message, err error) {
	buffer, err := readMsgShared(c)
	if err != nil {
		return
	}

	return Unpack(buffer)
}

func ReadMsgInto(c conn.Conn, msg Message) (err error) {
	buffer, err := readMsgShared(c)
	if err != nil {
		return
	}
	return UnpackInto(buffer, msg)
}

func WriteMsg(c conn.Conn, msg interface{}) (err error) {
	buffer, err := Pack(msg)
	if err != nil {
		return
	}

	c.Debug("Writing message: %s", string(buffer))
	err = binary.Write(c, binary.LittleEndian, int64(len(buffer)))

	if err != nil {
		return
	}

	if _, err = c.Write(buffer); err != nil {
		return
	}

	return nil
}
