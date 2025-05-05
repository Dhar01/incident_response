package router

import (
	"github.com/Dhar01/incident_resp/config"
	auth_gen "github.com/Dhar01/incident_resp/router/auth"
	"github.com/gin-gonic/gin"
)

var base string = "/api/v1"

func SetUpRouter(configure config.Configuration) (*gin.Engine, error) {
	if config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// auth routes
	authRoutes(&router.RouterGroup, base)

	if err := router.SetTrustedProxies(nil); err != nil {
		return router, err
	}

	return router, nil
}

func authRoutes(router *gin.RouterGroup, baseURL string) {
	middlewares := []auth_gen.MiddlewareFunc{}

	opt := auth_gen.GinServerOptions{
		BaseURL:     baseURL,
		Middlewares: middlewares,
	}

	api := newTestAPI()

	auth_gen.RegisterHandlersWithOptions(router, api, opt)
}
