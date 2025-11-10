package repository

import (
	"context"
	"errors"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"gorm.io/gorm"
)

type userServiceLogRepository struct {
	db *gorm.DB
}

func NewUserServiceLogRepository(db *gorm.DB) domain.UserServiceLogRepository {
	return &userServiceLogRepository{
		db: db,
	}
}

// Create inserts a new UserServiceLog in the database
func (r *userServiceLogRepository) Create(ctx context.Context, userServiceLog *domain.UserServiceLog) error {
	if err := r.db.WithContext(ctx).Create(userServiceLog).Error; err != nil {
		// Adjust to your error handling
		return domain.ErrDataBaseInternalError
	}
	return nil
}

// Fetch returns all UserServiceLog entries
func (r *userServiceLogRepository) Fetch(ctx context.Context) ([]domain.UserServiceLog, error) {
	var logs []domain.UserServiceLog
	if err := r.db.WithContext(ctx).Find(&logs).Error; err != nil {
		return nil, domain.ErrDataBaseInternalError
	}
	return logs, nil
}

// GetByID returns a UserServiceLog by its ID
func (r *userServiceLogRepository) GetByID(ctx context.Context, id uint) (domain.UserServiceLog, error) {
	var log domain.UserServiceLog
	if err := r.db.WithContext(ctx).First(&log, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return log, domain.ErrNotFound
		}
		return log, domain.ErrDataBaseInternalError
	}
	return log, nil
}

// GetByUserID returns a UserServiceLog by user ID
func (r *userServiceLogRepository) GetByUserID(ctx context.Context, userID uint) (domain.UserServiceLog, error) {
	var log domain.UserServiceLog
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return log, domain.ErrNotFound
		}
		return log, domain.ErrDataBaseInternalError
	}
	return log, nil
}

// GetByServiceID returns a UserServiceLog by service ID
func (r *userServiceLogRepository) GetByServiceID(ctx context.Context, serviceID uint) (domain.UserServiceLog, error) {
	var log domain.UserServiceLog
	if err := r.db.WithContext(ctx).Where("service_id = ?", serviceID).First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return log, domain.ErrNotFound
		}
		return log, domain.ErrDataBaseInternalError
	}
	return log, nil
}

// Update updates an existing UserServiceLog, add duration to the existing duration
func (r *userServiceLogRepository) UpdateDuration(ctx context.Context, userServiceLogID uint, duration int) error {
	if err := r.db.WithContext(ctx).Model(&domain.UserServiceLog{}).
		Where("id = ?", userServiceLogID).
		Update("duration", gorm.Expr("duration + ?", duration)).Error; err != nil {
		return domain.ErrDataBaseInternalError
	}
	return nil
}

// Delete removes a UserServiceLog by its ID (hard delete)
func (r *userServiceLogRepository) Delete(ctx context.Context, userServiceLogID uint) error {
	if err := r.db.WithContext(ctx).Delete(&domain.UserServiceLog{}, userServiceLogID).Error; err != nil {
		return domain.ErrDataBaseInternalError
	}
	return nil
}

// GetUsageStatistics returns aggregated usage statistics for admin dashboard
func (r *userServiceLogRepository) GetUsageStatistics(ctx context.Context) (domain.UsageStatistics, error) {
	var stats domain.UsageStatistics

	// Get total unique users
	var totalUsers int64
	if err := r.db.WithContext(ctx).
		Model(&domain.UserServiceLog{}).
		Distinct("user_id").
		Count(&totalUsers).Error; err != nil {
		return stats, domain.ErrDataBaseInternalError
	}
	stats.TotalUsers = int(totalUsers)

	// Get total duration (sum of all durations in nanoseconds, convert to seconds)
	type TotalDurationResult struct {
		TotalNanoseconds int64
	}
	var durationResult TotalDurationResult
	if err := r.db.WithContext(ctx).
		Model(&domain.UserServiceLog{}).
		Select("COALESCE(SUM(duration), 0) as total_nanoseconds").
		Scan(&durationResult).Error; err != nil {
		return stats, domain.ErrDataBaseInternalError
	}
	stats.TotalDuration = int(durationResult.TotalNanoseconds / 1000000000) // Convert nanoseconds to seconds

	// Get per-service statistics
	type ServiceStatsRow struct {
		ServiceID        uint
		ServiceName      string
		TotalUsers       int
		TotalNanoseconds int64
	}
	var serviceRows []ServiceStatsRow
	if err := r.db.WithContext(ctx).
		Table("user_service_logs").
		Select(`
			user_service_logs.service_id,
			services.name as service_name,
			COUNT(DISTINCT user_service_logs.user_id) as total_users,
			COALESCE(SUM(user_service_logs.duration), 0) as total_nanoseconds
		`).
		Joins("LEFT JOIN services ON services.id = user_service_logs.service_id").
		Where("user_service_logs.deleted_at IS NULL").
		Group("user_service_logs.service_id, services.name").
		Scan(&serviceRows).Error; err != nil {
		return stats, domain.ErrDataBaseInternalError
	}

	stats.ServiceStats = make([]domain.ServiceUsageStats, 0)
	for _, row := range serviceRows {
		totalSeconds := int(row.TotalNanoseconds / 1000000000) // Convert nanoseconds to seconds
		avgDuration := 0.0
		if row.TotalUsers > 0 {
			avgDuration = float64(totalSeconds) / float64(row.TotalUsers)
		}
		stats.ServiceStats = append(stats.ServiceStats, domain.ServiceUsageStats{
			ServiceID:    row.ServiceID,
			ServiceName:  row.ServiceName,
			TotalUsers:   row.TotalUsers,
			TotalSeconds: totalSeconds,
			AvgDuration:  avgDuration,
		})
	}

	// Get recent activity (last 10 entries)
	type RecentActivityRow struct {
		ID          uint
		UserID      uint
		UserEmail   string
		ServiceID   uint
		ServiceName string
		Duration    int64
		CreatedAt   string
	}
	var activityRows []RecentActivityRow
	if err := r.db.WithContext(ctx).
		Table("user_service_logs").
		Select(`
			user_service_logs.id,
			user_service_logs.user_id,
			users.email as user_email,
			user_service_logs.service_id,
			services.name as service_name,
			user_service_logs.duration / 1000000000 as duration,
			user_service_logs.created_at
		`).
		Joins("LEFT JOIN users ON users.id = user_service_logs.user_id").
		Joins("LEFT JOIN services ON services.id = user_service_logs.service_id").
		Where("user_service_logs.deleted_at IS NULL").
		Order("user_service_logs.created_at DESC").
		Limit(10).
		Scan(&activityRows).Error; err != nil {
		return stats, domain.ErrDataBaseInternalError
	}

	stats.RecentActivity = make([]domain.RecentActivityItem, 0)
	for _, row := range activityRows {
		stats.RecentActivity = append(stats.RecentActivity, domain.RecentActivityItem{
			ID:          row.ID,
			UserID:      row.UserID,
			UserEmail:   row.UserEmail,
			ServiceID:   row.ServiceID,
			ServiceName: row.ServiceName,
			Duration:    int(row.Duration),
			CreatedAt:   row.CreatedAt,
		})
	}

	return stats, nil
}
