package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielfmcoelho/platform-core/domain"
)

type ContactIntentUsecase struct {
	contactIntentRepository domain.ContactIntentRepository
	contextTimeout          time.Duration
}

func NewContactIntentUsecase(contactIntentRepository domain.ContactIntentRepository, timeout time.Duration) *ContactIntentUsecase {
	return &ContactIntentUsecase{
		contactIntentRepository: contactIntentRepository,
		contextTimeout:          timeout,
	}
}

func (ciu *ContactIntentUsecase) Create(c context.Context, createContactIntent *domain.CreateContactIntent) error {
	ctx, cancel := context.WithTimeout(c, ciu.contextTimeout)
	defer cancel()

	// Create the contact intent entity
	contactIntent := &domain.ContactIntent{
		Name:        createContactIntent.Name,
		Email:       createContactIntent.Email,
		Phone:       createContactIntent.Phone,
		Company:     createContactIntent.Company,
		Message:     createContactIntent.Message,
		ServiceName: createContactIntent.ServiceName,
		Status:      "pending", // Default status
	}

	err := ciu.contactIntentRepository.Create(ctx, contactIntent)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}

func (ciu *ContactIntentUsecase) Fetch(c context.Context) ([]domain.PublicContactIntent, error) {
	ctx, cancel := context.WithTimeout(c, ciu.contextTimeout)
	defer cancel()

	contactIntents, err := ciu.contactIntentRepository.Fetch(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return nil, domain.ErrDataBaseInternalError
		}
		return nil, domain.ErrInternalServerError
	}

	// Parse to PublicContactIntent
	publicContactIntents := make([]domain.PublicContactIntent, 0)
	for _, contactIntent := range contactIntents {
		publicContactIntents = append(publicContactIntents, domain.PublicContactIntent{
			ID:          contactIntent.ID,
			Name:        contactIntent.Name,
			Email:       contactIntent.Email,
			Phone:       contactIntent.Phone,
			Company:     contactIntent.Company,
			Message:     contactIntent.Message,
			ServiceName: contactIntent.ServiceName,
			Status:      contactIntent.Status,
			CreatedAt:   contactIntent.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return publicContactIntents, nil
}

func (ciu *ContactIntentUsecase) UpdateStatus(c context.Context, id uint, status string) error {
	ctx, cancel := context.WithTimeout(c, ciu.contextTimeout)
	defer cancel()

	// Validate status value
	validStatuses := map[string]bool{
		"pending":   true,
		"contacted": true,
		"completed": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return domain.ErrBadRequest
	}

	err := ciu.contactIntentRepository.UpdateStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}
