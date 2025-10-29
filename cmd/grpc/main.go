package main

import (
	"fmt"
	"log"
	"net"

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

	userController "discord/internal/user/controller"
	userRepo "discord/internal/user/repository"
	userService "discord/internal/user/service"

	voiceController "discord/internal/voice/controller"
	voiceRepo "discord/internal/voice/repository"
	voiceService "discord/internal/voice/service"

	syncController "discord/internal/sync/controller"
	syncRepo "discord/internal/sync/repository"
	syncService "discord/internal/sync/service"

	authPb "discord/gen/proto/service/auth"
	friendPb "discord/gen/proto/service/friend"
	messagePb "discord/gen/proto/service/message"
	serverPb "discord/gen/proto/service/server"
	syncPb "discord/gen/proto/service/sync"
	userPb "discord/gen/proto/service/user"
	voicePb "discord/gen/proto/service/voice_channel"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Application struct {
	Config *config.Config
	DB     *pgxpool.Pool

	// Repositories
	AuthRepo    *authRepo.AuthRepository
	FriendRepo  *friendRepo.FriendRepository
	MessageRepo *messageRepo.MessageRepository
	ServerRepo  *serverRepo.ServerRepository
	SyncRepo    *syncRepo.SyncRepository
	UserRepo    *userRepo.UserRepository
	VoiceRepo   *voiceRepo.VoiceRepository

	// Services
	AuthSvc    *authService.AuthService
	FriendSvc  *friendService.FriendService
	MessageSvc *messageService.MessageService
	ServerSvc  *serverService.ServerService
	SyncSvc    *syncService.SyncService
	UserSvc    *userService.UserService
	VoiceSvc   *voiceService.VoiceService

	// Controllers
	AuthCtrl    *authController.AuthController
	FriendCtrl  *friendPb.FriendServiceServer
	MessageCtrl *messagePb.MessageServiceServer
	ServerCtrl  *serverPb.ServerServiceServer
	SyncCtrl    *syncController.SyncController
	UserCtrl    *userPb.UserServiceServer
	VoiceCtrl   *voicePb.VoiceChannelServiceServer
}

func main() {
	app := &Application{}

	// Initialize application
	if err := app.Initialize(); err != nil {
		log.Fatalf("❌ Failed to initialize application: %v", err)
	}
	defer app.Shutdown()

	// Start gRPC server
	if err := app.StartServer(); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

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

func (app *Application) initRepositories() {
	app.AuthRepo = authRepo.NewAuthRepository(app.DB)
	app.FriendRepo = friendRepo.NewFriendRepository(app.DB)
	app.MessageRepo = messageRepo.NewMessageRepository(app.DB)
	app.ServerRepo = serverRepo.NewServerRepository(app.DB)
	app.SyncRepo = syncRepo.NewSyncRepository(app.DB)
	app.UserRepo = userRepo.NewUserRepository(app.DB)
	app.VoiceRepo = voiceRepo.NewVoiceRepository(app.DB)
}

func (app *Application) initServices() {
	app.AuthSvc = authService.NewAuthService(app.AuthRepo)
	app.FriendSvc = friendService.NewFriendService(app.FriendRepo)
	app.MessageSvc = messageService.NewMessageService(app.MessageRepo)
	app.ServerSvc = serverService.NewServerService(app.ServerRepo)
	app.SyncSvc = syncService.NewSyncService(app.SyncRepo)
	app.UserSvc = userService.NewUserService(app.UserRepo)
	app.VoiceSvc = voiceService.NewVoiceService(app.VoiceRepo)
}

func (app *Application) initControllers() {
	app.AuthCtrl = authController.NewAuthController(app.AuthSvc)
	app.FriendCtrl = friendController.NewFriendController(app.FriendSvc)
	app.MessageCtrl = messageController.NewMessageController(app.MessageSvc)
	app.ServerCtrl = serverController.NewServerController(app.ServerSvc)
	app.SyncCtrl = syncController.NewSyncController(app.SyncSvc)
	app.UserCtrl = userController.NewUserController(app.UserSvc)
	app.VoiceCtrl = voiceController.NewVoiceController(app.VoiceSvc)
}

func (app *Application) StartServer() error {
	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register all services
	authPb.RegisterAuthServiceServer(grpcServer, app.AuthCtrl)
	friendPb.RegisterFriendServiceServer(grpcServer, *app.FriendCtrl)
	messagePb.RegisterMessageServiceServer(grpcServer, *app.MessageCtrl)
	serverPb.RegisterServerServiceServer(grpcServer, *app.ServerCtrl)
	syncPb.RegisterSyncServiceServer(grpcServer, app.SyncCtrl)
	userPb.RegisterUserServiceServer(grpcServer, *app.UserCtrl)
	voicePb.RegisterVoiceChannelServiceServer(grpcServer, *app.VoiceCtrl)

	log.Println("✅ gRPC services registered")

	// Get port
	port := app.Config.Service.Port
	if port == "" {
		port = "50051"
	}

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	// Print startup info
	app.printStartupInfo(port)
	// Register reflection for development
	reflection.Register(grpcServer)
	// Start serving
	return grpcServer.Serve(listener)
}

func (app *Application) printStartupInfo(port string) {
	separator := "============================================================"
	log.Println("\n" + separator)
	log.Println("🚀 Discord gRPC Server")
	log.Println(separator)
	log.Printf("📡 Port:        %s", port)
	log.Printf("🌍 Environment: %s", app.Config.Service.Environment)
	log.Printf("🗄️  Database:    Connected")
	log.Println("\n📦 Registered Services:")
	log.Println("   ✓ AuthService         - User registration & authentication")
	log.Println("   ✓ FriendService       - Friend management & requests")
	log.Println("   ✓ MessageService      - Messages, reactions, attachments")
	log.Println("   ✓ ServerService       - Servers, members, roles, invites")
	log.Println("   ✓ SyncService         - Real-time data synchronization")
	log.Println("   ✓ UserService         - User profiles & settings")
	log.Println("   ✓ VoiceChannelService - Voice states & connections")
	log.Println("\n✨ Server is ready to accept connections!")
	log.Println(separator + "\n")
}

func (app *Application) Shutdown() {
	if app.DB != nil {
		app.DB.Close()
		log.Println("✅ Database connection closed")
	}
}
