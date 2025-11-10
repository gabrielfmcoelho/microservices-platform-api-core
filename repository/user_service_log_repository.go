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
// duration parameter is in seconds, but we store as nanoseconds (time.Duration)
func (r *userServiceLogRepository) UpdateDuration(ctx context.Context, userServiceLogID uint, duration int) error {
	// Convert seconds to nanoseconds (time.Duration stores nanoseconds)
	durationNanoseconds := int64(duration) * 1000000000
	if err := r.db.WithContext(ctx).Model(&domain.UserServiceLog{}).
		Where("id = ?", userServiceLogID).
		Update("duration", gorm.Expr("duration + ?", durationNanoseconds)).Error; err != nil {
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
func (r *userServiceLogRepository) GetUsageStatistics(ctx context.Context, organizationID *uint, startDate *string, endDate *string) (domain.UsageStatistics, error) {
	var stats domain.UsageStatistics

	// Build base query with optional filters
	baseQuery := r.db.WithContext(ctx).Model(&domain.UserServiceLog{})

	// Apply organization filter if provided
	if organizationID != nil {
		baseQuery = baseQuery.
			Joins("INNER JOIN users ON users.id = user_service_logs.user_id").
			Where("users.organization_id = ?", *organizationID)
	}

	// Apply date range filter if provided
	if startDate != nil && *startDate != "" {
		baseQuery = baseQuery.Where("user_service_logs.created_at >= ?", *startDate)
	}
	if endDate != nil && *endDate != "" {
		baseQuery = baseQuery.Where("user_service_logs.created_at <= ?", *endDate)
	}

	// Get total unique users (with activity)
	var totalUsers int64
	if err := baseQuery.
		Distinct("user_service_logs.user_id").
		Count(&totalUsers).Error; err != nil {
		return stats, domain.ErrDataBaseInternalError
	}
	stats.TotalUsers = int(totalUsers)

	// Get total organization users (if organization filter is applied)
	if organizationID != nil {
		var totalOrgUsers int64
		if err := r.db.WithContext(ctx).
			Model(&domain.User{}).
			Where("organization_id = ?", *organizationID).
			Count(&totalOrgUsers).Error; err != nil {
			return stats, domain.ErrDataBaseInternalError
		}
		stats.TotalOrgUsers = int(totalOrgUsers)
	} else {
		stats.TotalOrgUsers = 0
	}

	// Get total duration (sum of all durations in nanoseconds, convert to seconds)
	type TotalDurationResult struct {
		TotalNanoseconds int64
	}
	var durationResult TotalDurationResult

	// Rebuild query for duration with same filters
	durationQuery := r.db.WithContext(ctx).Model(&domain.UserServiceLog{})
	if organizationID != nil {
		durationQuery = durationQuery.
			Joins("INNER JOIN users ON users.id = user_service_logs.user_id").
			Where("users.organization_id = ?", *organizationID)
	}
	if startDate != nil && *startDate != "" {
		durationQuery = durationQuery.Where("user_service_logs.created_at >= ?", *startDate)
	}
	if endDate != nil && *endDate != "" {
		durationQuery = durationQuery.Where("user_service_logs.created_at <= ?", *endDate)
	}

	if err := durationQuery.
		Select("COALESCE(SUM(user_service_logs.duration), 0) as total_nanoseconds").
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

	// Build service stats query with filters
	serviceStatsQuery := r.db.WithContext(ctx).
		Table("user_service_logs").
		Select(`
			user_service_logs.service_id,
			services.name as service_name,
			COUNT(DISTINCT user_service_logs.user_id) as total_users,
			COALESCE(SUM(user_service_logs.duration), 0) as total_nanoseconds
		`).
		Joins("LEFT JOIN services ON services.id = user_service_logs.service_id").
		Where("user_service_logs.deleted_at IS NULL")

	if organizationID != nil {
		serviceStatsQuery = serviceStatsQuery.
			Joins("INNER JOIN users ON users.id = user_service_logs.user_id").
			Where("users.organization_id = ?", *organizationID)
	}
	if startDate != nil && *startDate != "" {
		serviceStatsQuery = serviceStatsQuery.Where("user_service_logs.created_at >= ?", *startDate)
	}
	if endDate != nil && *endDate != "" {
		serviceStatsQuery = serviceStatsQuery.Where("user_service_logs.created_at <= ?", *endDate)
	}

	if err := serviceStatsQuery.
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

	// Build recent activity query with filters
	activityQuery := r.db.WithContext(ctx).
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
		Where("user_service_logs.deleted_at IS NULL")

	if organizationID != nil {
		activityQuery = activityQuery.Where("users.organization_id = ?", *organizationID)
	}
	if startDate != nil && *startDate != "" {
		activityQuery = activityQuery.Where("user_service_logs.created_at >= ?", *startDate)
	}
	if endDate != nil && *endDate != "" {
		activityQuery = activityQuery.Where("user_service_logs.created_at <= ?", *endDate)
	}

	if err := activityQuery.
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

	// Get time-series data (usage by day and service)
	type TimeSeriesRow struct {
		Date         string
		ServiceName  string
		TotalSeconds int64
		AccessCount  int
	}
	var timeSeriesRows []TimeSeriesRow

	// Build time-series query with filters
	timeSeriesQuery := r.db.WithContext(ctx).
		Table("user_service_logs").
		Select(`
			DATE(user_service_logs.created_at) as date,
			services.name as service_name,
			FLOOR(COALESCE(SUM(user_service_logs.duration), 0) / 1000000000) as total_seconds,
			COUNT(*) as access_count
		`).
		Joins("LEFT JOIN services ON services.id = user_service_logs.service_id").
		Where("user_service_logs.deleted_at IS NULL")

	if organizationID != nil {
		timeSeriesQuery = timeSeriesQuery.
			Joins("INNER JOIN users ON users.id = user_service_logs.user_id").
			Where("users.organization_id = ?", *organizationID)
	}
	if startDate != nil && *startDate != "" {
		timeSeriesQuery = timeSeriesQuery.Where("user_service_logs.created_at >= ?", *startDate)
	}
	if endDate != nil && *endDate != "" {
		timeSeriesQuery = timeSeriesQuery.Where("user_service_logs.created_at <= ?", *endDate)
	}

	if err := timeSeriesQuery.
		Group("DATE(user_service_logs.created_at), services.name").
		Order("date ASC").
		Scan(&timeSeriesRows).Error; err != nil {
		return stats, domain.ErrDataBaseInternalError
	}

	// Transform rows into time-series data structure
	dateMap := make(map[string]map[string]domain.TimeSeriesService)
	for _, row := range timeSeriesRows {
		if _, exists := dateMap[row.Date]; !exists {
			dateMap[row.Date] = make(map[string]domain.TimeSeriesService)
		}
		dateMap[row.Date][row.ServiceName] = domain.TimeSeriesService{
			TotalSeconds: int(row.TotalSeconds),
			AccessCount:  row.AccessCount,
		}
	}

	// Convert map to array
	stats.TimeSeriesData = make([]domain.TimeSeriesDataPoint, 0)
	for date, services := range dateMap {
		stats.TimeSeriesData = append(stats.TimeSeriesData, domain.TimeSeriesDataPoint{
			Date:     date,
			Services: services,
		})
	}

	return stats, nil
}
