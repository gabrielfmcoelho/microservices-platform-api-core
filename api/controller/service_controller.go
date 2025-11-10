package controller

import (
	"fmt"
	"net/http"

	"strings"

	"github.com/gabrielfmcoelho/platform-core/bootstrap"
	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
	"github.com/gabrielfmcoelho/platform-core/internal/tokenutil"
	"github.com/gin-gonic/gin"
)

type ServiceController struct {
	ServiceUsecase domain.ServiceUsecase
	UserUsecase    domain.UserUsecase
	Env            *bootstrap.Env
}

// CreateService cria um novo serviço
// @Summary Create Service
// @Description Creates a new service
// @Tags Service
// @Accept json
// @Produce json
// @Param service body domain.Service true "Service data"
// @Success 201 {object} domain.SuccessResponse{data=domain.PublicService}
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /services [post]
func (sc *ServiceController) CreateService(c *gin.Context) {
	var service domain.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := sc.ServiceUsecase.Create(c, &service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	// Retornar o service criado (ou alguma versão pública dele)
	c.JSON(http.StatusCreated, parser.ToSuccessResponse(parser.ToPublicService(service)))
}

// FetchServices retorna todos os serviços
// @Summary Fetch Services
// @Description Gets all available services
// @Tags Service
// @Produce json
// @Success 200 {array} domain.SuccessResponse{data=[]domain.PublicService}
// @Failure 500 {object} domain.ErrorResponse
// @Router /services [get]
func (sc *ServiceController) FetchServices(c *gin.Context) {
	services, err := sc.ServiceUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, parser.ToSuccessResponse(services))
}

// GetServiceByIdentifier retorna um serviço por ID ou nome
// @Summary Get Service by Identifier
// @Description Gets service by numeric ID (e.g., /services/123) or name (/services/my-service)
// @Tags Service
// @Produce json
// @Param identifier path string true "Service ID or Name"
// @Success 200 {object} domain.SuccessResponse{data=domain.PublicService}
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /services/{identifier} [get]
func (sc *ServiceController) GetServiceByIdentifier(c *gin.Context) {
	identifier := c.Param("identifier")

	service, err := sc.ServiceUsecase.GetByIdentifier(c, identifier)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, parser.ToSuccessResponse(service))
}

// GetServicesByOrganization retorna os serviços de uma organização
// @Summary Get Services by Organization
// @Description Gets all services linked to an organization
// @Tags Service
// @Produce json
// @Param organizationID path int true "Organization ID"
// @Success 200 {array} []domain.HubService
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /services/organization [get]
func (sc *ServiceController) GetServicesByOrganization(c *gin.Context) {
	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("x-user-id")
	if !exists {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "User ID not found in context"})
		return
	}

	// Get the user to access their organization_id
	userIDStr := fmt.Sprintf("%d", userID)
	user, err := sc.UserUsecase.GetByIdentifier(c, userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "Failed to get user: " + err.Error()})
		return
	}

	fmt.Printf("=== BACKEND: GetServicesByOrganization ===\n")
	fmt.Printf("User ID: %d, Organization ID: %d\n", user.ID, user.OrganizationID)

	services, err := sc.ServiceUsecase.GetByOrganization(c, user.OrganizationID)
	if err != nil {
		fmt.Printf("Error getting services: %v\n", err)
		switch err {
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		}
		return
	}

	fmt.Printf("Number of services found: %d\n", len(services))
	for i, svc := range services {
		fmt.Printf("Service %d: ID=%d, Name=%s, Status=%s\n", i+1, svc.ID, svc.Name, svc.Status)
	}

	c.JSON(http.StatusOK, services)
}

// SetServiceAvailabilityToOrganization vincula um service a uma organização
// @Summary Set Service Availability
// @Description Links the service to an organization
// @Tags Service
// @Produce json
// @Param serviceID path int true "Service ID"
// @Param organizationID path int true "Organization ID"
// @Success 200 {object} domain.SuccessResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /services/{serviceID}/organization/{organizationID} [post]
func (sc *ServiceController) SetServiceAvailabilityToOrganization(c *gin.Context) {
	serviceID := c.Param("serviceID")
	organizationID := c.Param("organizationID")

	// converter para uint
	// (poderíamos extrair essa lógica para uma função utilitária)
	var sID, oID uint
	_, err := fmt.Sscanf(serviceID, "%d", &sID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid serviceID"})
		return
	}
	_, err = fmt.Sscanf(organizationID, "%d", &oID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid organizationID"})
		return
	}

	err = sc.ServiceUsecase.SetAvailabilityToOrganization(c, sID, oID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Service availability set successfully."})
}

// GetServiceApplication
// @Summary Start using a service (create a usage log)
// @Description Logs that a user started using a service, returns log ID and public service data
// @Tags Service
// @Accept json
// @Produce json
// @Param serviceID path int true "Service ID"
// @Success 200 {object} domain.UseService
// @Failure 400 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /services/{serviceID}/application [get]
func (sc *ServiceController) GetServiceApplication(c *gin.Context) {
	// 1) Parse serviceID from path
	serviceIDParam := c.Param("serviceID")
	var sID uint
	if _, err := fmt.Sscanf(serviceIDParam, "%d", &sID); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid serviceID"})
		return
	}

	// 2) extract UserID from jwt token
	authHeader := c.Request.Header.Get("Authorization")
	t := strings.Split(authHeader, " ")
	authToken := t[1]
	userID, err := tokenutil.ExtractIDFromToken(authToken, sc.Env.AccessTokenSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
	}

	// 3) Call usecase
	service, logID, err := sc.ServiceUsecase.Use(c, uint(userID), sID)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		}
		return
	}

	service.LogID = logID

	// 4) Return result
	c.JSON(http.StatusOK, service)
}

// HeartbeatService
// @Summary Heartbeat usage
// @Description Adds usage duration (in seconds) to a log record
// @Tags Service
// @Accept json
// @Produce json
// @Param heartbeat body domain.Heartbeat true "Heartbeat data"
// @Success 200 {object} domain.SuccessResponse
// @Failure 404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /services/heartbeat [patch]
func (sc *ServiceController) HeartbeatService(c *gin.Context) {
	// 1) Parse JSON body for logID and duration
	var req domain.Heartbeat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	err := sc.ServiceUsecase.Heartbeat(c, req.LogID, req.Duration)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		}
		return
	}

	// 3) Return success
	c.JSON(http.StatusOK, domain.SuccessResponse{Message: "Usage duration updated successfully"})
}

// UpdateService atualiza um service
// @Summary Update Service
// @Description Updates service data
// @Tags Service
// @Accept json
// @Produce json
// @Param serviceID path int true "Service ID"
// @Param service body domain.Service true "Service data"
// @Success 200 {object} domain.SuccessResponse{data=domain.PublicService}
// @Failure 500 {object} domain.ErrorResponse
// @Router /services/{serviceID} [put]
func (sc *ServiceController) UpdateService(c *gin.Context) {
	serviceID := c.Param("serviceID")

	var sID uint
	_, err := fmt.Sscanf(serviceID, "%d", &sID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid serviceID"})
		return
	}

	var service domain.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// Log what was received
	fmt.Printf("=== BACKEND: Received service update for ID %d ===\n", sID)
	fmt.Printf("MarketingName: '%s'\n", service.MarketingName)
	fmt.Printf("Name: '%s'\n", service.Name)
	fmt.Printf("Description: '%s'\n", service.Description)
	fmt.Printf("AppUrl: '%s'\n", service.AppUrl)
	fmt.Printf("IconUrl: '%s'\n", service.IconUrl)
	fmt.Printf("ScreenshotUrl: '%s'\n", service.ScreenshotUrl)
	fmt.Printf("TagLine: '%s'\n", service.TagLine)
	fmt.Printf("Benefits: '%s'\n", service.Benefits)
	fmt.Printf("Features: '%s'\n", service.Features)
	fmt.Printf("Tags: '%s'\n", service.Tags)
	fmt.Printf("LastUpdate: '%s'\n", service.LastUpdate)
	fmt.Printf("Status: '%s'\n", service.Status)
	fmt.Printf("Price: %f\n", service.Price)
	fmt.Printf("Version: '%s'\n", service.Version)
	fmt.Printf("IsMarketing: %t\n", service.IsMarketing)

	err = sc.ServiceUsecase.Update(c, sID, &service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(parser.ToPublicService(service)))
}

// DeleteService deleta um service
// @Summary Delete Service
// @Description Deletes a service
// @Tags Service
// @Produce json
// @Param serviceID path int true "Service ID"
// @Success 204 "No Content"
// @Failure 500 {object} domain.ErrorResponse
// @Router /services/{serviceID} [delete]
func (sc *ServiceController) DeleteService(c *gin.Context) {
	serviceID := c.Param("serviceID")

	var sID uint
	_, err := fmt.Sscanf(serviceID, "%d", &sID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid serviceID"})
		return
	}

	err = sc.ServiceUsecase.Delete(c, sID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	// Retornamos 204 No Content pois não há conteúdo no response
	c.Status(http.StatusNoContent)
}
