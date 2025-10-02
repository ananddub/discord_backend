package server

import (
	"fmt"
	"log"
	"net"

	"buf.build/go/protovalidate"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"discord/config"
	authpb "discord/gen/proto/service/auth"
	channelpb "discord/gen/proto/service/channel"
	friendpb "discord/gen/proto/service/friend"
	messagepb "discord/gen/proto/service/message"
	textchannelpb "discord/gen/proto/service/text_channel"
	userpb "discord/gen/proto/service/user"
	authHandler "discord/internal/auth/handler"
	channelHandler "discord/internal/channel/handler"
	friendHandler "discord/internal/friend/handler"
	messageHandler "discord/internal/message/handler"
	textChannelHandler "discord/internal/text_channel/handler"
	userHandler "discord/internal/user/handler"
)

type Server struct {
	config *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{config: cfg}
}

func (s *Server) Start() error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err = config.InitDB(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	fmt.Println("db connected")

	lis, err := net.Listen("tcp", cfg.Service.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to create protovalidate validator: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
	)
	registerHandlers(grpcServer)
	reflection.Register(grpcServer)

	fmt.Printf("gRPC server running on %s\n", cfg.Service.Port)
	return grpcServer.Serve(lis)
}

func registerHandlers(server *grpc.Server) {
	authpb.RegisterAuthServiceServer(server, authHandler.NewAuthHandler())
	userpb.RegisterUserServiceServer(server, userHandler.NewUserHandler())
	channelpb.RegisterChannelServiceServer(server, channelHandler.NewChannelHandler())
	messagepb.RegisterMessageServiceServer(server, messageHandler.NewMessageHandler())
	friendpb.RegisterFriendServiceServer(server, friendHandler.NewFriendHandler())
	textchannelpb.RegisterTextChannelServiceServer(server, textChannelHandler.NewTextChannelHandler())
}
