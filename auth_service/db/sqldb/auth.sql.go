package sqldb

import (
	"auth_service/db/mdata"
	"auth_service/util"
	"fmt"
	"github.com/google/uuid"
	"time"
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

	tx := sql.Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
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

	if err := tx.Table("auths").Select(cols).Where("email_verification_token = ?", token).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	updateMap := map[string]any{
		"email_verified":           true,
		"email_verification_token": nil,
	}

	if err := tx.Table("auths").Where("id = ?", user.ID).Updates(updateMap).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Use transaction for final query
	var updatedUser mdata.Auth
	if err := tx.Table("auths").Select(cols).Where("id = ?", user.ID).First(&updatedUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil

}

func (sql *SqlDB) GetUser(email string) (*mdata.Auth, error) {
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

	err := sql.FetchOne(&user, "auths", cols, "email=?", email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

type VerifyOTPParams struct {
	OTP         string `json:"otp"`
	BrowserName string `json:"browserName"`
	DeviceType  string `json:"deviceType"`
}

func (sql *SqlDB) VerifyOTP(arg VerifyOTPParams) (*mdata.Auth, error) {

	tx := sql.Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
		"otp",
		"otp_expiration_date",
	}

	// Use transaction for first query

	if err := tx.Table("auths").Select(cols).Where("otp = ?", arg.OTP).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if user.OtpExpirationDate != nil && user.OtpExpirationDate.Before(time.Now()) {
		tx.Rollback()
		return nil, fmt.Errorf("otp expired on %s", user.OtpExpirationDate)
	}

	// Update user: clear OTP and update browser/device info

	updateMap := map[string]any{
		"otp":                 nil,
		"otp_expiration_date": nil,
		"browser_name":        arg.BrowserName,
		"device_type":         arg.DeviceType,
	}

	if err := tx.Table("auths").Where("id = ?", user.ID).Updates(updateMap).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	// Use transaction for final query
	var updatedUser mdata.Auth
	if err := tx.Table("auths").Select(cols).Where("id = ?", user.ID).First(&updatedUser).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

type SignInParams struct {
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	UserAgent string    `json:"userAgent"`
	ClientIP  string    `json:"clientIP"`
	ExpiredAt time.Time `json:"expiredAt"`
}

type SignInResult struct {
	User    *mdata.Auth    `json:"user"`
	Session *mdata.Session `json:"session"`
}

func (sql *SqlDB) SignIn(arg SignInParams) (*SignInResult, error) {

	user, err := sql.GetUser(arg.Email)

	if err != nil {

		return nil, err
	}

	if user.BrowserName != arg.UserAgent || user.DeviceType != arg.ClientIP {
		// for testing verifyOTP DB operation
		//if true {
		// Generate OTP
		otpCode, err := util.GenerateOTP(4)

		if err != nil {
			return nil, err
		}
		// Set OTP expiration (10 minutes from now)
		otpExpiration := time.Now().Add(10 * time.Minute)

		// Update user with OTP
		updateMap := map[string]any{
			"otp":                 otpCode,
			"otp_expiration_date": otpExpiration,
		}

		var auth mdata.Auth
		err = sql.SingleUpdateDB(&auth, updateMap, "id=?", user.ID)
		if err != nil {
			return nil, err
		}

		// Return result without session for OTP verification
		return &SignInResult{
			User:    user,
			Session: nil, // No session created yet, waiting for OTP verification
		}, nil
	}
	//Create Session

	session := &mdata.Session{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: uuid.New().String(),
		UserAgent:    arg.UserAgent,
		ClientIP:     arg.ClientIP,
		IsBlocked:    false,
		LastUsedAt:   time.Now(),
		ExpiredAt:    arg.ExpiredAt,
	}

	// Insert session into database
	err = sql.InsertDB(session)
	if err != nil {
		return nil, err
	}

	return &SignInResult{
		User:    user,
		Session: session,
	}, nil
}

type UpdatePasswordTokenParams struct {
	AuthID          uint      `json:"authID"`
	Token           string    `json:"token"`
	TokenExpiration time.Time `json:"tokenExpiration"`
}

func (sql *SqlDB) UpdatePasswordToken(arg *UpdatePasswordTokenParams) error {
	tx := sql.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fixed: Added table name and proper model reference
	result := tx.Table("auths").
		Where("id = ?", arg.AuthID).
		Updates(map[string]interface{}{
			"password_reset_token":   arg.Token,
			"password_reset_expires": arg.TokenExpiration,
		})

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no user found with ID: %d", arg.AuthID)
	}

	return tx.Commit().Error
}

func (sql *SqlDB) GetUserByPasswordResetToken(token string) (*mdata.Auth, error) {
	var user mdata.Auth
	cols := []string{
		"id",
		"username",
		"email",
		"password",
		"password_reset_token",
		"password_reset_expires",
	}

	err := sql.FetchOne(&user, "auths", cols, "password_reset_token = ? AND password_reset_expires > ?", token, time.Now())
	if err != nil {
		return nil, err
	}

	return &user, nil
}

type UpdatePasswordParams struct {
	AuthID   uint   `json:"authId"`
	Password string `json:"password"`
}

func (sql *SqlDB) UpdatePassword(arg *UpdatePasswordParams) error {
	updateMap := map[string]any{
		"password":               arg.Password,
		"password_reset_token":   nil, // Clear the reset token
		"password_reset_expires": nil, // Clear the expiration
	}

	var auth mdata.Auth

	err := sql.SingleUpdateDB(&auth, updateMap, "id=?", arg.AuthID)

	if err != nil {
		return err
	}
	return nil
}
