package lobby

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	id     string
	lobby  *Lobby
	conn   *websocket.Conn
	logger echo.Logger
	out    chan LobbyMsg
}

func NewClient(clientId string, lobby *Lobby, conn *websocket.Conn, logger echo.Logger) *Client {
	c := &Client{
		id:     clientId,
		lobby:  lobby,
		conn:   conn,
		logger: logger,
		out:    make(chan LobbyMsg),
	}
	go c.Run()
	return c
}

func (c *Client) Run() {
	c.logger.Debug("Start Run for client ", c.id)
	c.lobby.Register <- c
	c.logger.Debug("Start register sent to lobby ", c.lobby.id)
	go c.readMessages()
	go c.writeMessages()
}

type LobbyInMsg struct {
	MsgType int
	Data    map[string]interface{}
}

func (c *Client) readMessages() {
	defer c.Close()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		msg := LobbyInMsg{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorf("error: %v", err)
			}
			break
		}

		c.logger.Debugf("Received message:\n type: %d\n data:%+v\n", msg.MsgType, msg.Data)
	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case msg, ok := <-c.out:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			outMsg := &OutMsg{MsgType: msg.Type(), Data: msg}
			err := c.conn.WriteJSON(outMsg)
			if err != nil {
				c.logger.Warnf("Error in client %s", err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	c.lobby.Unregister <- c
	c.conn.Close()
}
