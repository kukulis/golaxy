package api

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func bearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	return strings.TrimPrefix(header, "Bearer ")
}
