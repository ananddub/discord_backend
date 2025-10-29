package main

import (
	"fmt"
	"log"
	"net"

	authPb "discord/gen/proto/service/auth"
	friendPb "discord/gen/proto/service/friend"
	messagePb "discord/gen/proto/service/message"
	serverPb "discord/gen/proto/service/server"
	syncPb "discord/gen/proto/service/sync"
	userPb "discord/gen/proto/service/user"
	voicePb "discord/gen/proto/service/voice_channel"
	"discord/internal/common/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartServer starts the gRPC server
func (app *Application) StartServer() error {
	// Create gRPC server with interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RecoveryInterceptor(), // Panic recovery (first)
			middleware.LoggingInterceptor(),  // Request logging
			// middleware.RateLimitInterceptor(), // Rate limiting
			middleware.AuthInterceptor(), // Authentication (last)
		),
		grpc.ChainStreamInterceptor(
			middleware.StreamRecoveryInterceptor(), // Panic recovery for streams
			middleware.StreamLoggingInterceptor(),  // Logging for streams
		),
	)

	// Register all services
	app.registerServices(grpcServer)
	log.Println("✅ gRPC services registered")
	log.Println("✅ Interceptors enabled: Recovery, Logging, RateLimit, Auth")

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

// registerServices registers all gRPC services
func (app *Application) registerServices(grpcServer *grpc.Server) {
	authPb.RegisterAuthServiceServer(grpcServer, app.AuthCtrl)
	friendPb.RegisterFriendServiceServer(grpcServer, *app.FriendCtrl)
	messagePb.RegisterMessageServiceServer(grpcServer, *app.MessageCtrl)
	serverPb.RegisterServerServiceServer(grpcServer, *app.ServerCtrl)
	syncPb.RegisterSyncServiceServer(grpcServer, app.SyncCtrl)
	userPb.RegisterUserServiceServer(grpcServer, *app.UserCtrl)
	voicePb.RegisterVoiceChannelServiceServer(grpcServer, *app.VoiceCtrl)
}

// printStartupInfo prints server startup information
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
