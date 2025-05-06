package router

import (
	"github.com/Dhar01/incident_resp/config"
	"github.com/Dhar01/incident_resp/lib/middleware"
	auth_gen "github.com/Dhar01/incident_resp/router/auth"
	incident_gen "github.com/Dhar01/incident_resp/router/incidents"
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

	// incident routes
	incidentRoutes(&router.RouterGroup, base)

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

func incidentRoutes(router *gin.RouterGroup, baseURL string) {
	middlewares := []incident_gen.MiddlewareFunc{
		incident_gen.MiddlewareFunc(middleware.JWT()),
	}

	opt := incident_gen.GinServerOptions{
		BaseURL:     baseURL,
		Middlewares: middlewares,
	}

	api := newIncidentAPI()

	incident_gen.RegisterHandlersWithOptions(router, api, opt)
}
