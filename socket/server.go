package socket

import (
	"log"
	"net/http"

	"github.com/gobwas/ws"
)

type Server struct {
	rm      *roomManager
	handler *handler
}

func NewEngineSocket() *Server {
	return &Server{
		rm:      newRoomManager(),
		handler: newHandler(),
	}
}

func (e *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println(err)
		return
	}

	socket := newSocket(conn)

	bc := NewBroadcastSocket(socket, e.rm)

	go socket.write(e, bc)
	go socket.error(e, bc)
	go socket.read(e, bc)

	e.handler.connect(bc)
}

// Handle data received in the callback as a json string
func (e *Server) OnEvent(event string, f func(Socket, ContentType)) {
	e.handler.addEventHandler(event, f)
}

func (e *Server) OnError(f func(Socket, error)) {
	e.handler.addErrorHandler(f)
}

func (e *Server) OnConnect(f func(Socket) error) {
	e.handler.addConnectHandler(f)
}

func (e *Server) OnDisconnect(f func(Socket, string)) {
	e.handler.addDisconnectHandler(f)
}

func (e *Server) JoinRoom(roomName string, socket Socket) {
	e.rm.Join(roomName, socket)
}

func (e *Server) LeaveRoom(roomName string, socket Socket) bool {
	return e.rm.Leave(roomName, socket)
}

func (e *Server) LeaveAllRooms(socket Socket) bool {
	return e.rm.LeaveAll(socket)
}

func (e *Server) ClearRoom(roomName string) bool {
	return e.rm.Clear(roomName)
}

func (e *Server) SendMessageToClient(roomName string, clientId string, event string, data ContentType) bool {
	return e.rm.Send(roomName, clientId, event, data)
}

func (e *Server) BroadcastToRoom(roomName string, event string, data ContentType) bool {
	return e.rm.Broadcast(roomName, event, data)
}

func (s *Server) ForEach(room string, f EachFunc) {
	s.rm.ForEach(room, f)
}
