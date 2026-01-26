package config

import (
	"backend/internal/delivery/http"
	"backend/internal/delivery/http/middleware"
	"backend/internal/delivery/http/route"
	"backend/internal/repository"
	"backend/internal/usecase"
	"backend/internal/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB        *gorm.DB
	App       *fiber.App
	Redis     *redis.Client
	Log       *logrus.Logger
	Validator *validator.Validate
	Config    *viper.Viper
}

// Bootstrap initializes and configures the application components.
func Bootstrap(config *BootstrapConfig) {
	// Repositories
	userRepository := repository.NewUserRepository(config.Log)
	walletRepository := repository.NewWalletRepository(config.Log)
	transactionRepository := repository.NewTransactionRepository(config.Log)
	walletMutationRepository := repository.NewWalletMutationRepository(config.Log)

	// Utilities
	tokenUtil := util.NewTokenUtil(config.Config.GetString("JWT_SECRET"), config.Redis)

	// Use Cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validator, userRepository, walletRepository, tokenUtil)
	walletUseCase := usecase.NewWalletUseCase(config.DB, config.Log, config.Validator, walletRepository)
	transactionUseCase := usecase.NewTransactionUseCase(config.DB, config.Log, config.Validator, transactionRepository, walletRepository, walletMutationRepository)
	walletMutationUseCase := usecase.NewWalletMutationUseCase(config.DB, config.Log, config.Validator, walletMutationRepository, walletRepository)

	// Controllers
	userController := http.NewUserController(config.Log, userUseCase)
	walletController := http.NewWalletController(config.Log, walletUseCase)
	transactionController := http.NewTransactionController(config.Log, transactionUseCase)
	walletMutationController := http.NewWalletMutationController(config.Log, walletMutationUseCase)

	// Middleware
	authMiddleware := middleware.NewAuth(userUseCase, tokenUtil)

	routeConfig := route.ConfigRoute{
		App:                      config.App,
		UserController:           userController,
		WalletController:         walletController,
		TransactionController:    transactionController,
		WalletMutationController: walletMutationController,
		AuthMiddleware:           authMiddleware,
	}

	routeConfig.Setup()
}
