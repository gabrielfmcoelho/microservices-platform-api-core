package route

import (
	"time"

	"github.com/gabrielfmcoelho/platform-core/api/controller"
	"github.com/gabrielfmcoelho/platform-core/bootstrap"
	"github.com/gabrielfmcoelho/platform-core/repository"
	"github.com/gabrielfmcoelho/platform-core/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewOrganizationRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	or := repository.NewOrganizationRepository(db)
	oc := &controller.OrganizationController{
		OrganizationUsecase: usecase.NewOrganizationUsecase(or),
		Env:                 env,
	}

	group.POST("/organization", oc.CreateOrganization)            // Create a new organization
	group.GET("/organizations", oc.FetchOrganizations)            // Get all organizations
	group.GET("/organization/:identifier", oc.GetOrganization)    // Get organization by ID or name
	group.DELETE("/organization/:id", oc.DeleteOrganization)      // Delete organization
}
