package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/got-many-wheels/lemari/internal/config"
	directorynode "github.com/got-many-wheels/lemari/internal/directory_node"
	"github.com/got-many-wheels/lemari/internal/renderer"
	"github.com/got-many-wheels/lemari/views"
)

type app struct {
	Config        config.Config
	Engine        *gin.Engine
	DirectoryNode *directorynode.DirectoryNode
}

func newApp() (*app, error) {
	a := &app{
		DirectoryNode: directorynode.New(),
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	a.Config = *cfg

	// TODO: scan the transcoded media instead
	pwd, _ := os.Getwd()
	dirs, err := a.DirectoryNode.Scan(path.Join(pwd, "output"))
	if err != nil {
		return nil, err
	}
	a.DirectoryNode = dirs

	a.Engine = gin.Default()
	a.Engine.SetTrustedProxies(nil)

	// integrate gin with a-h's templ renderer
	ginHtmlRenderer := a.Engine.HTMLRender
	a.Engine.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	a.Engine.Static("public", path.Join(pwd, "public"))

	a.setupRoutes()

	return a, nil
}

func (a *app) setupRoutes() {
	a.Engine.GET("/", func(c *gin.Context) {
		dirs := a.DirectoryNode.DirFiles()
		pwd, _ := os.Getwd()
		// use lexically relative path equivalent to pwd/output to keep it clean
		for i := range dirs {
			rel, _ := filepath.Rel(path.Join(pwd, "output"), dirs[i])
			parts := strings.Split(rel, string(os.PathSeparator))
			dirs[i] = parts[0] // get the video title instead
		}
		c.HTML(http.StatusOK, "", views.Index(dirs))
	})

	a.Engine.GET("/:media", func(c *gin.Context) {
		media := c.Param("media")
		pwd, _ := os.Getwd()
		_, err := os.Stat(path.Join(pwd, "output", media))
		if err != nil {
			if errors.Is(os.ErrNotExist, err) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Media not found",
				})
				return
			}
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Something went wrong",
			})
			fmt.Println(err)
			return
		}
		c.HTML(http.StatusOK, "", views.Media(media))
	})

	a.Engine.GET("/manifest/:media/:manifest", func(c *gin.Context) {
		media := c.Param("media")
		manifest := c.Param("manifest")
		pwd, _ := os.Getwd()
		manifestPath := path.Join(pwd, "output", media, manifest)
		fmanifest, err := os.Open(manifestPath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Manifest not found",
			})
			return
		}
		defer fmanifest.Close()
		c.File(manifestPath)
	})
}

func (a *app) run() error {
	return a.Engine.Run(":" + strconv.Itoa(a.Config.Port))
}
