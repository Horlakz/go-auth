package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/horlakz/go-auth/constants"
	"github.com/horlakz/go-auth/database"
	"github.com/horlakz/go-auth/handler"
	user_repository "github.com/horlakz/go-auth/repository/user"
	"github.com/horlakz/go-auth/service"
)

func InitializeAuthRouter(router fiber.Router, db database.DatabaseInterface, env constants.Env) {
	// Repositories
	userRepository := user_repository.NewUserRepository(db)
	roleRepository := user_repository.NewRoleRepository(db)
	userRoleRepository := user_repository.NewUserRoleRepository(db)
	verificationCodeRepository := user_repository.NewVerificationCodeRepository(db)

	// config
	smtpService := service.NewSmtpService(env)
	// Services
	authService := service.NewAuthService(
		userRepository,
		roleRepository,
		userRoleRepository,
		verificationCodeRepository,
		smtpService,
	)

	// Handler
	authHandler := handler.NewAuthHandler(authService)

	// Routers
	authRoute := router.Group("/auth")

	// Routes
	authRoute.Post("/login", authHandler.Login)
	authRoute.Post("/register", authHandler.Register)
	authRoute.Post("/refresh-token", authHandler.RefreshAccessToken)
	authRoute.Post("/forgot-password", authHandler.ForgotPassword)
	authRoute.Post("/reset-password", authHandler.ResetPassword)

}
