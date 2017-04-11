package tcpio

import (
	"sync"
	"reflect"
	"fmt"
	"github.com/spf13/cast"
	"encoding/json"
)

type baseHandler struct {
	events map[string]*caller
}

func newBaseHandler() *baseHandler {
	return &baseHandler{
		events:    make(map[string]*caller),
	}
}

func (h *baseHandler) On(event string, f interface{}) error {
	c, err := newCaller(f)
	if err != nil {
		return err
	}
	h.events[event] = c
	return nil
}

type socketHandler struct {
	*baseHandler
	acksmu sync.Mutex
	acks   map[int]*caller
	socket *socket
	rooms  map[string]struct{}
}

func newSocketHandler(s *socket, base *baseHandler) *socketHandler {
	events := make(map[string]*caller)
	for k, v := range base.events {
		events[k] = v
	}
	return &socketHandler{
		baseHandler: &baseHandler{
			events:    events,
		},
		acks:   make(map[int]*caller),
		socket: s,
		rooms:  make(map[string]struct{}),
	}
}

func (h *socketHandler) Emit(event string, args ...interface{}) error {

	var c *caller
	var id int
	var err error
	if l := len(args); l > 0 {
		fv := reflect.ValueOf(args[l - 1])
		if fv.Kind() == reflect.Func {
			var err error
			c, err = newCaller(args[l - 1])
			if err != nil {
				return err
			}
			args = args[:l - 1]
		}
	}
	args = append([]interface{}{event}, args...)
	if c != nil {
		id, err = h.socket.sendId()
		if err != nil {
			return err
		}
		h.acksmu.Lock()
		h.acks[id] = c
		h.acksmu.Unlock()
	}

	fmt.Println(id)
	rs, err := json.Marshal(cast.ToStringSlice(args))
	fmt.Println(string(rs))
	fmt.Println(h.acks)

	s := "42" + cast.ToString(id) + string(rs)

	err = h.socket.conn.NextWriter(s)

	return err
}

func (h *socketHandler) onPacket(packet *packet) ([]interface{}, error) {

	c, ok := h.events[packet.Type]
	if !ok {
		// If the message is not recognized by the server, the decoder.currentCloser
		// needs to be closed otherwise the server will be stuck until the e
		return nil, nil
	}

	args := c.GetArgs()
	olen := len(args)

	if olen > 0 {
		args[0] = &packet.Data
	}

	retV := c.Call(h.socket, args)
	if len(retV) == 0 {
		return nil, nil
	}

	var err error

	if last, ok := retV[len(retV) - 1].Interface().(error); ok {
		err = last
		retV = retV[0 : len(retV) - 1]
	}
	ret := make([]interface{}, len(retV))
	for i, v := range retV {
		ret[i] = v.Interface()
	}

	return ret, err
}


// 触发发送消息端进行的回调
func (h *socketHandler) onAck(id int, packet *packet) error {
	h.acksmu.Lock()
	_, ok := h.acks[id]
	if !ok {
		h.acksmu.Unlock()
		return nil
	}
	delete(h.acks, id)
	h.acksmu.Unlock()
	//
	//args := c.GetArgs()
	//packet.Data = &args
	//c.Call(h.socket, args)
	return nil
}

