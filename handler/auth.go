package handler

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/Dhar01/incident_resp/config"
	"github.com/Dhar01/incident_resp/internal/database"
	"github.com/Dhar01/incident_resp/internal/model"
	"github.com/Dhar01/incident_resp/lib"
	"github.com/Dhar01/incident_resp/service"
	"github.com/pilinux/crypt"

	log "github.com/sirupsen/logrus"
)

var errInternalServer string = "internal server error"

// CreateUserAuth receives tasks from controller.CreateUserAuth.
// After email validation, it creates a new user account. It
// supports both the legacy way of saving user email in plaintext
// and the recommended way of applying encryption at rest.
func CreateUserAuth(auth model.Auth) (httpResponse model.HTTPResponse, httpStatusCode int) {

	db := database.GetDB()

	// user must not be able to manipulate all fields
	authFinal := new(model.Auth)
	authFinal.Email = auth.Email
	authFinal.Password = auth.Password

	// email validation
	if !lib.ValidateEmail(auth.Email) {
		return setMessage("wrong email address", http.StatusBadRequest)
	}

	// for backward compatibility
	// email must be unique
	err := db.Where("email = ?", auth.Email).First(&auth).Error
	if err != nil {
		if err.Error() != database.RecordNotFound {
			// db read error
			log.WithError(err).Error("error code: 1002.1")
			return setMessage(errInternalServer, http.StatusInternalServerError)
		}
	}

	if err == nil {
		return setMessage("email already registered", http.StatusBadRequest)
	}

	// downgrade must be avoided to prevent creating duplicate accounts
	// valid: non-encryption mode -> upgrade to encryption mode
	// invalid: encryption mode -> downgrade to non-encryption mode
	if !config.IsCipher() {

		err := db.Where("email_hash IS NOT NULL AND email_hash != ?", "").First(&auth).Error
		if err != nil {
			if err.Error() != database.RecordNotFound {
				// db read error
				log.WithError(err).Error("error code: 1002.2")
				return setMessage(errInternalServer, http.StatusInternalServerError)
			}
		}

		if err == nil {
			e := errors.New("check env: ACTIVATE_CIPHER")
			log.WithError(e).Error("error code: 1002.3")
			return setMessage(errInternalServer, http.StatusInternalServerError)
		}
	}

	// generate a fixed-sized BLAKE2b-256 hash of the email, used for auth purpose
	// when encryption at rest is used
	if config.IsCipher() {
		var err error

		// hash of the email in hexadecimal string format
		emailHash, err := service.CalcHash(
			[]byte(auth.Email),
			config.GetConfig().Security.Blake2bSec,
		)
		if err != nil {
			log.WithError(err).Error("error code: 1001.1")
			return setMessage(errInternalServer, http.StatusInternalServerError)
		}

		authFinal.EmailHash = hex.EncodeToString(emailHash)

		// email must be unique
		err = db.Where("email_hash = ?", authFinal.EmailHash).First(&auth).Error
		if err != nil {
			if err.Error() != database.RecordNotFound {
				// db read error
				log.WithError(err).Error("error code: 1002.4")
				return setMessage(errInternalServer, http.StatusInternalServerError)
			}
		}

		if err == nil {
			return setMessage("email already registered", http.StatusBadRequest)
		}
	}

	// send a verification email if required by the application
	emailDelivered, err := service.SendEmail(authFinal.Email, model.EmailTypeVerifyEmailNewAcc)
	if err != nil {
		log.WithError(err).Error("error code: 1002.5")
		return setMessage("email delivery service failed", http.StatusInternalServerError)
	}

	if emailDelivered {
		authFinal.VerifyEmail = model.EmailNotVerified
	}

	// encryption at rest for user email, mainly needed by system in future
	// to send verification or password recovery emails
	if config.IsCipher() {
		// encrypt the email
		cipherEmail, nonce, err := crypt.EncryptChacha20poly1305(
			config.GetConfig().Security.CipherKey,
			auth.Email,
		)

		if err != nil {
			log.WithError(err).Error("error code: 1001.2")
			return setMessage(errInternalServer, http.StatusInternalServerError)
		}

		// save email only in ciphertext
		authFinal.Email = ""
		authFinal.EmailCipher = hex.EncodeToString(cipherEmail)
		authFinal.EmailNonce = hex.EncodeToString(nonce)
	}

	// one unique email for each account
	tx := db.Begin()
	if err := tx.Create(&authFinal).Error; err != nil {
		tx.Rollback()
		log.WithError(err).Error("error code: 1001.3")
		return setMessage(errInternalServer, http.StatusInternalServerError)
	}

	tx.Commit()

	httpResponse.Message = *authFinal
	httpStatusCode = http.StatusCreated

	return
}

func setMessage(message string, statusCode int) (httpResponse model.HTTPResponse, httpStatusCode int) {
	httpResponse.Message = message
	httpStatusCode = statusCode
	return
}
