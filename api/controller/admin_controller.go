package controller

import (
	"net/http"
	"strconv"

	"github.com/gabrielfmcoelho/platform-core/bootstrap"
	"github.com/gabrielfmcoelho/platform-core/domain"
	"github.com/gabrielfmcoelho/platform-core/internal/parser"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	UserServiceLogUsecase      domain.UserServiceLogUsecase
	ContactIntentUsecase       domain.ContactIntentUsecase
	OrganizationRoleRepository domain.OrganizationRoleRepository
	UserRoleRepository         domain.UserRoleRepository
	UserUsecase                domain.UserUsecase
	ServiceUsecase             domain.ServiceUsecase
	Env                        *bootstrap.Env
}

// @Summary Get usage statistics
// @Description Get dashboard statistics for admin with optional filters
// @Tags Admin
// @ID getUsageStatistics
// @Security BearerAuth
// @Produce json
// @Param organization_id query int false "Organization ID filter"
// @Param start_date query string false "Start date filter (ISO format)"
// @Param end_date query string false "End date filter (ISO format)"
// @Success 200 {object} domain.SuccessResponse{data=domain.UsageStatistics} "Usage statistics"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/statistics [get]
func (ac *AdminController) GetUsageStatistics(c *gin.Context) {
	// Parse optional query parameters
	var organizationID *uint
	if orgIDStr := c.Query("organization_id"); orgIDStr != "" {
		orgID, err := strconv.Atoi(orgIDStr)
		if err == nil {
			orgIDUint := uint(orgID)
			organizationID = &orgIDUint
		}
	}

	var startDate *string
	if start := c.Query("start_date"); start != "" {
		startDate = &start
	}

	var endDate *string
	if end := c.Query("end_date"); end != "" {
		endDate = &end
	}

	statistics, err := ac.UserServiceLogUsecase.GetUsageStatistics(c, organizationID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch usage statistics: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(statistics))
}

// @Summary Get all contact intents
// @Description Get all contact intents for admin review
// @Tags Admin
// @ID getContactIntents
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.SuccessResponse{data=[]domain.PublicContactIntent} "List of contact intents"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/contact-intents [get]
func (ac *AdminController) GetContactIntents(c *gin.Context) {
	contactIntents, err := ac.ContactIntentUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch contact intents: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(contactIntents))
}

// @Summary Get all organization roles
// @Description Get all organization roles for dropdowns
// @Tags Admin
// @ID getOrganizationRoles
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.SuccessResponse{data=[]domain.PublicOrganizationRole} "List of organization roles"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/organization-roles [get]
func (ac *AdminController) GetOrganizationRoles(c *gin.Context) {
	roles, err := ac.OrganizationRoleRepository.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch organization roles: " + err.Error(),
		})
		return
	}

	publicRoles := parser.ToPublicOrganizationRoles(roles)
	c.JSON(http.StatusOK, parser.ToSuccessResponse(publicRoles))
}

// @Summary Get all user roles
// @Description Get all user roles for dropdowns
// @Tags Admin
// @ID getUserRoles
// @Security BearerAuth
// @Produce json
// @Success 200 {object} domain.SuccessResponse{data=[]domain.PublicUserRole} "List of user roles"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/user-roles [get]
func (ac *AdminController) GetUserRoles(c *gin.Context) {
	roles, err := ac.UserRoleRepository.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch user roles: " + err.Error(),
		})
		return
	}

	publicRoles := parser.ToPublicUserRoles(roles)
	c.JSON(http.StatusOK, parser.ToSuccessResponse(publicRoles))
}

// @Summary Get organization services
// @Description Get all services linked to an organization
// @Tags Admin
// @ID getOrganizationServices
// @Security BearerAuth
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} domain.SuccessResponse{data=[]domain.HubService} "List of organization services"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/organizations/{id}/services [get]
func (ac *AdminController) GetOrganizationServices(c *gin.Context) {
	orgID := c.Param("id")
	organizationID, err := strconv.Atoi(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid organization ID",
		})
		return
	}

	services, err := ac.ServiceUsecase.GetByOrganization(c, uint(organizationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch organization services: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(services))
}

// @Summary Get organization users
// @Description Get all users belonging to an organization
// @Tags Admin
// @ID getOrganizationUsers
// @Security BearerAuth
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} domain.SuccessResponse{data=[]domain.PublicUser} "List of organization users"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/organizations/{id}/users [get]
func (ac *AdminController) GetOrganizationUsers(c *gin.Context) {
	orgID := c.Param("id")
	organizationID, err := strconv.Atoi(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid organization ID",
		})
		return
	}

	users, err := ac.UserUsecase.Fetch(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to fetch users: " + err.Error(),
		})
		return
	}

	// Filter users by organization ID
	var orgUsers []domain.PublicUser
	for _, user := range users {
		if user.OrganizationID == uint(organizationID) {
			orgUsers = append(orgUsers, user)
		}
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(orgUsers))
}

// @Summary Unlink service from organization
// @Description Removes the link between a service and an organization
// @Tags Admin
// @ID unlinkServiceFromOrganization
// @Security BearerAuth
// @Produce json
// @Param organizationId path int true "Organization ID"
// @Param serviceId path int true "Service ID"
// @Success 200 {object} domain.SuccessResponse "Service unlinked successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/organizations/{organizationId}/services/{serviceId} [delete]
func (ac *AdminController) UnlinkServiceFromOrganization(c *gin.Context) {
	orgID := c.Param("organizationId")
	srvID := c.Param("serviceId")

	organizationID, err := strconv.Atoi(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid organization ID",
		})
		return
	}

	serviceID, err := strconv.Atoi(srvID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid service ID",
		})
		return
	}

	err = ac.ServiceUsecase.RemoveAvailabilityFromOrganization(c, uint(serviceID), uint(organizationID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to unlink service: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(map[string]string{
		"message": "Service unlinked successfully",
	}))
}

// @Summary Toggle user archive status
// @Description Archives or unarchives a user (soft delete)
// @Tags Admin
// @ID toggleUserArchiveStatus
// @Security BearerAuth
// @Produce json
// @Param userId path int true "User ID"
// @Param action path string true "Action: archive or unarchive"
// @Success 200 {object} domain.SuccessResponse "User status toggled successfully"
// @Failure 400 {object} domain.ErrorResponse "Bad Request"
// @Failure 500 {object} domain.ErrorResponse "Internal Server Error"
// @Router /admin/users/{userId}/{action} [patch]
func (ac *AdminController) ToggleUserArchiveStatus(c *gin.Context) {
	userIDStr := c.Param("userId")
	action := c.Param("action")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid user ID",
		})
		return
	}

	if action != "archive" && action != "unarchive" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "Invalid action. Must be 'archive' or 'unarchive'",
		})
		return
	}

	if action == "archive" {
		err = ac.UserUsecase.Archive(c, uint(userID))
	} else {
		err = ac.UserUsecase.Unarchive(c, uint(userID))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
			Message: "Failed to " + action + " user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, parser.ToSuccessResponse(map[string]string{
		"message": "User " + action + "d successfully",
	}))
}
