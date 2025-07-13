package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex"`
	Role      string         `json:"role" gorm:"default:'user'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Session struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	RefreshToken string    `json:"-" gorm:"uniqueIndex;not null"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProxyLog struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	UserID     *uint     `json:"user_id"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	RemoteAddr string    `json:"remote_addr"`
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	Host       string    `json:"host"`
	UserAgent  string    `json:"user_agent"`
	StatusCode int       `json:"status_code"`
	ResponseSize int64   `json:"response_size"`
	Duration   int64     `json:"duration"` // in milliseconds
	Timestamp  time.Time `json:"timestamp"`
}

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}