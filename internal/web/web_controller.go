package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type WebController struct {
}

func RenderDivision(c *gin.Context) {
	divisionId := c.Param("divisionId")
	c.HTML(http.StatusOK, "division_main", gin.H{"DivisionId": divisionId})
}
