package router

import (
	"net/http"
	"reflect"

	"github.com/Dhar01/incident_resp/config"
	"github.com/Dhar01/incident_resp/handler"
	"github.com/Dhar01/incident_resp/internal/model"
	"github.com/Dhar01/incident_resp/lib/middleware"
	"github.com/Dhar01/incident_resp/lib/renderer"
	auth_gen "github.com/Dhar01/incident_resp/router/auth"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type testAPI struct {
}

var _ auth_gen.ServerInterface = (*testAPI)(nil)

func newTestAPI() *testAPI {
	return &testAPI{}
}

func (api *testAPI) CreateUserAuth(c *gin.Context) {
	// delete existing auth cookie if present
	_, errAccessJWT := c.Cookie("accessJWT")
	_, errRefreshJWT := c.Cookie("refreshJWT")
	if errAccessJWT == nil || errRefreshJWT == nil {
		configSecurity := config.GetConfig().Security
		c.SetCookie(
			"accessJWT",
			"",
			-1,
			configSecurity.AuthCookiePath,
			configSecurity.AuthCookieDomain,
			configSecurity.AuthCookieSecure,
			configSecurity.AuthCookieHTTPOnly,
		)
		c.SetCookie(
			"refreshJWT",
			"",
			-1,
			configSecurity.AuthCookiePath,
			configSecurity.AuthCookieDomain,
			configSecurity.AuthCookieSecure,
			configSecurity.AuthCookieHTTPOnly,
		)
	}

	// verify that RDBMS is enabled in .env
	if !config.IsRDBMS() {
		renderer.Render(c, gin.H{"message": "relational database not enabled"}, http.StatusNotImplemented)
		return
	}

	// var auth model.Auth
	var auth model.AuthReq

	// bind JSON
	if err := c.ShouldBindJSON(&auth); err != nil {
		renderer.Render(c, gin.H{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, statusCode := handler.CreateUserAuth(model.Auth{
		Email:    string(auth.Email),
		Password: auth.Password,
	})

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		renderer.Render(c, resp, statusCode)
		return
	}

	renderer.Render(c, resp.Message, statusCode)
}

func (api *testAPI) LogIn(c *gin.Context) {
	// verify that RDBMS is enabled in .env
	if !config.IsRDBMS() {
		renderer.Render(c, gin.H{"message": "relational database not enabled"}, http.StatusNotImplemented)
		return
	}

	// // verify that JWT service is enabled in .env
	// if !config.IsJWT() {
	// 	renderer.Render(c, gin.H{"message": "JWT service not enabled"}, http.StatusNotImplemented)
	// 	return
	// }

	var payload model.AuthPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		renderer.Render(c, gin.H{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	resp, statusCode := handler.Login(payload)

	// auth verification failed
	if statusCode != http.StatusOK {
		renderer.Render(c, resp, statusCode)
		return
	}

	// auth verification OK
	// set cookie if the feature is enabled in app settings
	configSecurity := config.GetConfig().Security
	if configSecurity.AuthCookieActivate {
		tokens, ok := resp.Message.(middleware.JWTPayload)
		if ok {
			c.SetSameSite(configSecurity.AuthCookieSameSite)
			c.SetCookie(
				"accessJWT",
				tokens.AccessJWT,
				middleware.JWTParams.AccessKeyTTL*60,
				configSecurity.AuthCookiePath,
				configSecurity.AuthCookieDomain,
				configSecurity.AuthCookieSecure,
				configSecurity.AuthCookieHTTPOnly,
			)
			c.SetCookie(
				"refreshJWT",
				tokens.RefreshJWT,
				middleware.JWTParams.RefreshKeyTTL*60,
				configSecurity.AuthCookiePath,
				configSecurity.AuthCookieDomain,
				configSecurity.AuthCookieSecure,
				configSecurity.AuthCookieHTTPOnly,
			)

			if !configSecurity.ServeJwtAsResBody {
				resp.Message = "login successful"
				if configSecurity.Must2FA == config.Activated {
					tokens.AccessJWT = ""
					tokens.RefreshJWT = ""
					resp.Message = tokens
				}
			}
		}

		if !ok {
			log.Error("error code: 1011.1")
			resp.Message = "failed to prepare auth cookie"
			statusCode = http.StatusInternalServerError
		}
	}

	if reflect.TypeOf(resp.Message).Kind() == reflect.String {
		renderer.Render(c, resp, statusCode)
		return
	}

	renderer.Render(c, resp.Message, statusCode)
}
