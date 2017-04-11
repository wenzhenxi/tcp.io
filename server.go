package tcpio

import (
	"time"
	"net"
	"fmt"
	"sync"
	"crypto/md5"
	"encoding/base64"
	"bytes"
)

type Server struct {
	*baseHandler
	config            config
	currentConnection int32
	idOrSocket        map[string]*socket
	syncIdOrSocket    sync.RWMutex
}

type config struct {
	PingTimeout   time.Duration
	PingInterval  time.Duration
	MaxConnection int
	NewId         func(string) string
}

func NewServer() (*Server, error) {

	return &Server{
		idOrSocket  : make(map[string]*socket),
		baseHandler : newBaseHandler(),
		config      : config{
			PingTimeout:   60000 * time.Millisecond,
			PingInterval:  25000 * time.Millisecond,
			MaxConnection: 1000,
			NewId:         newId,
		},
	}, nil
}

// SetPingTimeout sets the timeout of a connection ping. When it times out, the server will close the connection with the client. Default is 60s.
func (s *Server) SetPingTimeout(t time.Duration) {
	s.config.PingTimeout = t
}

// SetPingInterval sets the interval of pings. Default is 25s.
func (s *Server) SetPingInterval(t time.Duration) {
	s.config.PingInterval = t
}


// SetMaxConnection sets the maximum number of connections with clients. Default is 1000.
func (s *Server) SetMaxConnection(n int) {
	s.config.MaxConnection = n
}

// GetMaxConnection returns the current max connection
func (s *Server) GetMaxConnection() int {
	return s.config.MaxConnection
}

// Count returns a count of current number of active connections in session
func (s *Server) Count() int {

	return int(s.currentConnection)
}

func newId(RemoteAddr string) string {
	hash := fmt.Sprintf("%s %s", RemoteAddr, time.Now())
	buf := bytes.NewBuffer(nil)
	sum := md5.Sum([]byte(hash))
	encoder := base64.NewEncoder(base64.URLEncoding, buf)
	encoder.Write(sum[:])
	encoder.Close()
	return buf.String()[:20]
}

func (s *Server)Run(ip string, port int) error {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(ip), port, ""})

	if err != nil {
		return err
	}

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}

		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		s.currentConnection++
		id := s.config.NewId(conn.RemoteAddr().String())
		fmt.Println(id)
		socket := newSocket(newServerConn(id, conn), s.baseHandler)
		s.setIdOrSocke(id, socket)
		defer s.Close(id)

		go socket.loop()
	}

	return nil
}

func (s *Server)Close(id string) {
	// 触发disconnection
	socket := s.getIdOrSocke(id)
	s.currentConnection--
	socket.conn.tcp.Close()
	s.delIdOrSocke(id)
}

func (this *Server)getIdOrSocke(id string) *socket {
	this.syncIdOrSocket.Lock()
	rs := this.idOrSocket[id]
	this.syncIdOrSocket.Unlock()
	return rs
}

func (this *Server)setIdOrSocke(id string, So *socket) {
	this.syncIdOrSocket.Lock()
	this.idOrSocket[id] = So
	this.syncIdOrSocket.Unlock()
}

func (this *Server)delIdOrSocke(id string) {
	this.syncIdOrSocket.Lock()
	delete(this.idOrSocket, id)
	this.syncIdOrSocket.Unlock()
}


