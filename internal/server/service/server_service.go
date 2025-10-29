package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"discord/gen/repo"
	commonErrors "discord/internal/common/errors"
	serverRepo "discord/internal/server/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type ServerService struct {
	serverRepo *serverRepo.ServerRepository
}

func NewServerService(serverRepo *serverRepo.ServerRepository) *ServerService {
	return &ServerService{
		serverRepo: serverRepo,
	}
}

// CreateServer creates a new server
func (s *ServerService) CreateServer(ctx context.Context, name string, ownerID int32, icon, banner, description, region *string) (repo.Server, error) {
	if name == "" {
		return repo.Server{}, commonErrors.ErrInvalidInput
	}

	server, err := s.serverRepo.CreateServer(ctx, name, ownerID, icon, banner, description, region)
	if err != nil {
		return repo.Server{}, err
	}

	// Add owner as first member
	_, err = s.serverRepo.AddServerMember(ctx, server.ID, ownerID, nil)
	if err != nil {
		return repo.Server{}, err
	}

	// Increment member count
	_ = s.serverRepo.IncrementMemberCount(ctx, server.ID)

	return server, nil
}

// GetServer retrieves a server by ID
func (s *ServerService) GetServer(ctx context.Context, serverID int32) (repo.Server, error) {
	return s.serverRepo.GetServerByID(ctx, serverID)
}

// UpdateServer updates server information
func (s *ServerService) UpdateServer(ctx context.Context, serverID, userID int32, name, icon, banner, description, region *string) (repo.Server, error) {
	// Verify user is server owner
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return repo.Server{}, commonErrors.ErrNotFound
	}

	if server.OwnerID != userID {
		return repo.Server{}, commonErrors.ErrPermissionDenied
	}

	return s.serverRepo.UpdateServer(ctx, serverID, name, icon, banner, description, region)
}

// DeleteServer deletes a server
func (s *ServerService) DeleteServer(ctx context.Context, serverID, userID int32) error {
	// Verify user is server owner
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if server.OwnerID != userID {
		return commonErrors.ErrPermissionDenied
	}

	return s.serverRepo.DeleteServer(ctx, serverID)
}

// GetUserServers retrieves all servers a user is in
func (s *ServerService) GetUserServers(ctx context.Context, userID int32) ([]repo.Server, error) {
	return s.serverRepo.GetUserServers(ctx, userID)
}

// TransferOwnership transfers server ownership
func (s *ServerService) TransferOwnership(ctx context.Context, serverID, currentOwnerID, newOwnerID int32) error {
	// Verify current owner
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if server.OwnerID != currentOwnerID {
		return commonErrors.ErrPermissionDenied
	}

	// Verify new owner is a member
	_, err = s.serverRepo.GetServerMember(ctx, serverID, newOwnerID)
	if err != nil {
		return errors.New("new owner must be a server member")
	}

	return s.serverRepo.UpdateServerOwner(ctx, serverID, newOwnerID)
}

// JoinServer adds a user to a server
func (s *ServerService) JoinServer(ctx context.Context, serverID, userID int32, inviteCode *string) error {
	// Check if user is banned
	banned, err := s.serverRepo.IsUserBanned(ctx, serverID, userID)
	if err == nil && banned {
		return errors.New("user is banned from this server")
	}

	// Check if already a member
	_, err = s.serverRepo.GetServerMember(ctx, serverID, userID)
	if err == nil {
		return errors.New("already a member of this server")
	}

	// Validate invite if provided
	if inviteCode != nil {
		invite, err := s.serverRepo.GetInviteByCode(ctx, *inviteCode)
		if err != nil {
			return errors.New("invalid invite code")
		}

		if invite.ServerID != serverID {
			return errors.New("invite is for a different server")
		}

		// Check if invite has expired
		if invite.ExpiresAt.Valid && invite.ExpiresAt.Time.Before(time.Now()) {
			return errors.New("invite has expired")
		}

		// Check max uses
		if invite.MaxUses.Valid && invite.MaxUses.Int32 > 0 && invite.Uses.Int32 >= invite.MaxUses.Int32 {
			return errors.New("invite has reached max uses")
		}

		// Increment invite uses
		_ = s.serverRepo.IncrementInviteUses(ctx, *inviteCode)
	}

	// Add member
	_, err = s.serverRepo.AddServerMember(ctx, serverID, userID, nil)
	if err != nil {
		return err
	}

	// Increment member count
	return s.serverRepo.IncrementMemberCount(ctx, serverID)
}

// LeaveServer removes a user from a server
func (s *ServerService) LeaveServer(ctx context.Context, serverID, userID int32) error {
	// Check if user is owner
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if server.OwnerID == userID {
		return errors.New("owner cannot leave server, transfer ownership first")
	}

	// Remove member
	err = s.serverRepo.RemoveServerMember(ctx, serverID, userID)
	if err != nil {
		return err
	}

	// Decrement member count
	return s.serverRepo.DecrementMemberCount(ctx, serverID)
}

// KickMember kicks a member from the server
func (s *ServerService) KickMember(ctx context.Context, serverID, moderatorID, targetUserID int32) error {
	// Verify server exists and moderator has permission
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Owner can kick anyone
	if server.OwnerID != moderatorID {
		// TODO: Check if moderator has KICK_MEMBERS permission
		return commonErrors.ErrPermissionDenied
	}

	// Cannot kick owner
	if server.OwnerID == targetUserID {
		return errors.New("cannot kick server owner")
	}

	// Remove member
	err = s.serverRepo.RemoveServerMember(ctx, serverID, targetUserID)
	if err != nil {
		return err
	}

	// Decrement member count
	return s.serverRepo.DecrementMemberCount(ctx, serverID)
}

// GetServerMembers retrieves server members
func (s *ServerService) GetServerMembers(ctx context.Context, serverID int32, limit, offset int32) ([]repo.ServerMember, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	return s.serverRepo.GetServerMembers(ctx, serverID, limit, offset)
}

// UpdateMemberNickname updates a member's nickname
func (s *ServerService) UpdateMemberNickname(ctx context.Context, serverID, userID int32, nickname *string) error {
	// Verify member exists
	_, err := s.serverRepo.GetServerMember(ctx, serverID, userID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	return s.serverRepo.UpdateMemberNickname(ctx, serverID, userID, nickname)
}

// CreateRole creates a new role
func (s *ServerService) CreateRole(ctx context.Context, serverID, userID int32, name string, color *string, hoist, mentionable bool, permissions int64, description *string) (repo.Role, error) {
	// Verify user has permission
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return repo.Role{}, commonErrors.ErrNotFound
	}

	if server.OwnerID != userID {
		// TODO: Check MANAGE_ROLES permission
		return repo.Role{}, commonErrors.ErrPermissionDenied
	}

	// Get current role count for position
	roles, _ := s.serverRepo.GetServerRoles(ctx, serverID)
	position := int32(len(roles))

	return s.serverRepo.CreateRole(ctx, serverID, name, color, hoist, mentionable, position, permissions, description)
}

// GetServerRoles retrieves all roles for a server
func (s *ServerService) GetServerRoles(ctx context.Context, serverID int32) ([]repo.Role, error) {
	return s.serverRepo.GetServerRoles(ctx, serverID)
}

// DeleteRole deletes a role
func (s *ServerService) DeleteRole(ctx context.Context, roleID, userID int32) error {
	role, err := s.serverRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Verify user has permission
	server, err := s.serverRepo.GetServerByID(ctx, role.ServerID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if server.OwnerID != userID {
		// TODO: Check MANAGE_ROLES permission
		return commonErrors.ErrPermissionDenied
	}

	return s.serverRepo.DeleteRole(ctx, roleID)
}

// CreateInvite creates a server invite
func (s *ServerService) CreateInvite(ctx context.Context, serverID, channelID, inviterID int32, maxUses, maxAge int32, temporary bool) (repo.Invite, error) {
	// Generate random invite code
	code := generateInviteCode()

	// Calculate expiration if maxAge is set
	var expiresAt *pgtype.Timestamp
	if maxAge > 0 {
		expTime := time.Now().Add(time.Duration(maxAge) * time.Second)
		expiresAt = &pgtype.Timestamp{
			Time:  expTime,
			Valid: true,
		}
	}

	return s.serverRepo.CreateInvite(ctx, code, serverID, channelID, inviterID, maxUses, maxAge, temporary, expiresAt)
}

// GetInvite retrieves an invite by code
func (s *ServerService) GetInvite(ctx context.Context, code string) (repo.Invite, error) {
	return s.serverRepo.GetInviteByCode(ctx, code)
}

// GetServerInvites retrieves all invites for a server
func (s *ServerService) GetServerInvites(ctx context.Context, serverID int32) ([]repo.Invite, error) {
	return s.serverRepo.GetServerInvites(ctx, serverID)
}

// DeleteInvite deletes an invite
func (s *ServerService) DeleteInvite(ctx context.Context, code string, userID int32) error {
	invite, err := s.serverRepo.GetInviteByCode(ctx, code)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	// Verify user has permission (owner or invite creator)
	server, err := s.serverRepo.GetServerByID(ctx, invite.ServerID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if server.OwnerID != userID && invite.InviterID != userID {
		// TODO: Check MANAGE_SERVER permission
		return commonErrors.ErrPermissionDenied
	}

	return s.serverRepo.DeleteInvite(ctx, code)
}

// BanMember bans a user from the server
func (s *ServerService) BanMember(ctx context.Context, serverID, moderatorID, targetUserID int32, reason *string, durationSeconds *int32) error {
	// Verify moderator has permission
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if server.OwnerID != moderatorID {
		// TODO: Check BAN_MEMBERS permission
		return commonErrors.ErrPermissionDenied
	}

	// Cannot ban owner
	if server.OwnerID == targetUserID {
		return errors.New("cannot ban server owner")
	}

	// Calculate expiration if duration is set
	var expiresAt *pgtype.Timestamp
	if durationSeconds != nil && *durationSeconds > 0 {
		expTime := time.Now().Add(time.Duration(*durationSeconds) * time.Second)
		expiresAt = &pgtype.Timestamp{
			Time:  expTime,
			Valid: true,
		}
	}

	// Create ban
	_, err = s.serverRepo.CreateBan(ctx, serverID, targetUserID, moderatorID, reason, expiresAt)
	if err != nil {
		return err
	}

	// Remove member if they're in the server
	_ = s.serverRepo.RemoveServerMember(ctx, serverID, targetUserID)
	_ = s.serverRepo.DecrementMemberCount(ctx, serverID)

	return nil
}

// UnbanMember unbans a user from the server
func (s *ServerService) UnbanMember(ctx context.Context, serverID, moderatorID, targetUserID int32) error {
	// Verify moderator has permission
	server, err := s.serverRepo.GetServerByID(ctx, serverID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if server.OwnerID != moderatorID {
		// TODO: Check BAN_MEMBERS permission
		return commonErrors.ErrPermissionDenied
	}

	return s.serverRepo.DeleteBan(ctx, serverID, targetUserID)
}

// GetServerBans retrieves all bans for a server
func (s *ServerService) GetServerBans(ctx context.Context, serverID int32) ([]repo.Ban, error) {
	return s.serverRepo.GetServerBans(ctx, serverID)
}

// generateInviteCode generates a random invite code
func generateInviteCode() string {
	b := make([]byte, 8)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:8]
}
