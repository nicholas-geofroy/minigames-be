package lobby

import (
	"minigames-be/src/game"

	"github.com/labstack/echo/v4"
)

type Lobby struct {
	id      string
	in      chan LobbyMsg
	clients map[string]*Client
	log     echo.Logger

	game *game.Game

	Register   chan *Client
	Unregister chan *Client
}

func NewLobby(id string, logger echo.Logger) *Lobby {
	l := &Lobby{
		id:         id,
		clients:    make(map[string]*Client, 0),
		log:        logger,
		in:         make(chan LobbyMsg),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	go l.Run()
	return l
}

func (l *Lobby) Run() {
	l.log.Debug("Run Start")
	l.game = game.NewGame(func(o []game.Object) {
		newState := MessageRender(o)
		l.broadcast(newState)
	})
	l.game.StartGame()

	for {
		select {
		case c := <-l.Register:
			l.log.Debug("Register Client", c)
			l.clients[c.id] = c

		case c := <-l.Unregister:
			l.log.Debug("Unregister Client", c)
			if _, ok := l.clients[c.id]; ok {
				delete(l.clients, c.id)
			} else {
				l.log.Warnf("Client %s was unregistered but never in the map", c.id)
			}

		case m := <-l.in:
			l.log.Info("received message %v", m)
		}
	}
}

func (l *Lobby) broadcast(msg LobbyMsg) {
	for _, c := range l.clients {
		c.out <- msg
	}
}

func (l *Lobby) close() {
	for _, c := range l.clients {
		close(c.out)
	}
}
