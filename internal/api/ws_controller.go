package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"glaktika.eu/galaktika/internal/ws"
)

type WsController struct {
	authenticationManager AuthenticationManager
	hub                   *ws.Hub
}

func NewWsController(authenticationManager AuthenticationManager, hub *ws.Hub) *WsController {
	return &WsController{authenticationManager: authenticationManager, hub: hub}
}

func (c *WsController) ServeWebSocketConnection(ctx *gin.Context) {
	token := ctx.Query("token")
	if !c.authenticationManager.TokenValid(token) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ws.RegisterWebSocketClient(c.hub, ctx.Writer, ctx.Request, token)
}
