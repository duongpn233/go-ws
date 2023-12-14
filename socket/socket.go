package socket

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
)

type Socket interface {
	io.Closer
	Broadcast
	ID() string
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

type socket struct {
	id       string
	conn     net.Conn
	sendChan chan Message
	quitChan chan struct{}
	errChan  chan error
	isClose  bool
}

func newSocket(conn net.Conn) *socket {
	return &socket{
		id:       uuid.NewString(),
		conn:     conn,
		sendChan: make(chan Message),
		quitChan: make(chan struct{}),
		errChan:  make(chan error),
		isClose:  false,
	}
}

func (s *socket) ID() string {
	return s.id
}

func (s *socket) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *socket) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *socket) Close() error {
	return s.conn.Close()
}

func (s *socket) read(e *Server, bc Socket) {
	defer func() {
		if err := s.Close(); err != nil {
			log.Println("close connect:", err)
		}
		e.LeaveAllRooms(bc)
	}()

	for {
		select {
		case <-s.quitChan:
			return
		default:
			msg, _, err := wsutil.ReadClientData(s.conn)
			if err != nil {
				s.errChan <- err
				if err == io.EOF {
					s.isClose = true
					close(s.quitChan)
					close(s.sendChan)
					close(s.errChan)
					msg := "connection close ID: " + s.id
					if e.handler.disconnect != nil {
						e.handler.disconnect(bc, msg)
					}
					return
				}
				continue
			}

			var body Message
			if err := json.Unmarshal(msg, &body); err != nil {
				s.errChan <- err
				continue
			}

			if e.handler.events[body.Event] != nil {
				go e.handler.events[body.Event](bc, body.Content)
			}
		}
	}
}

func (socket *socket) write(e *Server, bc Socket) {
	defer func() {
		if err := socket.Close(); err != nil {
			log.Println("close connect:", err)
		}
		e.LeaveAllRooms(bc)
	}()

	for {
		select {
		case <-socket.quitChan:
			return
		case msg, ok := <-socket.sendChan:
			if !ok {
				continue
			}

			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				socket.errChan <- err
				continue
			}

			err = wsutil.WriteServerMessage(socket.conn, ws.OpText, jsonMsg)
			if err != nil {
				socket.errChan <- err
				continue
			}
		}
	}
}

func (socket *socket) error(e *Server, bc Socket) {
	defer func() {
		if err := socket.Close(); err != nil {
			log.Println("close connect:", err)
		}
		e.LeaveAllRooms(bc)
	}()

	for {
		select {
		case <-socket.quitChan:
			return
		case err := <-socket.errChan:
			if e.handler.err != nil {
				go e.handler.err(bc, err)
			}
		}
	}
}

func (socket *socket) send(eventName string, data ContentType) {
	if !socket.isClose {
		socket.sendChan <- Message{
			Event:   eventName,
			Content: data,
		}
	}
}

func (socket *socket) ping(e *Server, bc Socket) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := wsutil.WriteServerMessage(socket.conn, ws.OpPing, nil); err != nil {
				close(socket.quitChan)
				msg := "socket ID: " + socket.id + "close " + err.Error()
				if e.handler.disconnect != nil {
					e.handler.disconnect(bc, msg)
				}
				return
			}
		}
	}
}
