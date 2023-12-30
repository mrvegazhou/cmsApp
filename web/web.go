package web

import (
	"cmsApp/configs"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

var StaticsFs http.FileSystem

func Init() error {
	StaticsFs = gin.Dir(configs.App.Upload.BasePath, false)
	return nil
}

func LoadTemplates() (render multitemplate.Renderer, err error) {
	templatesDir := configs.RootPath + "/web/views"
	render = multitemplate.NewRenderer()
	layouts, err := filepath.Glob(templatesDir + "/layout/*.html")
	if err != nil {
		return
	}
	includes, err := filepath.Glob(templatesDir + "/template/*/*.html")
	if err != nil {
		return
	}
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		dirSlice := strings.Split(include, string(filepath.Separator))
		// [main error.html]
		fileName := strings.Join(dirSlice[len(dirSlice)-2:], "/")
		render.AddFromFiles(fileName, files...)
	}
	return
}
