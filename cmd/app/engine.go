package main

import (
	"net/http"
	"strconv"

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
	dirs, err := a.DirectoryNode.Scan(a.Config.Target[0])
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
