package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
)

type serviceUsecase struct {
	serviceRepository        domain.ServiceRepository
	userServiceLogRepository domain.UserServiceLogRepository
	contextTimeout           time.Duration
}

// NewServiceUsecase cria um novo caso de uso para Service
func NewServiceUsecase(serviceRepository domain.ServiceRepository, userServiceLogRepository domain.UserServiceLogRepository, timeout time.Duration) domain.ServiceUsecase {
	return &serviceUsecase{
		serviceRepository:        serviceRepository,
		userServiceLogRepository: userServiceLogRepository,
		contextTimeout:           timeout,
	}
}

// Create cria um novo service
func (su *serviceUsecase) Create(ctx context.Context, service *domain.Service) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	err := su.serviceRepository.Create(ctx, service)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}

// Fetch retorna todos os serviços, convertidos em PublicService
func (su *serviceUsecase) Fetch(ctx context.Context) ([]domain.PublicService, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	services, err := su.serviceRepository.Fetch(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return nil, domain.ErrDataBaseInternalError
		}
		return nil, domain.ErrInternalServerError
	}

	publicServices := make([]domain.PublicService, 0, len(services))
	for _, s := range services {
		publicServices = append(publicServices, parser.ToPublicService(s))
	}
	return publicServices, nil
}

// GetByIdentifier obtém um serviço por ID (se o identifier for numérico) ou por nome (caso contrário)
func (su *serviceUsecase) GetByIdentifier(ctx context.Context, identifier string) (domain.PublicService, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	var service domain.Service
	var publicService domain.PublicService
	var err error

	if internal.IsNumeric(identifier) {
		// converter para uint
		id, convErr := internal.ParseUint(identifier)
		if convErr != nil {
			return publicService, domain.ErrInvalidIdentifier
		}
		service, err = su.serviceRepository.GetByID(ctx, id)
	} else {
		service, err = su.serviceRepository.GetByName(ctx, identifier)
	}

	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return publicService, domain.ErrNotFound
		}
		return publicService, domain.ErrInternalServerError
	}

	return parser.ToPublicService(service), nil
}

// GetByOrganization retorna todos os serviços vinculados a uma organização
func (su *serviceUsecase) GetByOrganization(ctx context.Context, organizationID uint) ([]domain.HubService, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	services, err := su.serviceRepository.GetByOrganization(ctx, organizationID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return nil, domain.ErrDataBaseInternalError
		}
		return nil, domain.ErrInternalServerError
	}

	hubServices := make([]domain.HubService, 0, len(services))
	for _, s := range services {
		hubServices = append(hubServices, parser.ToHubService(s))
	}
	return hubServices, nil
}

// GetMarketing retorna todos os serviços de marketing
func (su *serviceUsecase) GetMarketing(ctx context.Context) ([]domain.MarketingService, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	services, err := su.serviceRepository.GetMarketing(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrNotFound
		}
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return nil, domain.ErrDataBaseInternalError
		}
		return nil, domain.ErrInternalServerError
	}

	marketingServices := make([]domain.MarketingService, 0, len(services))
	for _, s := range services {
		marketingServices = append(marketingServices, parser.ToMarketingService(s))
	}
	return marketingServices, nil
}

// SetAvailabilityToOrganization vincula o service a uma organização
func (su *serviceUsecase) SetAvailabilityToOrganization(ctx context.Context, serviceID uint, organizationID uint) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	err := su.serviceRepository.SetAvailabilityToOrganization(ctx, serviceID, organizationID)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return domain.ErrInternalServerError
	}

	return nil
}

// RemoveAvailabilityFromOrganization remove o vínculo do service com uma organização
func (su *serviceUsecase) RemoveAvailabilityFromOrganization(ctx context.Context, serviceID uint, organizationID uint) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	err := su.serviceRepository.RemoveAvailabilityFromOrganization(ctx, serviceID, organizationID)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return domain.ErrInternalServerError
	}

	return nil
}

func (su *serviceUsecase) Use(ctx context.Context, userID uint, serviceID uint) (domain.UseService, uint, error) {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	var service domain.Service
	var useService domain.UseService
	var logID uint

	log := domain.UserServiceLog{
		UserID:    userID,
		ServiceID: serviceID,
	}

	err := su.userServiceLogRepository.Create(ctx, &log)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return useService, logID, domain.ErrDataBaseInternalError
		}
		return useService, logID, domain.ErrInternalServerError
	}

	logID = log.ID

	service, err = su.serviceRepository.GetByID(ctx, serviceID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return useService, logID, domain.ErrNotFound
		}
		return useService, logID, domain.ErrInternalServerError
	}

	return parser.ToUseService(service), logID, nil
}

func (su *serviceUsecase) Heartbeat(ctx context.Context, logID uint, duration int) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	err := su.userServiceLogRepository.UpdateDuration(ctx, logID, duration)
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

// Update atualiza os dados de um serviço
func (su *serviceUsecase) Update(ctx context.Context, serviceID uint, service *domain.Service) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	err := su.serviceRepository.Update(ctx, serviceID, service)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}

// Delete remove um serviço do banco
func (su *serviceUsecase) Delete(ctx context.Context, serviceID uint) error {
	ctx, cancel := context.WithTimeout(ctx, su.contextTimeout)
	defer cancel()

	err := su.serviceRepository.Delete(ctx, serviceID)
	if err != nil {
		if errors.Is(err, domain.ErrDataBaseInternalError) {
			return domain.ErrDataBaseInternalError
		}
		return domain.ErrInternalServerError
	}

	return nil
}
