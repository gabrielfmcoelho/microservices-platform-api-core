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

func NewServiceRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, group *gin.RouterGroup) {
	sr := repository.NewServiceRepository(db)
	uslr := repository.NewUserServiceLogRepository(db)
	ur := repository.NewUserRepository(db)
	sc := &controller.ServiceController{
		ServiceUsecase: usecase.NewServiceUsecase(sr, uslr, timeout),
		UserUsecase:    usecase.NewUserUsecase(ur, timeout),
		Env:            env,
	}

	group.POST("/service", sc.CreateService)
	group.GET("/services", sc.FetchServices)
	group.GET("/service/info", sc.GetServiceByIdentifier)
	group.GET("/service/:serviceID/application", sc.GetServiceApplication) // "start" usage
	group.POST("/service/:serviceID/organization/:organizationID", sc.SetServiceAvailabilityToOrganization)
	group.GET("/services/organization", sc.GetServicesByOrganization)
	group.PUT("/service/:serviceID", sc.UpdateService)
	group.DELETE("/service/:serviceID", sc.DeleteService)
	group.PATCH("/service/heartbeat", sc.HeartbeatService) // "update" usage duration
}
