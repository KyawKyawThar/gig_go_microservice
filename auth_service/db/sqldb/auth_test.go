package sqldb

import (
	"auth_service/db/mdata"
	"auth_service/util"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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

func TestGetUser(t *testing.T) {
	user := CreateRandomAccount(t)

	getUser, err := testStore.GetUser(user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, getUser)

	require.Equal(t, user.Username, getUser.Username)
	require.Equal(t, user.Email, getUser.Email)
	require.Equal(t, user.Password, getUser.Password)
	require.Equal(t, user.Country, getUser.Country)
	require.Equal(t, user.BrowserName, getUser.BrowserName)
	require.Equal(t, user.DeviceType, getUser.DeviceType)
	require.Equal(t, user.ProfilePicture, getUser.ProfilePicture)

}

func TestSignIn(t *testing.T) {
	user := CreateRandomAccount(t)

	signInParams := SignInParams{
		Email:     user.Email,
		Password:  user.Password,
		UserAgent: "Chrome",
		ClientIP:  "127.0.0.1",
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	result, err := testStore.SignIn(signInParams)

	require.NoError(t, err)
	require.NotEmpty(t, result)
	require.NotEmpty(t, result.User)
	//require.NotEmpty(t, result.Session)
}

func TestVerifyOTP(t *testing.T) {
	// TODO: please change condition true in SignIn function to test this function
	otpParams := VerifyOTPParams{
		OTP:         "7735",
		BrowserName: "Firefox",
		DeviceType:  "Web",
	}
	// Test OTP verification
	verifiedUser, err := testStore.VerifyOTP(otpParams)
	require.NoError(t, err)
	require.NotEmpty(t, verifiedUser)

}

func TestGetUserByPasswordResetToken(t *testing.T) {
	user := CreateRandomAccount(t)
	token := util.RandomString(20)
	tokenParams := UpdatePasswordTokenParams{
		AuthID:          user.ID,
		Token:           token,
		TokenExpiration: time.Now().Add(1 * time.Hour),
	}
	err := testStore.UpdatePasswordToken(&tokenParams)
	require.NoError(t, err)

	foundUser, err := testStore.GetUserByPasswordResetToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, foundUser)
	require.Equal(t, user.ID, foundUser.ID)
}

func TestUpdatePassword(t *testing.T) {
	user := CreateRandomAccount(t)

	newPassword := util.RandomString(10)
	hashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	updateParams := UpdatePasswordParams{
		AuthID:   user.ID,
		Password: hashedPassword,
	}

	// Test updating password
	err = testStore.UpdatePassword(&updateParams)
	require.NoError(t, err)
}

func TestGetCurrentUserById(t *testing.T) {
	user := CreateRandomAccount(t)

	// Test getting current user by ID
	currentUser, err := testStore.GetCurrentUserById(fmt.Sprintf("%d", user.ID))
	require.NoError(t, err)
	require.NotEmpty(t, currentUser)
	require.Equal(t, user.ID, currentUser.ID)
	require.Equal(t, user.Username, currentUser.Username)
}
