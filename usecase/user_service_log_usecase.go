package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
)

type userServiceLogUsecase struct {
	userServiceLogRepo domain.UserServiceLogRepository
	contextTimeout     time.Duration
}

func NewUserServiceLogUsecase(
	repo domain.UserServiceLogRepository,
	timeout time.Duration,
) domain.UserServiceLogUsecase {
	return &userServiceLogUsecase{
		userServiceLogRepo: repo,
		contextTimeout:     timeout,
	}
}

// Fetch all UserServiceLog entries
func (u *userServiceLogUsecase) Fetch(ctx context.Context) ([]domain.PublicUserServiceLog, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	logs, err := u.userServiceLogRepo.Fetch(c)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return nil, domain.ErrDataBaseInternalError
		}
		return nil, domain.ErrInternalServerError
	}

	PublicLogs := make([]domain.PublicUserServiceLog, 0, len(logs))
	for _, log := range logs {
		PublicLogs = append(PublicLogs, parser.ToPublicUserServiceLog(log))
	}

	return PublicLogs, nil
}

// GetByIdentifier tries to parse identifier to either:
// - a numeric ID -> GetByID
// - or if it starts with "user:" -> parse the rest as userID
// - or if it starts with "service:" -> parse the rest as serviceID
// Otherwise returns ErrInvalidIdentifier
func (u *userServiceLogUsecase) GetByIdentifier(ctx context.Context, identifier string) (domain.PublicUserServiceLog, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	var publicLog domain.PublicUserServiceLog
	var log domain.UserServiceLog
	var err error

	// Example rules:
	// 1) if it's numeric: interpret as ID
	// 2) if starts with "user:" => parse userID
	// 3) if starts with "service:" => parse serviceID
	// else => err

	switch {
	case internal.IsNumeric(identifier):
		id, convErr := internal.ParseUint(identifier)
		if convErr != nil {
			return publicLog, domain.ErrInvalidIdentifier
		}
		log, err = u.userServiceLogRepo.GetByID(c, id)

	case len(identifier) > 5 && identifier[:5] == "user:":
		userIDStr := identifier[5:]
		id, convErr := internal.ParseUint(userIDStr)
		if convErr != nil {
			return publicLog, domain.ErrInvalidIdentifier
		}
		log, err = u.userServiceLogRepo.GetByUserID(c, id)

	case len(identifier) > 8 && identifier[:8] == "service:":
		serviceIDStr := identifier[8:]
		id, convErr := internal.ParseUint(serviceIDStr)
		if convErr != nil {
			return publicLog, domain.ErrInvalidIdentifier
		}
		log, err = u.userServiceLogRepo.GetByServiceID(c, id)

	default:
		// not numeric, not user/service pattern
		return publicLog, domain.ErrInvalidIdentifier
	}

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return publicLog, domain.ErrNotFound
		}
		return publicLog, domain.ErrInternalServerError
	}

	return parser.ToPublicUserServiceLog(log), nil
}

// Delete removes a UserServiceLog by ID
func (u *userServiceLogUsecase) Delete(ctx context.Context, userServiceLogID uint) error {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	if err := u.userServiceLogRepo.Delete(c, userServiceLogID); err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}
	return nil
}

// GetUsageStatistics returns aggregated usage statistics for admin dashboard
func (u *userServiceLogUsecase) GetUsageStatistics(ctx context.Context, organizationID *uint, startDate *string, endDate *string) (domain.UsageStatistics, error) {
	c, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	stats, err := u.userServiceLogRepo.GetUsageStatistics(c, organizationID, startDate, endDate)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.UsageStatistics{}, domain.ErrDataBaseInternalError
		}
		return domain.UsageStatistics{}, domain.ErrInternalServerError
	}

	return stats, nil
}
