package main

import (
	"auth-service/api"
	"auth-service/cluster"
	"auth-service/cluster/userservice"
	"auth-service/database"
	"auth-service/service"
	"auth-service/settings"
	"context"
	"fmt"
	"github.com/GOAT-prod/goatlogger"
	"github.com/redis/go-redis/v9"
	"time"
)

type App struct {
	mainCtx context.Context
	cfg     settings.Config
	logger  goatlogger.Logger

	server            *api.Server
	redisClient       *redis.Client
	userServiceClient *userservice.Client

	hasher              service.PasswordHasher
	jwtService          service.JwtService
	authService         service.AuthService
	registrationService service.Registration

	refreshTokenRepository database.RefreshTokenRepository
}

func NewApp(ctx context.Context, cfg settings.Config, logger goatlogger.Logger) *App {
	return &App{
		mainCtx: ctx,
		cfg:     cfg,
		logger:  logger,
	}
}

func (a *App) Start() {
	go a.server.Start(a.cfg.Port)
}

func (a *App) Stop(ctx context.Context) {
	_, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := a.redisClient.Close(); err != nil {
		a.logger.Panic(fmt.Sprintf("не удалось закрыть подключение к redis: %v", err))
	}
}

func (a *App) initDatabases() {
	a.initRedis()
}

func (a *App) initRedis() {
	redisCtx, cancelFunc := context.WithTimeout(a.mainCtx, 15*time.Second)
	defer cancelFunc()

	redisClient, err := database.NewRedisClient(redisCtx, a.cfg.Databases.Redis)
	if err != nil {
		a.logger.Panic(fmt.Sprintf("не удалось подключиться к redis: %v", err))
	}

	a.redisClient = redisClient
}

func (a *App) initRepositories() {
	a.refreshTokenRepository = database.NewRefreshTokenRepository(a.redisClient)
}

func (a *App) initClusterClients() {
	a.userServiceClient = userservice.NewClient(cluster.NewBaseClient(a.cfg.ClusterClients.UserService))
}

func (a *App) initServices() {
	a.hasher = service.NewHasher(a.cfg.ClientSecret)
	a.jwtService = service.NewJwtService(a.cfg.JwtSecret, a.refreshTokenRepository)
	a.authService = service.NewAuthService(a.hasher, a.jwtService, a.userServiceClient)
	a.registrationService = service.NewRegistrationService(a.hasher, a.jwtService, a.userServiceClient)
}

func (a *App) initServer() {
	if a.server != nil {
		a.logger.Panic("сервер уже запущен")
	}

	a.server = api.NewServer()
	a.server.Setup(a.authService, a.registrationService)
}
