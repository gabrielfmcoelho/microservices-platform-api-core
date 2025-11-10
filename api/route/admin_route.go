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

func NewAdminRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	// Initialize repositories
	userServiceLogRepo := repository.NewUserServiceLogRepository(db)
	contactIntentRepo := repository.NewContactIntentRepository(db)
	organizationRoleRepo := repository.NewOrganizationRoleRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)
	userRepo := repository.NewUserRepository(db)
	serviceRepo := repository.NewServiceRepository(db)

	// Initialize admin controller
	ac := &controller.AdminController{
		UserServiceLogUsecase:      usecase.NewUserServiceLogUsecase(userServiceLogRepo, timeout),
		ContactIntentUsecase:       usecase.NewContactIntentUsecase(contactIntentRepo, timeout),
		OrganizationRoleRepository: organizationRoleRepo,
		UserRoleRepository:         userRoleRepo,
		UserUsecase:                usecase.NewUserUsecase(userRepo, timeout),
		ServiceUsecase:             usecase.NewServiceUsecase(serviceRepo, userServiceLogRepo, timeout),
		Env:                        env,
	}

	// Register admin routes (all protected)
	group.GET("/admin/statistics", ac.GetUsageStatistics)
	group.GET("/admin/contact-intents", ac.GetContactIntents)
	group.GET("/admin/organization-roles", ac.GetOrganizationRoles)
	group.GET("/admin/user-roles", ac.GetUserRoles)
	group.GET("/admin/organizations/:id/services", ac.GetOrganizationServices)
	group.GET("/admin/organizations/:id/users", ac.GetOrganizationUsers)
	group.DELETE("/admin/organizations/:organizationId/services/:serviceId", ac.UnlinkServiceFromOrganization)
	group.PATCH("/admin/users/:userId/:action", ac.ToggleUserArchiveStatus)
}
