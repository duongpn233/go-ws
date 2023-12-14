package socket

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Client struct {
	handler *handler
	*socket
	url     string
	context interface{}
}

func NewClient(addr string) *Client {
	return &Client{
		handler: newHandler(),
		url:     addr,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	conn, _, _, err := ws.DefaultDialer.Dial(ctx, c.url)
	if err != nil {
		return err
	}

	c.socket = newSocket(conn)

	go c.error()
	go c.write()
	go c.read()

	return nil
}

func (c *Client) OnConnect(f func(Socket) error) {
	c.handler.addConnectHandler(f)
}

func (c *Client) OnDisconnect(f func(Socket, string)) {
	c.handler.addDisconnectHandler(f)
}

func (c *Client) OnError(f func(Socket, error)) {
	c.handler.addErrorHandler(f)
}

func (c *Client) OnEvent(event string, f func(Socket, ContentType)) {
	c.handler.addEventHandler(event, f)
}

func (c *Client) write() {
	defer func() {
		if err := c.Close(); err != nil {
			log.Println("close connect:", err)
		}
	}()

	for {
		select {
		case <-c.socket.quitChan:
			return
		case msg, ok := <-c.socket.sendChan:
			if !ok {
				continue
			}

			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				c.socket.errChan <- err
				continue
			}

			err = wsutil.WriteClientMessage(c.socket.conn, ws.OpText, jsonMsg)
			if err != nil {
				c.socket.errChan <- err
				continue
			}
		}
	}
}

func (c *Client) read() {
	defer func() {
		if err := c.Close(); err != nil {
			log.Println("close connect:", err)
		}
	}()

	for {
		select {
		case <-c.socket.quitChan:
			return
		default:
			msg, _, err := wsutil.ReadServerData(c.socket.conn)
			if err != nil {
				c.socket.errChan <- err
				if err == io.EOF {
					c.socket.isClose = true
					close(c.socket.quitChan)
					close(c.socket.sendChan)
					close(c.socket.errChan)
					msg := "connection close ID: " + c.socket.ID()
					if c.handler.disconnect != nil {
						c.handler.disconnect(c, msg)
					}
					return
				}
				continue
			}

			var body Message
			if err := json.Unmarshal(msg, &body); err != nil {
				c.socket.errChan <- err
				continue
			}

			if c.handler.events[body.Event] != nil {
				go c.handler.events[body.Event](c, body.Content)
			}
		}
	}
}

func (c *Client) error() {
	defer func() {
		if err := c.Close(); err != nil {
			log.Println("close connect:", err)
		}
	}()

	for {
		select {
		case <-c.socket.quitChan:
			return
		case err := <-c.socket.errChan:
			if c.handler.err != nil {
				go c.handler.err(c, err)
			}
		}
	}
}

func (c *Client) Close() error {
	return c.socket.Close()
}

func (c *Client) SetContext(ctx interface{}) {
	c.context = ctx
}

func (c *Client) Context() interface{} {
	return c.context
}

func (c *Client) Emit(eventName string, data ContentType) {
	c.socket.send(eventName, data)
}

// Func doesn't work
func (c *Client) Join(room string) {

}

// Func doesn't work
func (c *Client) Leave(room string) {

}

// Func doesn't work
func (c *Client) LeaveAll() {

}
