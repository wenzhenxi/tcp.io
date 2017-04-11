package tcpio

import (
	"sync"
	"fmt"
)


// Socket is the socket object of socket.io.
type Socket interface {
	// Id returns the session id of socket.
	Id() string

	// On registers the function f to handle an event.
	On(event string, f interface{}) error

	// Emit emits an event with given args.
	Emit(event string, args ...interface{}) error
}

type socket struct {
	*socketHandler
	conn *serverConn
	id   int
	mu   sync.Mutex
}

func newSocket(conn *serverConn, base *baseHandler) *socket {
	ret := &socket{
		conn: conn,
	}
	ret.socketHandler = newSocketHandler(ret, base)
	return ret
}

func (s *socket) Emit(event string, args ...interface{}) error {
	if err := s.socketHandler.Emit(event, args...); err != nil {
		return err
	}
	if event == "disconnect" {
		s.conn.Close()
	}
	return nil
}

func (s *socket) On(event string, f interface{}) error {
	c, err := newCaller(f)
	if err != nil {
		return err
	}
	s.baseHandler.events[event] = c
	return nil
}

func (s *socket) Id() string {

	return s.conn.Id()
}

func (s *socket) loop() {
	defer func() {
		p := &packet{
			Type: _DISCONNECT,
		}
		// 触发断开事件
		s.onPacket(p)
	}()

	p := &packet{
		Type: _CONNECT,
	}
	// 触发事件connection
	s.onPacket(p)

	data := make([]byte, 128)
	for {
		i, err := s.conn.tcp.Read(data)
		fmt.Println("客户端发来数据:", string(data[0:i]))
		if err != nil {
			fmt.Println("读取客户端数据错误:", err.Error())
			break
		}
		//	s.conn.tcp.Write([]byte{'f', 'i', 'n', 'i', 's', 'h'})
	}
}

func (s *socket)trigger() {

}

func (s *socket) sendId() (int, error) {
	s.mu.Lock()

	s.id++
	if s.id < 0 {
		s.id = 0
	}
	s.mu.Unlock()

	return s.id, nil
}