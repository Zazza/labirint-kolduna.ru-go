package auth

import (
	"gamebook-backend/modules/auth/controller"
	"gamebook-backend/modules/auth/service"
	"gamebook-backend/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/samber/do"

	"gamebook-backend/middlewares"
)

func RegisterRoutes(server *gin.Engine, injector *do.Injector) {
	authController := do.MustInvoke[controller.AuthController](injector)
	jwtService := do.MustInvokeNamed[service.JWTService](injector, constants.JWTService)

	authRoutes := server.Group("/api/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh", authController.RefreshToken)
		authRoutes.POST("/logout", middlewares.Authenticate(jwtService), authController.Logout)
	}
}
