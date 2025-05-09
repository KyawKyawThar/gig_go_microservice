package sqldb

import (
	"auth_service/db/mdata"
)

// Querier defines all query methods
type Querier interface {
	MigrateDB(model any) error
	SignUp(arg *SignUpParams) (*mdata.Auth, error)
	VerifyEmail(token string) (*mdata.Auth, error)
}
