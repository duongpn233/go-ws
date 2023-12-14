package socket

import "sync"

type EachFunc func(Socket)

type RoomManager interface {
	Join(room string, s Socket)
	Leave(room string, s Socket) bool
	LeaveAll(s Socket) bool
	Clear(room string) bool
	Send(roomName string, clientId string, event string, data ContentType) bool
	Broadcast(roomName, event string, data ContentType) bool
	ForEach(room string, f EachFunc)
	Len(room string) int
}

type roomManager struct {
	rooms map[string]*room
	lock  sync.RWMutex
}

func newRoomManager() *roomManager {
	return &roomManager{
		rooms: make(map[string]*room),
	}
}

type room struct {
	name          string
	clients       map[string]Socket
	quitChan      chan struct{}
	broadcastChan chan Message
	stopped       bool
	lock          sync.RWMutex
}

func (room *room) join(socket Socket) {
	room.lock.Lock()
	defer room.lock.Unlock()
	room.clients[socket.ID()] = socket
}

func (room *room) leave(socket Socket) {
	room.lock.Lock()
	defer room.lock.Unlock()
	delete(room.clients, socket.ID())
}

func (room *room) run() {
	for {
		select {
		case <-room.quitChan:
			return
		case msg := <-room.broadcastChan:
			for _, client := range room.clients {
				go func(socket Socket) {
					socket.Emit(msg.Event, msg.Content)
				}(client)
			}
		}
	}
}

func (rm *roomManager) Join(roomName string, s Socket) {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	r, ok := rm.rooms[roomName]
	if !ok {
		r = &room{
			name:          roomName,
			clients:       make(map[string]Socket),
			broadcastChan: make(chan Message),
			quitChan:      make(chan struct{}),
			stopped:       false,
		}
		rm.rooms[roomName] = r
		go r.run()
	}

	r.join(s)
}

func (rm *roomManager) Leave(roomName string, socket Socket) bool {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	room, ok := rm.rooms[roomName]
	if !ok {
		return false
	}

	room.leave(socket)
	if len(room.clients) == 0 {
		room.stopped = true
		delete(rm.rooms, roomName)
		close(room.quitChan)
		close(room.broadcastChan)
	}
	return true
}

func (rm *roomManager) LeaveAll(socket Socket) bool {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	for _, room := range rm.rooms {
		room.leave(socket)
		if len(room.clients) == 0 {
			room.stopped = true
			delete(rm.rooms, room.name)
			close(room.quitChan)
			close(room.broadcastChan)
		}
	}

	return true
}

func (rm *roomManager) Clear(roomName string) bool {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	room, ok := rm.rooms[roomName]
	if !ok {
		return false
	}

	room.stopped = true
	delete(rm.rooms, roomName)
	close(room.quitChan)
	close(room.broadcastChan)

	return true
}

func (rm *roomManager) Send(roomName string, clientId string, event string, data ContentType) bool {
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	room, ok := rm.rooms[roomName]
	if !ok {
		return false
	}

	client, ok := room.clients[clientId]
	if !ok {
		return false

	}

	client.Emit(event, data)
	return true
}

func (rm *roomManager) Broadcast(roomName, event string, data ContentType) bool {
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	room, ok := rm.rooms[roomName]
	if !ok || room.stopped {
		return false
	}

	room.broadcastChan <- Message{
		Event:   event,
		Content: data,
	}
	return true
}

func (rm *roomManager) ForEach(room string, f EachFunc) {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	r, ok := rm.rooms[room]
	if !ok {
		return
	}

	for _, s := range r.clients {
		f(s)
	}
}

func (rm *roomManager) Len(room string) int {
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	return len(rm.rooms[room].clients)
}
