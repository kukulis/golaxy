package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

func main() {
	router := gin.Default()

	//var files []string
	//fs.WalkDir(os.DirFS("./pages"), ".", func(path string, d fs.DirEntry, err	error) error {
	//	if !d.IsDir() && strings.HasSuffix(path, ".html") {
	//	files = append(files, filepath.Join("./pages", path))
	//}
	//	return nil
	//})
	//tmpl := template.Must(template.ParseFiles(files...))

	//template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	testTmpl := template.Must(template.ParseFiles("./pages/layout/base.html", "./pages/test.html"))
	test2Tmpl := template.Must(template.ParseFiles("./pages/layout/base.html", "./pages/test2.html"))

	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
		testTmpl.ExecuteTemplate(c.Writer, "test.html", nil)
	})

	router.GET("/test2", func(c *gin.Context) {
		c.Status(http.StatusOK)
		test2Tmpl.ExecuteTemplate(c.Writer, "test2.html", nil)
	})

	_ = router.Run(":8088")
}
