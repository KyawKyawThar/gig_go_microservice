package mdata

import "time"

type Auth struct {
	ID                     uint    `gorm:"primaryKey;autoIncrement"`
	Username               string  `gorm:"type:varchar(255);not null;uniqueIndex"`
	Email                  string  `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password               string  `gorm:"type:varchar(255);not null"`
	Country                string  `gorm:"type:varchar(255);not null"`
	BrowserName            string  `gorm:"type:varchar(255);not null"`
	DeviceType             string  `gorm:"type:varchar(255);not null"`
	ProfilePicture         *string `gorm:"type:varchar(255)"`
	ProfilePublicId        string  `gorm:"type:varchar(255);not null;uniqueIndex"`
	EmailVerificationToken *string `gorm:"type:varchar(255);uniqueIndex"`
	EmailVerified          bool    `gorm:"not null;default:false"`
	Otp                    *string `gorm:"type:varchar(6)"` // Specific length for OTP
	OtpExpirationDate      *time.Time
	PasswordResetToken     *string `gorm:"type:varchar(255)"`
	PasswordResetExpires   *time.Time
	CreatedAt              time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt              time.Time `gorm:"not null;autoUpdateTime"`
}
