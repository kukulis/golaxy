package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()

	router.LoadHTMLFiles("pages/test2.html")

	router.GET("/test2", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test2.html", gin.H{
			"title": "Main website",
		})
	})

	//router.GET("/test2", func(c *gin.Context) { c.File("./pages/test2.html") })

	_ = router.Run(":8088")
}
