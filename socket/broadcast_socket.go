package socket

type Broadcast interface {
	Context() interface{}
	SetContext(ctx interface{})
	Emit(eventName string, data ContentType)
	Join(room string)
	Leave(room string)
	LeaveAll()
}

type broadcastSocket struct {
	*socket
	rm      RoomManager
	context interface{}
}

func NewBroadcastSocket(sk *socket, rm RoomManager) *broadcastSocket {
	return &broadcastSocket{
		socket: sk,
		rm:     rm,
	}
}

func (bc *broadcastSocket) SetContext(ctx interface{}) {
	bc.context = ctx
}

func (bc *broadcastSocket) Context() interface{} {
	return bc.context
}

func (bc *broadcastSocket) Emit(eventName string, data ContentType) {
	bc.socket.send(eventName, data)
}

func (bc *broadcastSocket) Join(room string) {
	bc.rm.Join(room, bc)
}

func (bc *broadcastSocket) Leave(room string) {
	bc.rm.Leave(room, bc)
}

func (bc *broadcastSocket) LeaveAll() {
	bc.rm.LeaveAll(bc)
}
