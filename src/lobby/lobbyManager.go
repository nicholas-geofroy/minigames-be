package lobby

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

type LobbyManager struct {
	ctx     echo.Context
	lobbies map[string]*Lobby

	register   chan *Lobby
	unregister chan *Lobby
}

func (m *LobbyManager) run() {
	for {
		select {
		case lobby := <-m.register:
			m.lobbies[lobby.id] = lobby

		case lobby := <-m.unregister:
			if _, ok := m.lobbies[lobby.id]; ok {
				delete(m.lobbies, lobby.id)
			} else {
				m.ctx.Error(
					errors.New(
						fmt.Sprintf(
							"Attempted to unregister lobby <%s> which was never registered",
							lobby.id)))
			}
		}
	}
}
