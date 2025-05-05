package router

import (
	"github.com/Dhar01/incident_resp/config"
	"github.com/gin-gonic/gin"
)

func SetUpRouter(configure config.Configuration) (*gin.Engine, error) {
	if config.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		return router, err
	}

	return router, nil
}
