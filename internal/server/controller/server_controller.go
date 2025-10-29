package controller

import (
	"context"

	"discord/gen/proto/schema"
	serverPb "discord/gen/proto/service/server"
	commonErrors "discord/internal/common/errors"
	serverService "discord/internal/server/service"
)

type ServerController struct {
	serverPb.UnimplementedServerServiceServer
	serverService *serverService.ServerService
}

func NewServerController(serverService *serverService.ServerService) *serverPb.ServerServiceServer {
	controller := &ServerController{
		serverService: serverService,
	}
	var grpcController serverPb.ServerServiceServer = controller
	return &grpcController
}

// CreateServer creates a new server
func (c *ServerController) CreateServer(ctx context.Context, req *serverPb.CreateServerRequest) (*serverPb.CreateServerResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetName() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	var icon, banner, description, region *string
	if req.GetIcon() != "" {
		i := req.GetIcon()
		icon = &i
	}
	if req.GetDescription() != "" {
		d := req.GetDescription()
		description = &d
	}
	if req.GetRegion() != "" {
		r := req.GetRegion()
		region = &r
	}

	server, err := c.serverService.CreateServer(ctx, req.GetName(), userID, icon, banner, description, region)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.CreateServerResponse{
		Server: &schema.Server{
			Id:      server.ID,
			Name:    server.Name,
			OwnerId: server.OwnerID,
		},
		Success: true,
	}, nil
}

// GetServer retrieves a server
func (c *ServerController) GetServer(ctx context.Context, req *serverPb.GetServerRequest) (*serverPb.GetServerResponse, error) {
	if req.GetServerId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	server, err := c.serverService.GetServer(ctx, req.GetServerId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbServer := &schema.Server{
		Id:      server.ID,
		Name:    server.Name,
		OwnerId: server.OwnerID,
	}

	if server.Icon.Valid {
		pbServer.Icon = server.Icon.String
	}
	if server.Banner.Valid {
		pbServer.Banner = server.Banner.String
	}
	if server.Description.Valid {
		pbServer.Description = server.Description.String
	}
	if server.Region.Valid {
		pbServer.Region = server.Region.String
	}

	return &serverPb.GetServerResponse{
		Server: pbServer,
	}, nil
}

// UpdateServer updates server information
func (c *ServerController) UpdateServer(ctx context.Context, req *serverPb.UpdateServerRequest) (*serverPb.UpdateServerResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetServerId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	var name, icon, banner, description, region *string
	if req.GetName() != "" {
		n := req.GetName()
		name = &n
	}
	if req.GetIcon() != "" {
		i := req.GetIcon()
		icon = &i
	}
	if req.GetBanner() != "" {
		b := req.GetBanner()
		banner = &b
	}
	if req.GetDescription() != "" {
		d := req.GetDescription()
		description = &d
	}
	if req.GetRegion() != "" {
		r := req.GetRegion()
		region = &r
	}

	server, err := c.serverService.UpdateServer(ctx, req.GetServerId(), userID, name, icon, banner, description, region)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.UpdateServerResponse{
		Server: &schema.Server{
			Id:      server.ID,
			Name:    server.Name,
			OwnerId: server.OwnerID,
		},
		Success: true,
	}, nil
}

// DeleteServer deletes a server
func (c *ServerController) DeleteServer(ctx context.Context, req *serverPb.DeleteServerRequest) (*serverPb.DeleteServerResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetServerId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.serverService.DeleteServer(ctx, req.GetServerId(), userID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.DeleteServerResponse{
		Success: true,
	}, nil
}

// GetUserServers retrieves all servers a user is in
func (c *ServerController) GetUserServers(ctx context.Context, req *serverPb.GetUserServersRequest) (*serverPb.GetUserServersResponse, error) {
	if req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	servers, err := c.serverService.GetUserServers(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbServers := make([]*schema.Server, len(servers))
	for i, server := range servers {
		pbServers[i] = &schema.Server{
			Id:      server.ID,
			Name:    server.Name,
			OwnerId: server.OwnerID,
		}
		if server.Icon.Valid {
			pbServers[i].Icon = server.Icon.String
		}
	}

	return &serverPb.GetUserServersResponse{
		Servers: pbServers,
	}, nil
}

// AddMember adds a member to a server (not typically used - JoinServerWithInvite is preferred)
func (c *ServerController) AddMember(ctx context.Context, req *serverPb.AddMemberRequest) (*serverPb.AddMemberResponse, error) {
	if req.GetServerId() == 0 || req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Use JoinServer without invite code
	err := c.serverService.JoinServer(ctx, req.GetServerId(), req.GetUserId(), nil)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.AddMemberResponse{
		Success: true,
	}, nil
}

// RemoveMember removes a member from a server
func (c *ServerController) RemoveMember(ctx context.Context, req *serverPb.RemoveMemberRequest) (*serverPb.RemoveMemberResponse, error) {
	if req.GetServerId() == 0 || req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.serverService.LeaveServer(ctx, req.GetServerId(), req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.RemoveMemberResponse{
		Success: true,
	}, nil
}

// GetMembers retrieves server members (streaming response)
func (c *ServerController) GetMembers(req *serverPb.GetMembersRequest, stream serverPb.ServerService_GetMembersServer) error {
	ctx := stream.Context()

	if req.GetServerId() == 0 {
		return commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	limit := req.GetLimit()
	if limit == 0 {
		limit = 100
	}

	members, err := c.serverService.GetServerMembers(ctx, req.GetServerId(), limit, req.GetOffset())
	if err != nil {
		return commonErrors.ToGRPCError(err)
	}

	// Stream members to client
	for _, member := range members {
		pbMember := &schema.ServerMember{
			ServerId: member.ServerID,
			UserId:   member.UserID,
			JoinedAt: member.JoinedAt.Time.Unix(),
		}

		if member.Nickname.Valid {
			pbMember.Nickname = member.Nickname.String
		}

		if err := stream.Send(&serverPb.GetMembersResponse{
			Members: []*schema.ServerMember{pbMember},
		}); err != nil {
			return err
		}
	}

	return nil
}

// UpdateMember updates member information
func (c *ServerController) UpdateMember(ctx context.Context, req *serverPb.UpdateMemberRequest) (*serverPb.UpdateMemberResponse, error) {
	if req.GetServerId() == 0 || req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	var nickname *string
	if req.GetNickname() != "" {
		n := req.GetNickname()
		nickname = &n
	}

	err := c.serverService.UpdateMemberNickname(ctx, req.GetServerId(), req.GetUserId(), nickname)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.UpdateMemberResponse{
		Success: true,
	}, nil
}

// KickMember kicks a member from the server
func (c *ServerController) KickMember(ctx context.Context, req *serverPb.KickMemberRequest) (*serverPb.KickMemberResponse, error) {
	// Get moderator ID from context
	moderatorID := ctx.Value("user_id").(int32)

	if req.GetServerId() == 0 || req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.serverService.KickMember(ctx, req.GetServerId(), moderatorID, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.KickMemberResponse{
		Success: true,
	}, nil
}

// BanMember bans a member from the server
func (c *ServerController) BanMember(ctx context.Context, req *serverPb.BanMemberRequest) (*serverPb.BanMemberResponse, error) {
	// Get moderator ID from context
	moderatorID := ctx.Value("user_id").(int32)

	if req.GetServerId() == 0 || req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	var reason *string
	if req.GetReason() != "" {
		r := req.GetReason()
		reason = &r
	}

	// Note: delete_message_days is ignored in this implementation
	err := c.serverService.BanMember(ctx, req.GetServerId(), moderatorID, req.GetUserId(), reason, nil)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.BanMemberResponse{
		Success: true,
	}, nil
}

// UnbanMember unbans a member from the server
func (c *ServerController) UnbanMember(ctx context.Context, req *serverPb.UnbanMemberRequest) (*serverPb.UnbanMemberResponse, error) {
	// Get moderator ID from context
	moderatorID := ctx.Value("user_id").(int32)

	if req.GetServerId() == 0 || req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.serverService.UnbanMember(ctx, req.GetServerId(), moderatorID, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.UnbanMemberResponse{
		Success: true,
	}, nil
}

// CreateInvite creates a server invite
func (c *ServerController) CreateInvite(ctx context.Context, req *serverPb.CreateInviteRequest) (*serverPb.CreateInviteResponse, error) {
	// Get inviter ID from context
	inviterID := ctx.Value("user_id").(int32)

	if req.GetServerId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	channelID := req.GetChannelId()
	if channelID == 0 {
		channelID = 1 // Default channel, adjust as needed
	}

	maxAge := req.GetMaxAge()
	maxUses := req.GetMaxUses()
	temporary := req.GetTemporary()

	invite, err := c.serverService.CreateInvite(ctx, req.GetServerId(), channelID, inviterID, maxUses, maxAge, temporary)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbInvite := &schema.Invite{
		Id:        invite.ID,
		Code:      invite.Code,
		ServerId:  invite.ServerID,
		InviterId: invite.InviterID,
		CreatedAt: invite.CreatedAt.Time.Unix(),
	}

	if invite.ChannelID.Valid {
		pbInvite.ChannelId = invite.ChannelID.Int32
	}
	if invite.MaxUses.Valid {
		pbInvite.MaxUses = invite.MaxUses.Int32
	}
	if invite.Uses.Valid {
		pbInvite.Uses = invite.Uses.Int32
	}

	return &serverPb.CreateInviteResponse{
		Invite:  pbInvite,
		Success: true,
	}, nil
}

// GetInvite retrieves an invite by code
func (c *ServerController) GetInvite(ctx context.Context, req *serverPb.GetInviteRequest) (*serverPb.GetInviteResponse, error) {
	if req.GetCode() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	invite, err := c.serverService.GetInvite(ctx, req.GetCode())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbInvite := &schema.Invite{
		Id:        invite.ID,
		Code:      invite.Code,
		ServerId:  invite.ServerID,
		InviterId: invite.InviterID,
		CreatedAt: invite.CreatedAt.Time.Unix(),
	}

	if invite.ChannelID.Valid {
		pbInvite.ChannelId = invite.ChannelID.Int32
	}
	if invite.MaxUses.Valid {
		pbInvite.MaxUses = invite.MaxUses.Int32
	}
	if invite.Uses.Valid {
		pbInvite.Uses = invite.Uses.Int32
	}

	return &serverPb.GetInviteResponse{
		Invite: pbInvite,
	}, nil
}

// DeleteInvite deletes an invite
func (c *ServerController) DeleteInvite(ctx context.Context, req *serverPb.DeleteInviteRequest) (*serverPb.DeleteInviteResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetCode() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.serverService.DeleteInvite(ctx, req.GetCode(), userID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.DeleteInviteResponse{
		Success: true,
	}, nil
}

// GetServerInvites retrieves all invites for a server
func (c *ServerController) GetServerInvites(ctx context.Context, req *serverPb.GetServerInvitesRequest) (*serverPb.GetServerInvitesResponse, error) {
	if req.GetServerId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	invites, err := c.serverService.GetServerInvites(ctx, req.GetServerId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbInvites := make([]*schema.Invite, len(invites))
	for i, invite := range invites {
		pbInvites[i] = &schema.Invite{
			Id:        invite.ID,
			Code:      invite.Code,
			ServerId:  invite.ServerID,
			InviterId: invite.InviterID,
			CreatedAt: invite.CreatedAt.Time.Unix(),
		}
		if invite.ChannelID.Valid {
			pbInvites[i].ChannelId = invite.ChannelID.Int32
		}
		if invite.MaxUses.Valid {
			pbInvites[i].MaxUses = invite.MaxUses.Int32
		}
		if invite.Uses.Valid {
			pbInvites[i].Uses = invite.Uses.Int32
		}
	}

	return &serverPb.GetServerInvitesResponse{
		Invites: pbInvites,
	}, nil
}

// JoinServerWithInvite allows a user to join a server using an invite code
func (c *ServerController) JoinServerWithInvite(ctx context.Context, req *serverPb.JoinServerWithInviteRequest) (*serverPb.JoinServerWithInviteResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetCode() == "" {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	// Get invite to find server
	invite, err := c.serverService.GetInvite(ctx, req.GetCode())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Join server with invite code
	code := req.GetCode()
	err = c.serverService.JoinServer(ctx, invite.ServerID, userID, &code)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	// Get server info
	server, err := c.serverService.GetServer(ctx, invite.ServerID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &serverPb.JoinServerWithInviteResponse{
		Server: &schema.Server{
			Id:      server.ID,
			Name:    server.Name,
			OwnerId: server.OwnerID,
		},
		Success: true,
	}, nil
}

// CreateEmoji creates a custom emoji
func (c *ServerController) CreateEmoji(ctx context.Context, req *serverPb.CreateEmojiRequest) (*serverPb.CreateEmojiResponse, error) {
	// TODO: Implement emoji creation
	return &serverPb.CreateEmojiResponse{
		Success: true,
	}, nil
}

// DeleteEmoji deletes a custom emoji
func (c *ServerController) DeleteEmoji(ctx context.Context, req *serverPb.DeleteEmojiRequest) (*serverPb.DeleteEmojiResponse, error) {
	// TODO: Implement emoji deletion
	return &serverPb.DeleteEmojiResponse{
		Success: true,
	}, nil
}

// GetServerEmojis retrieves all emojis for a server
func (c *ServerController) GetServerEmojis(ctx context.Context, req *serverPb.GetServerEmojisRequest) (*serverPb.GetServerEmojisResponse, error) {
	// TODO: Implement emoji retrieval
	return &serverPb.GetServerEmojisResponse{
		Emojis: []*schema.Emoji{},
	}, nil
}
