package sqldb

import (
	"auth_service/db/mdata"
	"auth_service/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateRandomAccount(t *testing.T) *mdata.Auth {

	password, err := util.HashPassword(util.RandomString(7))

	profilePicture := util.RandomString(5)
	require.NoError(t, err)
	//
	verifyUser := &SignUpParams{
		Username:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		Password:       password,
		Country:        "Myanmar",
		BrowserName:    "Chrome",
		DeviceType:     "Mobile",
		ProfilePicture: &profilePicture,
	}

	user, err := testStore.SignUp(verifyUser)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, verifyUser.Username, user.Username)
	require.Equal(t, verifyUser.Email, user.Email)
	require.Equal(t, verifyUser.Password, user.Password)
	require.Equal(t, verifyUser.Country, user.Country)
	require.Equal(t, verifyUser.BrowserName, user.BrowserName)
	require.Equal(t, verifyUser.DeviceType, user.DeviceType)
	require.Equal(t, verifyUser.ProfilePicture, user.ProfilePicture)

	return user

}

func TestSignUp(t *testing.T) {
	CreateRandomAccount(t)
}

func TestVerifyEmail(t *testing.T) {
	user := CreateRandomAccount(t)

	verifyUser, err := testStore.VerifyEmail(*user.EmailVerificationToken)

	require.NoError(t, err)
	require.NotEmpty(t, verifyUser)

	require.Equal(t, verifyUser.Username, user.Username)
	require.Equal(t, verifyUser.Email, user.Email)
	require.Equal(t, verifyUser.Password, user.Password)
	require.Equal(t, verifyUser.Country, user.Country)
	require.Equal(t, verifyUser.BrowserName, user.BrowserName)
	require.Equal(t, verifyUser.DeviceType, user.DeviceType)
	require.Equal(t, verifyUser.ProfilePicture, user.ProfilePicture)
}
