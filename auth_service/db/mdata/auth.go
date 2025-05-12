package mdata

import (
	"github.com/google/uuid"
	"time"
)

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

type Session struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	Username     string    `gorm:"type:varchar(255);not null;index;uniqueIndex:index_sessions_username"`
	RefreshToken string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	UserAgent    string    `gorm:"type:varchar(255);not null"`
	ClientIP     string    `gorm:"type:varchar(255);not null"`
	IsBlocked    bool      `gorm:"not null;default:false"`
	LastUsedAt   time.Time `gorm:"not null;default:(CURRENT_TIMESTAMP)"`
	CreatedAt    time.Time `gorm:"not null;autoCreateTime"`
	ExpiredAt    time.Time `gorm:"not null"`

	// Foreign key relationship
	Auth Auth `gorm:"foreignKey:Username;references:Username;constraint:OnUpdate:CASCADE,onDelete:CASCADE"`
}
