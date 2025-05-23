package model

import (
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Email verification statuses
const (
	EmailNotVerified       int8 = -1
	EmailVerifyNotRequired int8 = 0
	EmailVerified          int8 = 1
)

// Email type
const (
	EmailTypeVerifyEmailNewAcc  int = 1 // verify email of newly registered user
	EmailTypePassRecovery       int = 2 // password recovery code
	EmailTypeVerifyUpdatedEmail int = 3 // verify request of updating user email
)

// Redis key prefixes
const (
	EmailVerificationKeyPrefix string = "gorest-email-verification-"
	EmailUpdateKeyPrefix       string = "gorest-email-update-"
	PasswordRecoveryKeyPrefix  string = "gorest-pass-recover-"
)

// type Auth struct {
// 	AuthID    uint64         `gorm:"column:auth_id;primaryKey" json:"authID,omitempty"`
// 	CreatedAt time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
// 	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
// 	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
// 	Email       string `gorm:"column:email" json:"email"`
// 	EmailCipher string `gorm:"column:email_cipher" json:"-"`
// 	EmailNonce  string `gorm:"column:email_nonce" json:"-"`
// 	EmailHash   string `gorm:"column:email_hash" json:"-"`
// 	Password    string `gorm:"column:password_hash" json:"password"`
// 	VerifyEmail int8   `gorm:"column:verify_email" json:"-"`
// }

type Auth struct {
	AuthID    uint64         `gorm:"column:auth_id;primaryKey;autoIncrement" json:"authID,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`

	Email       string `gorm:"column:email;type:varchar(4);not null;default:''" json:"email"`
	EmailCipher string `gorm:"column:email_cipher;type:text" json:"-"`
	EmailNonce  string `gorm:"column:email_nonce;type:text" json:"-"`
	EmailHash   string `gorm:"column:email_hash;type:varchar(255);unique;not null" json:"-"`
	Password    string `gorm:"column:password_hash;type:varchar(255);not null" json:"password"`
	VerifyEmail int8   `gorm:"column:verify_email;type:smallint;default:0" json:"-"`
}


type AuthReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// // UnmarshalJSON ...
// func (v *Auth) UnmarshalJSON(b []byte) error {
// 	aux := struct {
// 		AuthID   uint64 `json:"authID"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}{}
// 	if err := json.Unmarshal(b, &aux); err != nil {
// 		return err
// 	}
// 	configSecurity := config.GetConfig().Security
// 	// check password length
// 	// if more checks are required i.e. password pattern,
// 	// add all conditions here
// 	if len(aux.Password) < configSecurity.UserPassMinLength {
// 		return errors.New("short password")
// 	}
// 	v.AuthID = aux.AuthID
// 	v.Email = strings.TrimSpace(aux.Email)
// 	config := lib.HashPassConfig{
// 		Memory:      configSecurity.HashPass.Memory,
// 		Iterations:  configSecurity.HashPass.Iterations,
// 		Parallelism: configSecurity.HashPass.Parallelism,
// 		SaltLength:  configSecurity.HashPass.SaltLength,
// 		KeyLength:   configSecurity.HashPass.KeyLength,
// 	}
// 	pass, err := lib.HashPass(config, aux.Password, configSecurity.HashSec)
// 	if err != nil {
// 		return err
// 	}
// 	v.Password = pass
// 	return nil
// }

// MarshalJSON ...
func (v Auth) MarshalJSON() ([]byte, error) {
	aux := struct {
		AuthID uint64 `json:"authID"`
		Email  string `json:"email"`
	}{
		AuthID: v.AuthID,
		Email:  strings.TrimSpace(v.Email),
	}

	return json.Marshal(aux)
}

// AuthPayload - struct to handle all auth data
type AuthPayload struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`

	VerificationCode string `json:"verificationCode,omitempty"`

	OTP string `json:"otp,omitempty"`

	SecretCode  string `json:"secretCode,omitempty"`
	RecoveryKey string `json:"recoveryKey,omitempty"`

	PassNew    string `json:"passNew,omitempty"`
	PassRepeat string `json:"passRepeat,omitempty"`
}

// TempEmail - 'temp_emails' table to hold data temporarily
// during the process of replacing a user's email address
// with a new one
type TempEmail struct {
	ID          uint64    `gorm:"primaryKey" json:"-"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
	Email       string    `gorm:"index" json:"emailNew"`
	Password    string    `gorm:"-" json:"password,omitempty"`
	EmailCipher string    `json:"-"`
	EmailNonce  string    `json:"-"`
	EmailHash   string    `gorm:"index" json:"-"`
	IDAuth      uint64    `gorm:"index" json:"-"`
}
