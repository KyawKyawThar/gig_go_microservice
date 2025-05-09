package sqldb

import (
	"auth_service/db/mdata"
	"auth_service/util"
	"github.com/google/uuid"
)

type SignUpParams struct {
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	Country        string  `json:"country"`
	BrowserName    string  `json:"browserName"`
	DeviceType     string  `json:"deviceType"`
	ProfilePicture *string `json:"profilePicture"`
}

func (sql *SqlDB) SignUp(arg *SignUpParams) (*mdata.Auth, error) {

	profilePublicId := uuid.New().String()
	emailVarifyToken, err := util.RandomHex(20)
	if err != nil {
		return nil, err
	}
	user := &mdata.Auth{
		Username:               arg.Username,
		Email:                  arg.Email,
		Password:               arg.Password,
		Country:                arg.Country,
		ProfilePublicId:        profilePublicId,
		BrowserName:            arg.BrowserName,
		DeviceType:             arg.DeviceType,
		ProfilePicture:         arg.ProfilePicture,
		EmailVerificationToken: &emailVarifyToken,
	}
	err = sql.InsertDB(user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (sql *SqlDB) VerifyEmail(token string) (*mdata.Auth, error) {
	var user mdata.Auth

	cols := []string{
		"id",
		"username",
		"email",
		"password",
		"country",
		"browser_name",
		"device_type",
		"profile_picture",
		"email_verification_token",
		"email_verified",
	}

	err := sql.FetchOne(&user, "auths", cols, "email_verification_token=?", token)

	if err != nil {
		return nil, err
	}

	updateMap := map[string]any{
		"email_verified":           true,
		"email_verification_token": nil,
	}

	err = sql.SingleUpdateDB(&user, updateMap, "id=?", user.ID)
	if err != nil {
		return nil, err
	}
	var updatedUser mdata.Auth
	err = sql.FetchOne(&updatedUser, "auths", cols, "id=?", user.ID)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
