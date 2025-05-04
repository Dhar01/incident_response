package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pilinux/gorest/config"
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
