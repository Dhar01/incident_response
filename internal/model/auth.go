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

// Auth model - `auths` table
type Auth struct {
	AuthID    uint64         `gorm:"primaryKey" json:"authID,omitempty"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email       string `gorm:"type:varchar(100)" json:"email"`
	EmailCipher string `json:"-"`
	EmailNonce  string `json:"-"`
	EmailHash   string `gorm:"uniqueIndex:idx_email_hash" json:"-"`
	Password    string `gorm:"type:varchar(255);not null" json:"password"`
	VerifyEmail int8   `json:"-"`

	CreateIncidents   []Incident `gorm:"foreignKey:AuthID"`
	AssignedIncidents []Incident `gorm:"foreignKey:AssignedTo"`
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
