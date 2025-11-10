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

func NewContactIntentRouter(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, publicGroup *gin.RouterGroup, protectedGroup *gin.RouterGroup) {
	cir := repository.NewContactIntentRepository(db)
	cic := &controller.ContactIntentController{
		ContactIntentUsecase: usecase.NewContactIntentUsecase(cir, timeout),
		Env:                  env,
	}

	// Public route - anyone can submit a contact intent
	publicGroup.POST("/contact-intent", cic.CreateContactIntent)

	// Protected routes - admin only
	protectedGroup.GET("/contact-intents", cic.FetchContactIntents)
	protectedGroup.PATCH("/contact-intent/:id/status", cic.UpdateContactIntentStatus)
}
