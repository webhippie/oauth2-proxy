package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/webhippie/oauth2-proxy/assets"
	"github.com/webhippie/oauth2-proxy/config"
	"github.com/webhippie/oauth2-proxy/router/middleware/header"
	"github.com/webhippie/oauth2-proxy/router/middleware/logger"
	"github.com/webhippie/oauth2-proxy/router/middleware/recovery"
	"github.com/webhippie/oauth2-proxy/templates"
)

// Load initializes the routing of the application.
func Load(middleware ...gin.HandlerFunc) http.Handler {
	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.SetHTMLTemplate(
		templates.Load(),
	)

	e.Use(middleware...)
	e.Use(logger.SetLogger())
	e.Use(recovery.SetRecovery())
	e.Use(header.SetCache())
	e.Use(header.SetOptions())
	e.Use(header.SetSecure())
	e.Use(header.SetVersion())

	e.StaticFS(
		"/oauth2-proxy/assets",
		assets.Load(),
	)

	e.NoRoute(func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{},
		)
	})

	if config.Server.Pprof {
		pprof.Register(
			e,
			&pprof.Options{
				RoutePrefix: "/oauth2-proxy/debug/pprof",
			},
		)
	}

	return e
}
