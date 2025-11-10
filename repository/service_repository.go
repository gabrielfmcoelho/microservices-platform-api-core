package repository

import (
	"context"
	"errors"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) domain.ServiceRepository {
	return &serviceRepository{
		db: db,
	}
}

// Create insere um novo service no banco de dados
func (r *serviceRepository) Create(ctx context.Context, service *domain.Service) error {
	if err := r.db.WithContext(ctx).Create(service).Error; err != nil {
		return domain.ErrDataBaseInternalError
	}
	return nil
}

// Fetch retorna todos os serviços cadastrados
func (r *serviceRepository) Fetch(ctx context.Context) ([]domain.Service, error) {
	var services []domain.Service
	if err := r.db.WithContext(ctx).Find(&services).Error; err != nil {
		return nil, domain.ErrDataBaseInternalError
	}
	return services, nil
}

// GetByID retorna um service específico com base no ID
func (r *serviceRepository) GetByID(ctx context.Context, id uint) (domain.Service, error) {
	var service domain.Service
	if err := r.db.WithContext(ctx).First(&service, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service, domain.ErrNotFound
		}
		return service, domain.ErrDataBaseInternalError
	}
	return service, nil
}

// GetByName retorna um service específico com base no nome
func (r *serviceRepository) GetByName(ctx context.Context, name string) (domain.Service, error) {
	var service domain.Service
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return service, domain.ErrNotFound
		}
		return service, domain.ErrDataBaseInternalError
	}
	return service, nil
}

// GetByOrganization retorna todos os serviços vinculados a uma organização
func (r *serviceRepository) GetByOrganization(ctx context.Context, organizationID uint) ([]domain.Service, error) {
	var services []domain.Service
	if err := r.db.WithContext(ctx).
		Preload("Organization").
		Joins("JOIN organization_services ON services.id = organization_services.service_id").
		Where("organization_services.organization_id = ?", organizationID).
		Find(&services).Error; err != nil {
		return nil, domain.ErrDataBaseInternalError
	}
	return services, nil
}

// GetMarketing retorna todos os serviços de marketing
func (r *serviceRepository) GetMarketing(ctx context.Context) ([]domain.Service, error) {
	var services []domain.Service
	if err := r.db.WithContext(ctx).
		Where("is_marketing = ?", true).
		Find(&services).Error; err != nil {
		return nil, domain.ErrDataBaseInternalError
	}
	return services, nil
}

// SetAvailabilityToOrganization vincula o service a uma organização na tabela pivô (many2many)
func (r *serviceRepository) SetAvailabilityToOrganization(ctx context.Context, serviceID uint, organizationID uint) error {
	// Para associar, precisamos obter primeiro o service e a organization
	var service domain.Service
	if err := r.db.WithContext(ctx).First(&service, serviceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return domain.ErrDataBaseInternalError
	}

	var organization domain.Organization
	if err := r.db.WithContext(ctx).First(&organization, organizationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return domain.ErrDataBaseInternalError
	}

	// GORM many2many association
	if err := r.db.WithContext(ctx).Model(&service).Association("Organization").Append(&organization); err != nil {
		return domain.ErrDataBaseInternalError
	}

	return nil
}

// RemoveAvailabilityFromOrganization remove o vínculo do service com uma organização
func (r *serviceRepository) RemoveAvailabilityFromOrganization(ctx context.Context, serviceID uint, organizationID uint) error {
	var service domain.Service
	if err := r.db.WithContext(ctx).First(&service, serviceID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return domain.ErrDataBaseInternalError
	}

	var organization domain.Organization
	if err := r.db.WithContext(ctx).First(&organization, organizationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return domain.ErrDataBaseInternalError
	}

	// GORM many2many association - Delete removes the association
	if err := r.db.WithContext(ctx).Model(&service).Association("Organization").Delete(&organization); err != nil {
		return domain.ErrDataBaseInternalError
	}

	return nil
}

// Update atualiza os dados de um service no banco
func (r *serviceRepository) Update(ctx context.Context, serviceID uint, serviceData *domain.Service) error {
	// A forma de atualização depende de como você deseja aplicar as mudanças.
	// Exemplo simples de updates:
	if err := r.db.WithContext(ctx).
		Model(&domain.Service{}).
		Where("id = ?", serviceID).
		Updates(serviceData).Error; err != nil {
		return domain.ErrDataBaseInternalError
	}
	return nil
}

// Delete remove um service do banco de dados
func (r *serviceRepository) Delete(ctx context.Context, serviceID uint) error {
	// Exemplo: deleção hard (exclui permanentemente)
	if err := r.db.WithContext(ctx).Delete(&domain.Service{}, serviceID).Error; err != nil {
		return domain.ErrDataBaseInternalError
	}
	return nil
}
