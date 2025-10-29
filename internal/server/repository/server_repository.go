package repository

import (
	"context"

	"discord/gen/repo"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerRepository struct {
	db      *pgxpool.Pool
	queries *repo.Queries
}

func NewServerRepository(db *pgxpool.Pool) *ServerRepository {
	return &ServerRepository{
		db:      db,
		queries: repo.New(db),
	}
}

// CreateServer creates a new server
func (r *ServerRepository) CreateServer(ctx context.Context, name string, ownerID int32, icon, banner, description, region *string) (repo.Server, error) {
	var iconType, bannerType, descType, regionType pgtype.Text

	if icon != nil {
		iconType = pgtype.Text{String: *icon, Valid: true}
	}
	if banner != nil {
		bannerType = pgtype.Text{String: *banner, Valid: true}
	}
	if description != nil {
		descType = pgtype.Text{String: *description, Valid: true}
	}
	if region != nil {
		regionType = pgtype.Text{String: *region, Valid: true}
	}

	return r.queries.CreateServer(ctx, repo.CreateServerParams{
		Name:        name,
		Icon:        iconType,
		Banner:      bannerType,
		Description: descType,
		OwnerID:     ownerID,
		Region:      regionType,
	})
}

// GetServerByID retrieves a server by ID
func (r *ServerRepository) GetServerByID(ctx context.Context, serverID int32) (repo.Server, error) {
	return r.queries.GetServerByID(ctx, serverID)
}

// UpdateServer updates server information
func (r *ServerRepository) UpdateServer(ctx context.Context, serverID int32, name, icon, banner, description, region *string) (repo.Server, error) {
	var nameType, iconType, bannerType, descType, regionType pgtype.Text

	if name != nil {
		nameType = pgtype.Text{String: *name, Valid: true}
	}
	if icon != nil {
		iconType = pgtype.Text{String: *icon, Valid: true}
	}
	if banner != nil {
		bannerType = pgtype.Text{String: *banner, Valid: true}
	}
	if description != nil {
		descType = pgtype.Text{String: *description, Valid: true}
	}
	if region != nil {
		regionType = pgtype.Text{String: *region, Valid: true}
	}

	return r.queries.UpdateServer(ctx, repo.UpdateServerParams{
		ID:          serverID,
		Name:        nameType,
		Icon:        iconType,
		Banner:      bannerType,
		Description: descType,
		Region:      regionType,
	})
}

// DeleteServer deletes a server
func (r *ServerRepository) DeleteServer(ctx context.Context, serverID int32) error {
	return r.queries.SoftDeleteServer(ctx, serverID)
}

// GetUserServers retrieves all servers a user is a member of
func (r *ServerRepository) GetUserServers(ctx context.Context, userID int32) ([]repo.Server, error) {
	return r.queries.GetUserServers(ctx, userID)
}

// IncrementMemberCount increments the server member count
func (r *ServerRepository) IncrementMemberCount(ctx context.Context, serverID int32) error {
	return r.queries.IncrementMemberCount(ctx, serverID)
}

// DecrementMemberCount decrements the server member count
func (r *ServerRepository) DecrementMemberCount(ctx context.Context, serverID int32) error {
	return r.queries.DecrementMemberCount(ctx, serverID)
}

// GetServersByOwner retrieves servers owned by a user
func (r *ServerRepository) GetServersByOwner(ctx context.Context, ownerID int32) ([]repo.Server, error) {
	return r.queries.GetServersByOwner(ctx, ownerID)
}

// UpdateServerOwner updates the server owner
func (r *ServerRepository) UpdateServerOwner(ctx context.Context, serverID, newOwnerID int32) error {
	return r.queries.UpdateServerOwner(ctx, repo.UpdateServerOwnerParams{
		ID:      serverID,
		OwnerID: newOwnerID,
	})
}

// AddServerMember adds a member to a server
func (r *ServerRepository) AddServerMember(ctx context.Context, serverID, userID int32, nickname *string) (repo.ServerMember, error) {
	var nicknameType pgtype.Text
	if nickname != nil {
		nicknameType = pgtype.Text{String: *nickname, Valid: true}
	}

	return r.queries.AddServerMember(ctx, repo.AddServerMemberParams{
		ServerID: serverID,
		UserID:   userID,
		Nickname: nicknameType,
	})
}

// GetServerMember retrieves a specific server member
func (r *ServerRepository) GetServerMember(ctx context.Context, serverID, userID int32) (repo.ServerMember, error) {
	return r.queries.GetServerMember(ctx, repo.GetServerMemberParams{
		ServerID: serverID,
		UserID:   userID,
	})
}

// GetServerMembers retrieves all members of a server with pagination
func (r *ServerRepository) GetServerMembers(ctx context.Context, serverID int32, limit, offset int32) ([]repo.ServerMember, error) {
	return r.queries.GetServerMembers(ctx, repo.GetServerMembersParams{
		ServerID: serverID,
		Limit:    limit,
		Offset:   offset,
	})
}

// UpdateMemberNickname updates a member's nickname
func (r *ServerRepository) UpdateMemberNickname(ctx context.Context, serverID, userID int32, nickname *string) error {
	var nicknameType pgtype.Text
	if nickname != nil {
		nicknameType = pgtype.Text{String: *nickname, Valid: true}
	}

	return r.queries.UpdateMemberNickname(ctx, repo.UpdateMemberNicknameParams{
		ServerID: serverID,
		UserID:   userID,
		Nickname: nicknameType,
	})
}

// RemoveServerMember removes a member from a server
func (r *ServerRepository) RemoveServerMember(ctx context.Context, serverID, userID int32) error {
	return r.queries.RemoveServerMember(ctx, repo.RemoveServerMemberParams{
		ServerID: serverID,
		UserID:   userID,
	})
}

// CountServerMembers counts the number of members in a server
func (r *ServerRepository) CountServerMembers(ctx context.Context, serverID int32) (int64, error) {
	return r.queries.CountServerMembers(ctx, serverID)
}

// CreateRole creates a new role
func (r *ServerRepository) CreateRole(ctx context.Context, serverID int32, name string, color *string, hoist, mentionable bool, position int32, permissions int64, description *string) (repo.Role, error) {
	var colorType, descType pgtype.Text

	if color != nil {
		colorType = pgtype.Text{String: *color, Valid: true}
	}
	if description != nil {
		descType = pgtype.Text{String: *description, Valid: true}
	}

	return r.queries.CreateRole(ctx, repo.CreateRoleParams{
		ServerID:    serverID,
		Name:        name,
		Color:       colorType,
		Hoist:       pgtype.Bool{Bool: hoist, Valid: true},
		Position:    pgtype.Int4{Int32: position, Valid: true},
		Permissions: pgtype.Int8{Int64: permissions, Valid: true},
		Mentionable: pgtype.Bool{Bool: mentionable, Valid: true},
		Description: descType,
	})
} // GetRoleByID retrieves a role by ID
func (r *ServerRepository) GetRoleByID(ctx context.Context, roleID int32) (repo.Role, error) {
	return r.queries.GetRoleByID(ctx, roleID)
}

// GetServerRoles retrieves all roles for a server
func (r *ServerRepository) GetServerRoles(ctx context.Context, serverID int32) ([]repo.Role, error) {
	return r.queries.GetServerRoles(ctx, serverID)
}

// DeleteRole deletes a role
func (r *ServerRepository) DeleteRole(ctx context.Context, roleID int32) error {
	return r.queries.SoftDeleteRole(ctx, roleID)
}

// CreateInvite creates a server invite
func (r *ServerRepository) CreateInvite(ctx context.Context, code string, serverID, channelID, inviterID int32, maxUses, maxAge int32, temporary bool, expiresAt *pgtype.Timestamp) (repo.Invite, error) {
	var expiresAtType pgtype.Timestamp
	if expiresAt != nil {
		expiresAtType = *expiresAt
	}

	return r.queries.CreateInvite(ctx, repo.CreateInviteParams{
		Code:      code,
		ServerID:  serverID,
		ChannelID: pgtype.Int4{Int32: channelID, Valid: true},
		InviterID: inviterID,
		MaxUses:   pgtype.Int4{Int32: maxUses, Valid: true},
		MaxAge:    pgtype.Int4{Int32: maxAge, Valid: true},
		Temporary: pgtype.Bool{Bool: temporary, Valid: true},
		ExpiresAt: expiresAtType,
	})
}

// GetInviteByCode retrieves an invite by code
func (r *ServerRepository) GetInviteByCode(ctx context.Context, code string) (repo.Invite, error) {
	return r.queries.GetInviteByCode(ctx, code)
}

// GetServerInvites retrieves all invites for a server
func (r *ServerRepository) GetServerInvites(ctx context.Context, serverID int32) ([]repo.Invite, error) {
	return r.queries.GetServerInvites(ctx, serverID)
}

// IncrementInviteUses increments the invite usage count
func (r *ServerRepository) IncrementInviteUses(ctx context.Context, code string) error {
	return r.queries.IncrementInviteUses(ctx, code)
}

// DeleteInvite deletes an invite
func (r *ServerRepository) DeleteInvite(ctx context.Context, code string) error {
	return r.queries.SoftDeleteInvite(ctx, code)
}

// CreateBan creates a ban
func (r *ServerRepository) CreateBan(ctx context.Context, serverID, userID, moderatorID int32, reason *string, expiresAt *pgtype.Timestamp) (repo.Ban, error) {
	var reasonType pgtype.Text
	var expiresAtType pgtype.Timestamp

	if reason != nil {
		reasonType = pgtype.Text{String: *reason, Valid: true}
	}
	if expiresAt != nil {
		expiresAtType = *expiresAt
	}

	return r.queries.CreateBan(ctx, repo.CreateBanParams{
		ServerID:    serverID,
		UserID:      userID,
		ModeratorID: moderatorID,
		Reason:      reasonType,
		ExpiresAt:   expiresAtType,
	})
} // GetBan retrieves a ban
func (r *ServerRepository) GetBan(ctx context.Context, serverID, userID int32) (repo.Ban, error) {
	return r.queries.GetBan(ctx, repo.GetBanParams{
		ServerID: serverID,
		UserID:   userID,
	})
}

// GetServerBans retrieves all bans for a server
func (r *ServerRepository) GetServerBans(ctx context.Context, serverID int32) ([]repo.Ban, error) {
	return r.queries.GetServerBans(ctx, serverID)
}

// DeleteBan removes a ban
func (r *ServerRepository) DeleteBan(ctx context.Context, serverID, userID int32) error {
	return r.queries.SoftDeleteBan(ctx, repo.SoftDeleteBanParams{
		ServerID: serverID,
		UserID:   userID,
	})
}

// IsUserBanned checks if a user is banned
func (r *ServerRepository) IsUserBanned(ctx context.Context, serverID, userID int32) (bool, error) {
	return r.queries.IsUserBanned(ctx, repo.IsUserBannedParams{
		ServerID: serverID,
		UserID:   userID,
	})
}
