package app

import (
	"discord/config"

	authController "discord/internal/auth/controller"
	authRepo "discord/internal/auth/repository"
	authService "discord/internal/auth/service"

	friendRepo "discord/internal/friend/repository"
	friendService "discord/internal/friend/service"

	messageRepo "discord/internal/message/repository"
	messageService "discord/internal/message/service"

	serverRepo "discord/internal/server/repository"
	serverService "discord/internal/server/service"

	// syncController "discord/internal/sync/controller"
	// syncRepo "discord/internal/sync/repository"
	// syncService "discord/internal/sync/service"

	userRepo "discord/internal/user/repository"
	userService "discord/internal/user/service"

	voiceRepo "discord/internal/voice/repository"
	voiceService "discord/internal/voice/service"

	friendPb "discord/gen/proto/service/friend"
	messagePb "discord/gen/proto/service/message"
	serverPb "discord/gen/proto/service/server"
	userPb "discord/gen/proto/service/user"
	voicePb "discord/gen/proto/service/voice_channel"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Application holds all application dependencies
type Application struct {
	Config *config.Config
	DB     *pgxpool.Pool

	// Repositories
	AuthRepo    *authRepo.AuthRepository
	FriendRepo  *friendRepo.FriendRepository
	MessageRepo *messageRepo.MessageRepository
	ServerRepo  *serverRepo.ServerRepository
	// SyncRepo    *syncRepo.SyncRepository
	UserRepo  *userRepo.UserRepository
	VoiceRepo *voiceRepo.VoiceRepository

	// Services
	AuthSvc    *authService.AuthService
	FriendSvc  *friendService.FriendService
	MessageSvc *messageService.MessageService
	ServerSvc  *serverService.ServerService
	// SyncSvc    *syncService.SyncService
	UserSvc  *userService.UserService
	VoiceSvc *voiceService.VoiceService

	// Controllers
	AuthCtrl    *authController.AuthController
	FriendCtrl  *friendPb.FriendServiceServer
	MessageCtrl *messagePb.MessageServiceServer
	ServerCtrl  *serverPb.ServerServiceServer
	// SyncCtrl    *syncController.SyncController
	UserCtrl  *userPb.UserServiceServer
	VoiceCtrl *voicePb.VoiceChannelServiceServer
}
