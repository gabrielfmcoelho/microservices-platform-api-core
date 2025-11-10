package controller

import (
	"net/http"

	"github.com/gabrielfmcoelho/platform-core/bootstrap"
	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	OrganizationUsecase domain.OrganizationUsecase
	Env                 *bootstrap.Env
}

// @Summary Create a new organization
// @Description Create a new organization with the input payload
// @Tags Organization
// @ID createOrganization
// @Accept json
// @Produce json
// @Param organization body domain.CreateOrganization true "Organization object"
// @Success 201 {object} domain.SuccessResponse{data=domain.PublicOrganization}
// @Failure 400 {object} domain.ErrorResponse "Bad Request - Invalid input"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /organization [post]
func (oc *OrganizationController) CreateOrganization(c *gin.Context) {
	var createOrg domain.CreateOrganization
	if err := c.ShouldBindJSON(&createOrg); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid input: " + err.Error(),
		})
		return
	}

	// Create Organization object from CreateOrganization
	organization := domain.Organization{
		Name:   createOrg.Name,
		RoleID: createOrg.OrganizationRoleID,
	}

	err := oc.OrganizationUsecase.Create(c, &organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to create organization: " + err.Error(),
		})
		return
	}

	// Return the created organization
	publicOrg := domain.PublicOrganization{
		ID:      organization.ID,
		Name:    organization.Name,
		LogoUrl: organization.LogoUrl,
	}

	c.JSON(http.StatusCreated, parser.ToSuccessResponse(publicOrg))
}

// @Summary Get all organizations
// @Description Get all organizations from the database
// @Tags Organization
// @ID fetchOrganizations
// @Produce json
// @Success 200 {object} domain.SuccessResponse{data=[]domain.PublicOrganization} "List of organizations"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /organizations [get]
func (oc *OrganizationController) FetchOrganizations(c *gin.Context) {
	organizations, err := oc.OrganizationUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch organizations: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(organizations))
}

// @Summary Get organization by ID or name
// @Description Get organization by ID or name
// @Tags Organization
// @ID getOrganization
// @Produce json
// @Param identifier path string true "Organization ID or name"
// @Success 200 {object} domain.SuccessResponse{data=domain.PublicOrganization} "Organization object"
// @Failure 404 {object} domain.ErrorResponse "Not Found - Organization not found"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /organization/{identifier} [get]
func (oc *OrganizationController) GetOrganization(c *gin.Context) {
	identifier := c.Param("identifier")

	organization, err := oc.OrganizationUsecase.GetByIdentifier(c, identifier)
	if err != nil {
		switch err {
		case domain.ErrNotFound:
			c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Organization not found"})
		default:
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(organization))
}

// @Summary Delete organization
// @Description Delete an organization by ID
// @Tags Organization
// @ID deleteOrganization
// @Param id path int true "Organization ID"
// @Success 200 {object} domain.SuccessResponse
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /organization/{id} [delete]
func (oc *OrganizationController) DeleteOrganization(c *gin.Context) {
	idParam := c.Param("id")
	id, err := internal.ParseUint(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "Invalid organization ID"})
		return
	}

	err = oc.OrganizationUsecase.Delete(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to delete organization: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(gin.H{"message": "Organization deleted successfully"}))
}
