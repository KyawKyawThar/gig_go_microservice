package sqldb

import (
	"auth_service/db/mdata"
)

// Querier defines all query methods
type Querier interface {
	DBMigrate
	AuthQuerier
	UserQuerier
}

type DBMigrate interface {
	MigrateDB() error
}

type AuthQuerier interface {
	SignUp(arg *SignUpParams) (*mdata.Auth, error)
	VerifyEmail(token string) (*mdata.Auth, error)
	GetUser(email string) (*mdata.Auth, error)
	VerifyOTP(arg VerifyOTPParams) (*mdata.Auth, error)
	SignIn(arg SignInParams) (*SignInResult, error)
	UpdatePasswordToken(arg *UpdatePasswordTokenParams) error
	GetUserByPasswordResetToken(token string) (*mdata.Auth, error)
	UpdatePassword(arg *UpdatePasswordParams) error
}

type UserQuerier interface {
	GetCurrentUserById(authId string) (*mdata.Auth, error)
}
