package main

import (
	"minigames-be/src/lobby"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return checkOrigin(r.Header.Get("Origin"))
		},
	}
)

var (
	lobbies   = make(map[string]*lobby.Lobby)
	lobbyLock sync.Mutex
)

func checkOrigin(origin string) bool {
	switch origin {
	case "http://localhost:3000", "localhost:3000":
		return true
	}
	return false
}

func registerClient(c echo.Context, ws *websocket.Conn, clientId string, lbyId string) {
	c.Logger().Info("register client")
	lobbyLock.Lock()

	lby, ok := lobbies[lbyId]
	if !ok {
		c.Logger().Infof("Creating new lobby %s", lbyId)
		lby = lobby.NewLobby(lbyId, c.Logger())
		lobbies[lbyId] = lby
	}

	lobbyLock.Unlock()

	_ = lobby.NewClient(clientId, lby, ws, c.Logger())
}

func handleWs(c echo.Context) error {
	c.Logger().Debug("handleWs")
	lobbyId := c.Param("lobbyId")
	if lobbyId == "" {
		return c.String(http.StatusBadRequest, "LobbyId Cannot be an empty string")
	}
	clientId := c.Request().URL.Query().Get("clientId")
	if clientId == "" {
		return c.String(http.StatusBadRequest, "clientId must be set")
	}
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	go registerClient(c, ws, clientId, lobbyId)
	return nil
}

func main() {
	var e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(o string) (bool, error) { return checkOrigin(o), nil },
	}))
	println("starting server")
	e.Logger.SetLevel(log.DEBUG)
	e.GET("/lobby/:lobbyId/ws", handleWs)
	e.Logger.Fatal(e.Start(":8080"))
}
