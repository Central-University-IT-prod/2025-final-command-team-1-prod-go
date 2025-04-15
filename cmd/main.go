package main

import (
	"context"
	"os"
	"runtime/debug"
	"time"

	"example.com/m/internal/api/v1/adapters/controllers"
	"example.com/m/internal/api/v1/adapters/repositories"
	"example.com/m/internal/api/v1/core/application/services/auth_service"
	"example.com/m/internal/api/v1/core/application/services/booking_service"
	"example.com/m/internal/api/v1/core/application/services/gpt_service"
	"example.com/m/internal/api/v1/core/application/services/place_service"
	"example.com/m/internal/api/v1/core/application/services/post_service"
	"example.com/m/internal/api/v1/core/application/services/review_service"
	"example.com/m/internal/api/v1/core/application/services/user_service"
	"example.com/m/internal/api/v1/infrastructure/cache"
	database "example.com/m/internal/api/v1/infrastructure/database"
	"example.com/m/internal/api/v1/infrastructure/logger"
	"example.com/m/internal/api/v1/infrastructure/middlewares"
	"example.com/m/internal/api/v1/infrastructure/notifications"
	"example.com/m/internal/api/v1/infrastructure/router"
	object_storage "example.com/m/internal/api/v1/infrastructure/s3"
	"example.com/m/internal/api/v1/utils"
	"example.com/m/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Logger.Warn(("No .env file found"))
	}
}

func handlePanic() {
	if r := recover(); r != nil {
		logger.Logger.Fatal(
			"Handle panic",
			zap.Any("caller", debug.Stack()),
			zap.Any("error", r),
		)
		os.Exit(1)
	}
}

func startTokenUpdater(logger *zap.Logger) {
	err := utils.UpdateYandexIAMToken()
	if err != nil {
		logger.Fatal("Failed to update IAM token", zap.Error(err))
	}
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	defer handlePanic()

	for {
		select {
		case <-ticker.C:
			err := utils.UpdateYandexIAMToken()
			if err != nil {
				logger.Fatal("Failed to update IAM token", zap.Error(err))
			}
		}
	}
}

func main() {
	logger.NewLogger()
	loadEnv()
	config.InitConfig()
	database.ConnectToDatabase()
	object_storage.InitS3()
	database.MigrateDB()
	cache.ConnectToRedis()
	notifications.InitClient(context.Background())
	go startTokenUpdater(logger.Logger)
	metrics := middlewares.NewPrometheusMetrics()
	logger.Logger.Info("SERVICE STARTED")
	defer database.Db.Close()
	defer logger.Logger.Sync()

	fcmRepository := repositories.NewFcmRepository(notifications.FirebaseClient, logger.Logger)
	userRepository := repositories.NewUserRepository(database.Db, logger.Logger)
	tokenRepository := repositories.NewTokenRepository(cache.Redis, logger.Logger)
	postRepository := repositories.NewPostRepository(database.Db, logger.Logger, &object_storage.S3Client)
	placeRepository := repositories.NewPlaceRepository(database.Db, logger.Logger)
	bookingRepository := repositories.NewBookingRepository(database.Db, logger.Logger)
	reviewRepository := repositories.NewReviewRepository(database.Db, logger.Logger)
	pushTokenRepository := repositories.NewPushTokenRepository(cache.Redis, logger.Logger)
	chatRepository := repositories.NewChatRepository(database.Db, logger.Logger)

	gptService := gpt_service.NewGPTService(config.Config.YandexCatalogID, logger.Logger, chatRepository)
	userService := user_service.NewUserService(userRepository, pushTokenRepository)
	placeService := place_service.NewPlaceService(placeRepository)
	authService := auth_service.NewAuthService(userService, tokenRepository)
	postService := post_service.NewPostService(postRepository, placeService, userService, gptService)
	bookingService := booking_service.NewBookingService(*bookingRepository, *userService, *postService, fcmRepository, pushTokenRepository)
	reviewService := review_service.NewReviewService(reviewRepository, userService)

	authMiddleware := middlewares.NewAuthMiddleware(authService)
	adminMiddleware := middlewares.NewAdminMiddleware(userService)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	metricController := controllers.NewMetricController()
	postController := controllers.NewPostController(postService)
	placeController := controllers.NewPlaceController(placeService)
	bookingController := controllers.NewBookingController(bookingService)
	reviewController := controllers.NewReviewController(reviewService)
	chatBotController := controllers.NewChatBotController(gptService)

	// gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(metrics.Middleware())
	router := router.NewRouter(engine, authMiddleware, adminMiddleware)

	router.BindAuthRoutes(authController)
	router.BindMetricsRoutes(metricController)
	router.BindUserRoutes(userController)
	router.BindPlaceRoutes(placeController)
	router.BindSwaggerRoutes()
	router.BindPostRoutes(postController)
	router.BindBookingRoutes(bookingController)
	router.BindReviewRoutes(reviewController)
	router.BindChatBotRoutes(chatBotController)

	engine.Run(":8000")
}
