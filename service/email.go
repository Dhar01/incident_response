package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Dhar01/incident_resp/config"
	"github.com/Dhar01/incident_resp/internal/database"
	"github.com/Dhar01/incident_resp/internal/model"
	"github.com/pilinux/gorest/lib"
	"github.com/google/uuid"
	"github.com/mediocregopher/radix/v4"
	"github.com/pilinux/libgo/timestring"

	log "github.com/sirupsen/logrus"
)

// SendEmail sends a verification/password recovery email if
//
// - required by the application
//
// - an external email service is configured
//
// - a redis database is configured
//
// {true, nil} => email delivered successfully
//
// {false, nil} => email delivery not required/service not configured
//
// {false, error} => email delivery failed
func SendEmail(email string, emailType int, opts ...string) (bool, error) {
	// send email if required by the application
	appConfig := config.GetConfig()

	// is external email service activated
	if appConfig.EmailConf.Activate != config.Activated {
		return false, nil
	}

	// is verification/password recovery email required
	doSendEmail := (appConfig.Security.VerifyEmail &&
		(emailType == model.EmailTypeVerifyEmailNewAcc || emailType == model.EmailTypeVerifyUpdatedEmail)) ||

		(appConfig.Security.RecoverPass && emailType == model.EmailTypePassRecovery)
	if !doSendEmail {
		return false, nil
	}

	// is redis database activated
	if appConfig.Database.REDIS.Activate != config.Activated {
		return false, nil
	}

	data := struct {
		key   string
		value string
	}{}
	var keyTTL uint64
	var emailTag string
	var code uint64
	var codeUUIDv4 string

	// generate verification/password recovery code
	if emailType == model.EmailTypeVerifyEmailNewAcc || emailType == model.EmailTypeVerifyUpdatedEmail {
		if emailType == model.EmailTypeVerifyEmailNewAcc {
			data.key = model.EmailVerificationKeyPrefix
		}

		if emailType == model.EmailTypeVerifyUpdatedEmail {
			data.key = model.EmailUpdateKeyPrefix
		}

		if config.IsEmailVerificationCodeUUIDv4() {
			codeUUIDv4 = uuid.NewString()
			data.key += codeUUIDv4
		}

		if !config.IsEmailVerificationCodeUUIDv4() {
			code = lib.SecureRandomNumber(appConfig.EmailConf.EmailVerificationCodeLength)
			data.key += strconv.FormatUint(code, 10)
		}

		keyTTL = appConfig.EmailConf.EmailVerifyValidityPeriod
		emailTag = appConfig.EmailConf.EmailVerificationTag
	}

	if emailType == model.EmailTypePassRecovery {
		if config.IsPasswordRecoverCodeUUIDv4() {
			codeUUIDv4 = uuid.NewString()
			data.key = model.PasswordRecoveryKeyPrefix + codeUUIDv4
		}

		if !config.IsPasswordRecoverCodeUUIDv4() {
			code = lib.SecureRandomNumber(appConfig.EmailConf.PasswordRecoverCodeLength)
			data.key = model.PasswordRecoveryKeyPrefix + strconv.FormatUint(code, 10)
		}

		keyTTL = appConfig.EmailConf.PassRecoverValidityPeriod
		emailTag = appConfig.EmailConf.PasswordRecoverTag
	}

	data.value = email

	// when encryption at rest is used
	if config.IsCipher() {
		var err error

		// hash of the email in hexadecimal string format
		value, err := CalcHash(
			[]byte(email),
			config.GetConfig().Security.Blake2bSec,
		)
		if err != nil {
			log.WithError(err).Error("error code: 406.1")
			return false, err
		}

		data.value = hex.EncodeToString(value)
	}

	// save in redis with expiry time
	client := *database.GetRedis()
	redisConnTTL := appConfig.Database.REDIS.Conn.ConnTTL

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(redisConnTTL)*time.Second)
	defer cancel()

	// Set key in Redis
	r1 := ""
	if err := client.Do(ctx, radix.FlatCmd(&r1, "SET", data.key, data.value)); err != nil {
		log.WithError(err).Error("error code: 401")
		return false, err
	}

	if r1 != "OK" {
		log.Error("error code: 402")
		return false, errors.New("failed to save in redis")
	}

	// Set expiry time
	r2 := 0
	if err := client.Do(ctx, radix.FlatCmd(&r2, "EXPIRE", data.key, keyTTL)); err != nil {
		log.WithError(err).Error("error code: 403")
	}

	if r2 != 1 {
		log.Error("error code: 404")
	}

	// check which email service
	// for Postmark
	if appConfig.EmailConf.Provider == "postmark" {
		htmlModel := lib.HTMLModel(lib.StrArrHTMLModel(appConfig.EmailConf.HTMLModel))
		if code != 0 {
			htmlModel["secret_code"] = code
		}

		if code == 0 {
			htmlModel["secret_code"] = codeUUIDv4
		}

		htmlModel["email_validity_period"] = timestring.HourMinuteSecond(keyTTL)

		optsLen := len(opts)
		if optsLen > 0 {
			for i := 0; i < optsLen; i++ {
				key := fmt.Sprintf("additional_info_%d", i)
				htmlModel[key] = opts[i]
			}
		}

		params := PostmarkParams{}
		params.ServerToken = appConfig.EmailConf.APIToken

		if emailType == model.EmailTypeVerifyEmailNewAcc {
			params.TemplateID = appConfig.EmailConf.EmailVerificationTemplateID
		}

		if emailType == model.EmailTypePassRecovery {
			params.TemplateID = appConfig.EmailConf.PasswordRecoverTemplateID
		}

		if emailType == model.EmailTypeVerifyUpdatedEmail {
			params.TemplateID = appConfig.EmailConf.EmailUpdateVerifyTemplateID
		}

		params.From = appConfig.EmailConf.AddrFrom
		params.To = email
		params.Tag = emailTag
		params.TrackOpens = appConfig.EmailConf.TrackOpens
		params.TrackLinks = appConfig.EmailConf.TrackLinks
		params.MessageStream = appConfig.EmailConf.DeliveryType
		params.HTMLModel = htmlModel

		// send the email
		res, err := Postmark(params)
		if err != nil {
			log.WithError(err).Error("error code: 405")
			return false, err
		}

		if res.Message != "OK" {
			return false, errors.New("email delivery failed")
		}

		return true, nil
	}

	e := errors.New(
		"email delivery service provider: '" + appConfig.EmailConf.Provider + "' is unknown",
	)

	log.WithError(e).Error("error code: 406")

	return false, e
}
