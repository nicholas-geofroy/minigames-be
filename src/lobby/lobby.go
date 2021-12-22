package lobby

import (
	"minigames-be/src/game"

	"github.com/labstack/echo/v4"
)

type Lobby struct {
	id      string
	in      chan ClientMsg
	clients map[string]*Client
	log     echo.Logger

	game *game.Game

	Register   chan *Client
	Unregister chan *Client

	receiveFunc func(ClientMsg)
}

func NewLobby(id string, logger echo.Logger) *Lobby {
	l := &Lobby{
		id:         id,
		clients:    make(map[string]*Client, 0),
		log:        logger,
		in:         make(chan ClientMsg),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	l.receiveFunc = l.lobbyReceive
	go l.Run()
	return l
}

func (l *Lobby) Run() {
	l.log.Debug("Run Start")
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
			l.log.Info("received message", m)
			l.receiveFunc(m)
		}
	}
}

func (l *Lobby) lobbyReceive(m ClientMsg) {
	switch m.msg.(type) {
	case *StartMsg:
		pIds := make([]string, 0, len(l.clients))
		for pId := range l.clients {
			pIds = append(pIds, pId)
		}
		l.log.Infof("start game with players %+v", pIds)
		l.game = game.NewGame(func(o []game.Object) {
			newState := MessageRender(o)
			l.broadcast(newState)
		}, pIds)
		l.game.StartGame()
		l.receiveFunc = l.gameReceive
	default:
		l.clients[m.clientId].out <- &ErrorMsg{
			errorType: 0,
			msg:       "Invalid Move",
		}
	}
}

func (l *Lobby) gameReceive(m ClientMsg) {
	switch msg := m.msg.(type) {
	case *MoveMsg:
		l.game.MakeMove(m.clientId, msg.ToVec())
	default:
		l.clients[m.clientId].out <- &ErrorMsg{
			errorType: 0,
			msg:       "Invalid Move",
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
