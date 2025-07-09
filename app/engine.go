package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	directorynode "github.com/got-many-wheels/lemari/internal/directory_node"
	"github.com/got-many-wheels/lemari/internal/renderer"
	"github.com/got-many-wheels/lemari/views"
)

type app struct {
	Config        config
	Engine        *gin.Engine
	DirectoryNode *directorynode.DirectoryNode
}

func newApp() (*app, error) {
	a := &app{
		DirectoryNode: directorynode.New(),
	}

	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	a.Config = *cfg

	// scan media folder from `media_path`
	dirs, err := a.DirectoryNode.Scan(a.Config.MediaPath)
	if err != nil {
		return nil, err
	}
	a.DirectoryNode = dirs

	a.Engine = gin.Default()
	a.Engine.SetTrustedProxies(nil)

	// integrate gin with a-h's templ renderer
	ginHtmlRenderer := a.Engine.HTMLRender
	a.Engine.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}

	a.setupRoutes()

	return a, nil
}

func (a *app) setupRoutes() {
	a.Engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "", views.Index())
	})
}

func (a *app) run() error {
	return a.Engine.Run(":" + strconv.Itoa(a.Config.Port))
}
