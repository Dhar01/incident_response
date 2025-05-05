package handler

import (
	"net/http"
	"strings"

	"github.com/Dhar01/incident_resp/config"
	"github.com/Dhar01/incident_resp/internal/database"
	"github.com/Dhar01/incident_resp/internal/model"
	"github.com/Dhar01/incident_resp/lib"
	"github.com/Dhar01/incident_resp/lib/middleware"
	"github.com/Dhar01/incident_resp/service"
	"github.com/pilinux/argon2"

	log "github.com/sirupsen/logrus"
)

// Login receives tasks from controller.Login.
// After authentication, it returns new access and refresh tokens.
func Login(payload model.AuthPayload) (httpResponse model.HTTPResponse, httpStatusCode int) {
	payload.Email = strings.TrimSpace(payload.Email)
	if !lib.ValidateEmail(payload.Email) {
		return setMessage("wrong email address", http.StatusBadRequest)
	}

	v, err := service.GetUserByEmail(payload.Email, false)
	if err != nil {
		if err.Error() != database.RecordNotFound {
			// db read error
			log.WithError(err).Error("error code: 1013.1")
			return setMessage(errInternalServer, http.StatusInternalServerError)
		}

		return setMessage("email not found", http.StatusNotFound)
	}

	// app settings
	configSecurity := config.GetConfig().Security

	// check whether email verification is required
	if configSecurity.VerifyEmail {
		if v.VerifyEmail != model.EmailVerified {
			return setMessage("email verification required", http.StatusUnauthorized)
		}
	}

	verifyPass, err := argon2.ComparePasswordAndHash(payload.Password, configSecurity.HashSec, v.Password)
	if err != nil {
		log.WithError(err).Error("error code: 1013.2")
		return setMessage(errInternalServer, http.StatusInternalServerError)
	}

	if !verifyPass {
		return setMessage("wrong credentials", http.StatusUnauthorized)
	}

	// custom claims
	claims := middleware.MyCustomClaims{}
	claims.AuthID = v.AuthID
	// claims.Email
	// claims.Role
	// claims.Scope
	// claims.TwoFA
	// claims.SiteLan
	// claims.Custom1
	// claims.Custom2

	// when 2FA is enabled for this application (ACTIVATE_2FA=yes)
	if configSecurity.Must2FA == config.Activated {
		db := database.GetDB()
		twoFA := model.TwoFA{}

		// have the user configured 2FA
		err := db.Where("id_auth = ?", v.AuthID).First(&twoFA).Error
		if err != nil {
			if err.Error() != database.RecordNotFound {
				// db read error
				log.WithError(err).Error("error code: 1013.3")
				return setMessage(errInternalServer, http.StatusInternalServerError)
			}
		}
		if err == nil {
			claims.TwoFA = twoFA.Status

			// 2FA ON
			if twoFA.Status == configSecurity.TwoFA.Status.On {
				// hash user's pass
				hashPass, err := service.GetHash([]byte(payload.Password))
				if err != nil {
					log.WithError(err).Error("error code: 1013.4")
					return setMessage(errInternalServer, http.StatusInternalServerError)
				}

				// save the hashed pass in memory for OTP validation step
				data2FA := model.Secret2FA{}
				data2FA.PassSHA = hashPass
				model.InMemorySecret2FA[claims.AuthID] = data2FA
			}
		}
	}

	// // issue new tokens
	// accessJWT, _, err := middleware.GetJWT(claims, "access")
	// if err != nil {
	// 	log.WithError(err).Error("error code: 1013.5")
	// 	return setMessage(errInternalServer, http.StatusInternalServerError)
	// }
	// refreshJWT, _, err := middleware.GetJWT(claims, "refresh")
	// if err != nil {
	// 	log.WithError(err).Error("error code: 1013.6")
	// 	return setMessage(errInternalServer, http.StatusInternalServerError)
	// }

	// jwtPayload := middleware.JWTPayload{}
	// jwtPayload.AccessJWT = accessJWT
	// jwtPayload.RefreshJWT = refreshJWT
	// jwtPayload.TwoAuth = claims.TwoFA

	// httpResponse.Message = jwtPayload
	httpResponse.Message = "ok"
	httpStatusCode = http.StatusOK
	return
}
