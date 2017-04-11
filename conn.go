package tcpio

import (
	"net"
	"github.com/wenzhenxi/tcp.io/errors"
	"github.com/spf13/cast"
)


// Conn is the connection object of engine.io.
type Conn interface {
	// Id returns the session id of connection.
	Id() string

	// Close closes the connection.
	Close() error

	// NextReader returns the next message type, reader. If no message received, it will block.
	NextReader()

	// NextWriter returns the next message writer with given message type.
	NextWriter()
}

type serverConn struct {
	id  string
	tcp net.Conn
}

func newServerConn(id string, tcp net.Conn) *serverConn {
	conn := &serverConn{
		id : id,
		tcp : tcp,
	}

	return conn
}

func (c *serverConn) Id() string {
	return c.id
}

func (c *serverConn) Close() error {
	return c.tcp.Close()
}

func (c *serverConn) NextReader() {

}

func (c *serverConn) NextWriter(msg string) error {
	requestLen := cast.ToString(len(msg))

	if len(requestLen) > 5 {
		return errors.ErrRequestLong
	} else if len(requestLen) < 5 {
		for i := len(requestLen); i < 5; i++ {
			requestLen = "0" + requestLen
		}
	}

	msg = requestLen + msg

	_, err := c.tcp.Write([]byte(msg))

	return err
}
