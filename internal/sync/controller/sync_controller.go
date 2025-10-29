package controller

import (
	"context"

	syncPb "discord/gen/proto/service/sync"
	"discord/internal/sync/service"
	"discord/internal/sync/util"
)

type SyncController struct {
	syncPb.UnimplementedSyncServiceServer
	syncService *service.SyncService
}

func NewSyncController(syncService *service.SyncService) *SyncController {
	return &SyncController{
		syncService: syncService,
	}
}

// SyncFriends handles friend sync requests
func (c *SyncController) SyncFriends(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncFriendsResponse, error) {
	// Extract user ID from context
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	// Set defaults
	limit := req.Limit
	if limit == 0 {
		limit = 100
	}
	offset := req.Offset

	// Call service
	result, err := c.syncService.SyncFriends(ctx, userID, req.LastUpdatedAt, limit, offset)
	if err != nil {
		return nil, err
	}

	// Extract data from result map
	data := result.(map[string]interface{})

	return &syncPb.SyncFriendsResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int64)),
		Message:         data["message"].(string),
	}, nil
}

// SyncMessages handles message sync requests
func (c *SyncController) SyncMessages(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncMessagesResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	limit := req.Limit
	if limit == 0 {
		limit = 100
	}
	offset := req.Offset

	result, err := c.syncService.SyncMessages(ctx, userID, req.LastUpdatedAt, limit, offset)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncMessagesResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int64)),
		HasMore:         data["has_more"].(bool),
		Message:         data["message"].(string),
	}, nil
}

// SyncServers handles server sync requests
func (c *SyncController) SyncServers(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncServersResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	result, err := c.syncService.SyncServers(ctx, userID, req.LastUpdatedAt)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncServersResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int64)),
		Message:         data["message"].(string),
	}, nil
}

// SyncChannels handles channel sync requests
func (c *SyncController) SyncChannels(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncChannelsResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	result, err := c.syncService.SyncChannels(ctx, userID, req.LastUpdatedAt)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncChannelsResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int64)),
		Message:         data["message"].(string),
	}, nil
}

// SyncUserProfile handles user profile sync requests
func (c *SyncController) SyncUserProfile(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncUserProfileResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	result, err := c.syncService.SyncUserProfile(ctx, userID, req.LastUpdatedAt)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncUserProfileResponse{
		ServerTimestamp:   data["server_timestamp"].(int64),
		IsSynced:          data["is_synced"].(bool),
		LastProfileUpdate: data["last_profile_update"].(int64),
		Message:           data["message"].(string),
	}, nil
}

// SyncVoiceChannels handles voice channel sync requests
func (c *SyncController) SyncVoiceChannels(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncVoiceChannelsResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	result, err := c.syncService.SyncVoiceChannels(ctx, userID, req.LastUpdatedAt)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncVoiceChannelsResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int)),
		Message:         data["message"].(string),
	}, nil
}

// SyncTextChannels handles text channel sync requests
func (c *SyncController) SyncTextChannels(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncTextChannelsResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	result, err := c.syncService.SyncTextChannels(ctx, userID, req.LastUpdatedAt)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncTextChannelsResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int)),
		Message:         data["message"].(string),
	}, nil
}

// SyncDirectMessages handles direct message sync requests
func (c *SyncController) SyncDirectMessages(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncDirectMessagesResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	limit := req.Limit
	if limit == 0 {
		limit = 50
	}
	offset := req.Offset

	result, err := c.syncService.SyncDirectMessages(ctx, userID, req.LastUpdatedAt, limit, offset)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncDirectMessagesResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int)),
		HasMore:         data["has_more"].(bool),
		Message:         data["message"].(string),
	}, nil
}

// SyncPermissions handles permission sync requests
func (c *SyncController) SyncPermissions(ctx context.Context, req *syncPb.SyncDataRequest) (*syncPb.SyncPermissionsResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	result, err := c.syncService.SyncPermissions(ctx, userID, req.LastUpdatedAt)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	return &syncPb.SyncPermissionsResponse{
		ServerTimestamp: data["server_timestamp"].(int64),
		IsSynced:        data["is_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int)),
		Message:         data["message"].(string),
	}, nil
}

// SyncAll handles sync all request
func (c *SyncController) SyncAll(ctx context.Context, req *syncPb.SyncAllRequest) (*syncPb.SyncAllResponse, error) {
	userID, ok := ctx.Value("user_id").(int32)
	if !ok {
		userID = req.UserId
	}

	// Build timestamps map
	timestamps := map[string]int64{
		"friends":         req.FriendsLastUpdatedAt,
		"messages":        req.MessagesLastUpdatedAt,
		"servers":         req.ServersLastUpdatedAt,
		"channels":        req.ChannelsLastUpdatedAt,
		"user_profile":    req.UserProfileLastUpdatedAt,
		"voice_channels":  req.VoiceChannelsLastUpdatedAt,
		"text_channels":   req.TextChannelsLastUpdatedAt,
		"direct_messages": req.DirectMessagesLastUpdatedAt,
		"permissions":     req.PermissionsLastUpdatedAt,
	}

	result, err := c.syncService.SyncAll(ctx, userID, timestamps)
	if err != nil {
		return nil, err
	}

	data := result.(map[string]interface{})

	// Convert sub-responses
	friends := data["friends"].(map[string]interface{})
	messages := data["messages"].(map[string]interface{})
	servers := data["servers"].(map[string]interface{})
	channels := data["channels"].(map[string]interface{})
	userProfile := data["user_profile"].(map[string]interface{})
	voiceChannels := data["voice_channels"].(map[string]interface{})
	textChannels := data["text_channels"].(map[string]interface{})
	directMessages := data["direct_messages"].(map[string]interface{})
	permissions := data["permissions"].(map[string]interface{})

	return &syncPb.SyncAllResponse{
		Friends: &syncPb.SyncFriendsResponse{
			ServerTimestamp: friends["server_timestamp"].(int64),
			IsSynced:        friends["is_synced"].(bool),
			TotalUpdates:    int32(friends["total_updates"].(int64)),
			Message:         friends["message"].(string),
		},
		Messages: &syncPb.SyncMessagesResponse{
			ServerTimestamp: messages["server_timestamp"].(int64),
			IsSynced:        messages["is_synced"].(bool),
			TotalUpdates:    int32(messages["total_updates"].(int64)),
			HasMore:         messages["has_more"].(bool),
			Message:         messages["message"].(string),
		},
		Servers: &syncPb.SyncServersResponse{
			ServerTimestamp: servers["server_timestamp"].(int64),
			IsSynced:        servers["is_synced"].(bool),
			TotalUpdates:    int32(servers["total_updates"].(int64)),
			Message:         servers["message"].(string),
		},
		Channels: &syncPb.SyncChannelsResponse{
			ServerTimestamp: channels["server_timestamp"].(int64),
			IsSynced:        channels["is_synced"].(bool),
			TotalUpdates:    int32(channels["total_updates"].(int64)),
			Message:         channels["message"].(string),
		},
		UserProfile: &syncPb.SyncUserProfileResponse{
			ServerTimestamp:   userProfile["server_timestamp"].(int64),
			IsSynced:          userProfile["is_synced"].(bool),
			LastProfileUpdate: userProfile["last_profile_update"].(int64),
			Message:           userProfile["message"].(string),
		},
		VoiceChannels: &syncPb.SyncVoiceChannelsResponse{
			ServerTimestamp: voiceChannels["server_timestamp"].(int64),
			IsSynced:        voiceChannels["is_synced"].(bool),
			TotalUpdates:    int32(voiceChannels["total_updates"].(int)),
			Message:         voiceChannels["message"].(string),
		},
		TextChannels: &syncPb.SyncTextChannelsResponse{
			ServerTimestamp: textChannels["server_timestamp"].(int64),
			IsSynced:        textChannels["is_synced"].(bool),
			TotalUpdates:    int32(textChannels["total_updates"].(int)),
			Message:         textChannels["message"].(string),
		},
		DirectMessages: &syncPb.SyncDirectMessagesResponse{
			ServerTimestamp: directMessages["server_timestamp"].(int64),
			IsSynced:        directMessages["is_synced"].(bool),
			TotalUpdates:    int32(directMessages["total_updates"].(int)),
			HasMore:         directMessages["has_more"].(bool),
			Message:         directMessages["message"].(string),
		},
		Permissions: &syncPb.SyncPermissionsResponse{
			ServerTimestamp: permissions["server_timestamp"].(int64),
			IsSynced:        permissions["is_synced"].(bool),
			TotalUpdates:    int32(permissions["total_updates"].(int)),
			Message:         permissions["message"].(string),
		},
		ServerTimestamp: data["server_timestamp"].(int64),
		IsFullySynced:   data["is_fully_synced"].(bool),
		TotalUpdates:    int32(data["total_updates"].(int)),
		Message:         data["message"].(string),
	}, nil
}

// Helper function used internally
func getUserIDFromContext(ctx context.Context, fallbackID int32) int32 {
	if userID, ok := ctx.Value("user_id").(int32); ok {
		return userID
	}
	return fallbackID
}

// Initialize helper reference
var _ = util.ConvertPgTimestampToMillis
