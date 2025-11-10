package domain

import (
	"context"

	"time"

	"gorm.io/gorm"
)

// MANY TO ONE WITH USER
// MANY TO ONE WITH SERVICE

type UserServiceLog struct {
	gorm.Model
	UserID    uint          `gorm:"not null;Index"`
	ServiceID uint          `gorm:"not null;Index"`
	Duration  time.Duration `gorm:"default:0"` // Duration in seconds;
}

type PublicUserServiceLog struct {
	ID        uint `json:"id"`
	UserID    uint `json:"user_id"`
	ServiceID uint `json:"service_id"`
	Duration  int  `json:"duration"`
}

type UserServiceLogRepository interface {
	Create(ctx context.Context, UserServiceLog *UserServiceLog) error
	Fetch(ctx context.Context) ([]UserServiceLog, error)
	GetByID(ctx context.Context, id uint) (UserServiceLog, error)
	GetByUserID(ctx context.Context, userID uint) (UserServiceLog, error)
	GetByServiceID(ctx context.Context, serviceID uint) (UserServiceLog, error)
	UpdateDuration(ctx context.Context, UserServiceLogID uint, duration int) error
	Delete(ctx context.Context, UserServiceLogID uint) error
	GetUsageStatistics(ctx context.Context, organizationID *uint, startDate *string, endDate *string) (UsageStatistics, error)
}

type UserServiceLogUsecase interface {
	Fetch(ctx context.Context) ([]PublicUserServiceLog, error)
	GetByIdentifier(ctx context.Context, identifier string) (PublicUserServiceLog, error)
	Delete(ctx context.Context, UserServiceLogID uint) error
	GetUsageStatistics(ctx context.Context, organizationID *uint, startDate *string, endDate *string) (UsageStatistics, error)
}
