package repository

import (
	"context"
	"errors"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"gorm.io/gorm"
)

type contactIntentRepository struct {
	db *gorm.DB
}

func NewContactIntentRepository(db *gorm.DB) domain.ContactIntentRepository {
	return &contactIntentRepository{
		db: db,
	}
}

// Create creates a new contact intent in the database
func (r *contactIntentRepository) Create(ctx context.Context, contactIntent *domain.ContactIntent) error {
	if err := r.db.WithContext(ctx).Create(contactIntent).Error; err != nil {
		return domain.ErrDataBaseInternalError
	}
	return nil
}

// Fetch returns all contact intents from the database
func (r *contactIntentRepository) Fetch(ctx context.Context) ([]domain.ContactIntent, error) {
	var contactIntents []domain.ContactIntent
	if err := r.db.WithContext(ctx).Order("created_at DESC").Find(&contactIntents).Error; err != nil {
		return nil, domain.ErrDataBaseInternalError
	}
	return contactIntents, nil
}

// GetByID returns a specific contact intent by ID
func (r *contactIntentRepository) GetByID(ctx context.Context, id uint) (domain.ContactIntent, error) {
	var contactIntent domain.ContactIntent
	if err := r.db.WithContext(ctx).First(&contactIntent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return contactIntent, domain.ErrNotFound
		}
		return contactIntent, domain.ErrDataBaseInternalError
	}
	return contactIntent, nil
}

// UpdateStatus updates only the status field of a contact intent
func (r *contactIntentRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
	// Check if the contact intent exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update only the status field
	if err := r.db.WithContext(ctx).
		Model(&domain.ContactIntent{}).
		Where("id = ?", id).
		Update("status", status).Error; err != nil {
		return domain.ErrDataBaseInternalError
	}

	return nil
}
