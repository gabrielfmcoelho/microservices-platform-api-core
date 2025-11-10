package domain

import (
	"context"

	"gorm.io/gorm"
)

// MANY TO ONE WITH ORGANIZATION

type User struct {
	gorm.Model
	Name           string           `gorm:"size:255"`
	Email          string           `gorm:"size:255;uniqueIndex;not null"`
	Password       string           `gorm:"size:255;not null"`
	OrganizationID uint             `gorm:"not nul;Index"`
	Organization   Organization     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Relationship to Organization
	RoleID         uint             `gorm:"not null;Index"`
	Role           UserRole         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"` // Relationship to UserRole
	Bio            UserBio          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Metrics        UserMetrics      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Configs        UserConfig       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Logs           []UserLog        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ServiceLogs    []UserServiceLog `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CreateUser struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required"`
	OrganizationID uint   `json:"organization_id" binding:"required"`
	RoleID         uint   `json:"role" binding:"required"`
}

type PublicUser struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	OrganizationID   uint   `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
	RoleID           uint   `json:"role_id"`
	RoleName         string `json:"role_name"`
	CreatedAt        string `json:"created_at"`
	LastLogin        string `json:"last_login"`
	IsArchived       bool   `json:"is_archived"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Fetch(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id uint) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, userID uint, user *User) error
	Archive(ctx context.Context, userID uint) error
	Unarchive(ctx context.Context, userID uint) error
}

type UserUsecase interface {
	Create(ctx context.Context, user *CreateUser) error
	Fetch(ctx context.Context) ([]PublicUser, error)
	GetByIdentifier(ctx context.Context, identifier string) (PublicUser, error)
	Update(ctx context.Context, userID uint, user *User) error
	Archive(ctx context.Context, userID uint) error
	Unarchive(ctx context.Context, userID uint) error
}

// EXEMPLE TIP: To access Bio from a User, use the following:

// var user domain.User
// err := db.Preload("Bio").First(&user, 1).Error
// if err != nil {
// 	return err
// }
