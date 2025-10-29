package main

import (
	"fmt"
	"log"

	"discord/config"

	authController "discord/internal/auth/controller"
	authRepo "discord/internal/auth/repository"
	authService "discord/internal/auth/service"

	friendController "discord/internal/friend/controller"
	friendRepo "discord/internal/friend/repository"
	friendService "discord/internal/friend/service"

	messageController "discord/internal/message/controller"
	messageRepo "discord/internal/message/repository"
	messageService "discord/internal/message/service"

	serverController "discord/internal/server/controller"
	serverRepo "discord/internal/server/repository"
	serverService "discord/internal/server/service"

	syncController "discord/internal/sync/controller"
	syncRepo "discord/internal/sync/repository"
	syncService "discord/internal/sync/service"

	userController "discord/internal/user/controller"
	userRepo "discord/internal/user/repository"
	userService "discord/internal/user/service"

	voiceController "discord/internal/voice/controller"
	voiceRepo "discord/internal/voice/repository"
	voiceService "discord/internal/voice/service"
)

// Initialize initializes all application dependencies
func (app *Application) Initialize() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	app.Config = cfg
	log.Println("✅ Configuration loaded")

	// Initialize database
	if err := config.InitDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	app.DB = config.DB
	log.Println("✅ Database connected")
	if err := config.InitValKey(); err != nil {
		return fmt.Errorf("failed to initialize validation key: %w", err)
	}
	// Initialize repositories
	app.initRepositories()
	log.Println("✅ Repositories initialized")

	// Initialize services
	app.initServices()
	log.Println("✅ Services initialized")

	// Initialize controllers
	app.initControllers()
	log.Println("✅ Controllers initialized")

	return nil
}

// initRepositories initializes all repository instances
func (app *Application) initRepositories() {
	app.AuthRepo = authRepo.NewAuthRepository(app.DB)
	app.FriendRepo = friendRepo.NewFriendRepository(app.DB)
	app.MessageRepo = messageRepo.NewMessageRepository(app.DB)
	app.ServerRepo = serverRepo.NewServerRepository(app.DB)
	app.SyncRepo = syncRepo.NewSyncRepository(app.DB)
	app.UserRepo = userRepo.NewUserRepository(app.DB)
	app.VoiceRepo = voiceRepo.NewVoiceRepository(app.DB)
}

// initServices initializes all service instances
func (app *Application) initServices() {
	app.AuthSvc = authService.NewAuthService(app.AuthRepo)
	app.FriendSvc = friendService.NewFriendService(app.FriendRepo)
	app.MessageSvc = messageService.NewMessageService(app.MessageRepo)
	app.ServerSvc = serverService.NewServerService(app.ServerRepo)
	app.SyncSvc = syncService.NewSyncService(app.SyncRepo)
	app.UserSvc = userService.NewUserService(app.UserRepo)
	app.VoiceSvc = voiceService.NewVoiceService(app.VoiceRepo)
}

// initControllers initializes all controller instances
func (app *Application) initControllers() {
	app.AuthCtrl = authController.NewAuthController(app.AuthSvc)
	app.FriendCtrl = friendController.NewFriendController(app.FriendSvc)
	app.MessageCtrl = messageController.NewMessageController(app.MessageSvc)
	app.ServerCtrl = serverController.NewServerController(app.ServerSvc)
	app.SyncCtrl = syncController.NewSyncController(app.SyncSvc)
	app.UserCtrl = userController.NewUserController(app.UserSvc)
	app.VoiceCtrl = voiceController.NewVoiceController(app.VoiceSvc)
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() {
	if app.DB != nil {
		app.DB.Close()
		log.Println("✅ Database connection closed")
	}
}
