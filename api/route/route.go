package route

import (
	"time"

	"github.com/gabrielfmcoelho/platform-core/api/middleware"
	"github.com/gabrielfmcoelho/platform-core/bootstrap"

	//_ "github.com/gabrielfmcoelho/platform-coredocs"
	"github.com/gin-gonic/gin"
	"github.com/mvrilo/go-redoc"
	ginredoc "github.com/mvrilo/go-redoc/gin"
	_ "github.com/swaggo/swag"
	"gorm.io/gorm"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db *gorm.DB, router *gin.Engine) {
	// Router documentation binding
	doc := redoc.Redoc{
		Title:       "Platform Core API",
		Description: "Platform Core API Description",
		SpecFile:    "./docs/swagger.json",
		SpecPath:    "/docs/swagger.json",
		DocsPath:    "/docs/",
	}
	router.GET("/docs/*any", ginredoc.New(doc))

	// All Public APIs
	publicRouter := router.Group("/")
	//NewSignupRouter(env, timeout, db, publicRouter)
	NewAuthRouter(env, timeout, db, publicRouter)
	NewPublicWebsiteRouter(env, timeout, db, publicRouter)
	//NewRefreshTokenRouter(env, timeout, db, publicRouter)

	// All Private APIs
	protectedRouter := router.Group("/")
	/// Middleware to verify AccessToken
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	NewUserRouter(env, timeout, db, protectedRouter)
	NewOrganizationRouter(env, timeout, db, protectedRouter)
	NewServiceRouter(env, timeout, db, protectedRouter)
	//NewProfileRouter(env, timeout, db, protectedRouter)
	//NewTaskRouter(env, timeout, db, protectedRouter)

	// Contact Intent Routes (both public and protected)
	NewContactIntentRouter(env, timeout, db, publicRouter, protectedRouter)

	// Admin Routes (all protected)
	NewAdminRouter(env, timeout, db, protectedRouter)
}
