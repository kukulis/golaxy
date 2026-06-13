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

func (c *WsController) ServeWs(ctx *gin.Context) {
	token := ctx.Query("token")
	if !c.authenticationManager.TokenValid(token) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	ws.ServeWs(c.hub, ctx.Writer, ctx.Request)
}
