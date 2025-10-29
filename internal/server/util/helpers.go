package util

import (
	"discord/gen/repo"
	"fmt"
	"strings"
)

// ValidateServerName validates server name
func ValidateServerName(name string) bool {
	// Check if not empty
	if name == "" {
		return false
	}

	// Check length (2-100 characters)
	if len(name) < 2 || len(name) > 100 {
		return false
	}

	return true
}

// ValidateRoleName validates role name
func ValidateRoleName(name string) bool {
	if name == "" {
		return false
	}

	if len(name) > 100 {
		return false
	}

	return true
}

// IsServerOwner checks if user is server owner
func IsServerOwner(server repo.Server, userID int32) bool {
	return server.OwnerID == userID
}

// IsServerMember checks if user is in the server members list
func IsServerMember(members []repo.ServerMember, userID int32) bool {
	for _, member := range members {
		if member.UserID == userID {
			return true
		}
	}
	return false
}

// HasPermission checks if user has a specific permission
func HasPermission(permissions int64, permission int64) bool {
	return permissions&permission == permission
}

// AddPermission adds a permission to the permissions bitfield
func AddPermission(permissions int64, permission int64) int64 {
	return permissions | permission
}

// RemovePermission removes a permission from the permissions bitfield
func RemovePermission(permissions int64, permission int64) int64 {
	return permissions &^ permission
}

// Permission constants (bitfield)
const (
	PermissionCreateInvite       int64 = 1 << 0  // 1
	PermissionKickMembers        int64 = 1 << 1  // 2
	PermissionBanMembers         int64 = 1 << 2  // 4
	PermissionAdministrator      int64 = 1 << 3  // 8
	PermissionManageChannels     int64 = 1 << 4  // 16
	PermissionManageServer       int64 = 1 << 5  // 32
	PermissionAddReactions       int64 = 1 << 6  // 64
	PermissionViewAuditLog       int64 = 1 << 7  // 128
	PermissionViewChannel        int64 = 1 << 10 // 1024
	PermissionSendMessages       int64 = 1 << 11 // 2048
	PermissionManageMessages     int64 = 1 << 13 // 8192
	PermissionEmbedLinks         int64 = 1 << 14 // 16384
	PermissionAttachFiles        int64 = 1 << 15 // 32768
	PermissionReadMessageHistory int64 = 1 << 16 // 65536
	PermissionMentionEveryone    int64 = 1 << 17 // 131072
	PermissionConnect            int64 = 1 << 20 // 1048576
	PermissionSpeak              int64 = 1 << 21 // 2097152
	PermissionMuteMembers        int64 = 1 << 22 // 4194304
	PermissionDeafenMembers      int64 = 1 << 23 // 8388608
	PermissionMoveMembers        int64 = 1 << 24 // 16777216
	PermissionManageRoles        int64 = 1 << 28 // 268435456
)

// GetDefaultPermissions returns default permissions for new members
func GetDefaultPermissions() int64 {
	return PermissionViewChannel |
		PermissionSendMessages |
		PermissionReadMessageHistory |
		PermissionAddReactions |
		PermissionConnect |
		PermissionSpeak |
		PermissionEmbedLinks |
		PermissionAttachFiles
}

// GetAdministratorPermissions returns all permissions
func GetAdministratorPermissions() int64 {
	return PermissionAdministrator | (1<<30 - 1) // All permissions
}

// FormatServerInfo formats server information for display
func FormatServerInfo(server repo.Server) string {
	memberCount := int32(0)
	if server.MemberCount.Valid {
		memberCount = server.MemberCount.Int32
	}
	return fmt.Sprintf("Server: %s (ID: %d, Owner: %d, Members: %d)",
		server.Name, server.ID, server.OwnerID, memberCount)
}

// FormatRoleInfo formats role information
func FormatRoleInfo(role repo.Role) string {
	position := int32(0)
	permissions := int64(0)
	if role.Position.Valid {
		position = role.Position.Int32
	}
	if role.Permissions.Valid {
		permissions = role.Permissions.Int64
	}
	return fmt.Sprintf("Role: %s (ID: %d, Position: %d, Permissions: %d)",
		role.Name, role.ID, position, permissions)
}

// FormatInviteCode formats an invite code for display
func FormatInviteCode(code string, serverName string) string {
	return fmt.Sprintf("Invite: %s for %s", code, serverName)
}

// ParseInviteCode extracts the code from various invite formats
func ParseInviteCode(input string) string {
	// Remove common URL prefixes
	input = strings.TrimPrefix(input, "https://discord.gg/")
	input = strings.TrimPrefix(input, "http://discord.gg/")
	input = strings.TrimPrefix(input, "discord.gg/")

	return strings.TrimSpace(input)
}

// ValidateInviteCode validates invite code format
func ValidateInviteCode(code string) bool {
	if len(code) < 4 || len(code) > 32 {
		return false
	}

	// Basic alphanumeric check
	for _, c := range code {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}

	return true
}

// CalculateRolePosition calculates the position for a new role
func CalculateRolePosition(roles []repo.Role) int32 {
	if len(roles) == 0 {
		return 0
	}

	maxPos := int32(0)
	for _, role := range roles {
		if role.Position.Valid && role.Position.Int32 > maxPos {
			maxPos = role.Position.Int32
		}
	}

	return maxPos + 1
}

// SortRolesByPosition sorts roles by their position (descending)
func SortRolesByPosition(roles []repo.Role) []repo.Role {
	// Simple bubble sort - for production, use sort.Slice
	n := len(roles)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			posJ := int32(0)
			posJ1 := int32(0)
			if roles[j].Position.Valid {
				posJ = roles[j].Position.Int32
			}
			if roles[j+1].Position.Valid {
				posJ1 = roles[j+1].Position.Int32
			}
			if posJ < posJ1 {
				roles[j], roles[j+1] = roles[j+1], roles[j]
			}
		}
	}
	return roles
}

// GetMemberHighestRole gets the highest role of a member
func GetMemberHighestRole(memberRoles []repo.Role) *repo.Role {
	if len(memberRoles) == 0 {
		return nil
	}

	highest := &memberRoles[0]
	highestPos := int32(0)
	if highest.Position.Valid {
		highestPos = highest.Position.Int32
	}

	for i := 1; i < len(memberRoles); i++ {
		if memberRoles[i].Position.Valid && memberRoles[i].Position.Int32 > highestPos {
			highest = &memberRoles[i]
			highestPos = memberRoles[i].Position.Int32
		}
	}

	return highest
}

// CanManageRole checks if a user can manage a role based on position hierarchy
func CanManageRole(userHighestRolePosition, targetRolePosition int32) bool {
	return userHighestRolePosition > targetRolePosition
}

// GetPermissionNames returns human-readable names for permissions
func GetPermissionNames(permissions int64) []string {
	names := []string{}

	permissionMap := map[int64]string{
		PermissionAdministrator:      "Administrator",
		PermissionManageServer:       "Manage Server",
		PermissionManageChannels:     "Manage Channels",
		PermissionManageRoles:        "Manage Roles",
		PermissionKickMembers:        "Kick Members",
		PermissionBanMembers:         "Ban Members",
		PermissionCreateInvite:       "Create Invite",
		PermissionViewChannel:        "View Channel",
		PermissionSendMessages:       "Send Messages",
		PermissionManageMessages:     "Manage Messages",
		PermissionReadMessageHistory: "Read Message History",
		PermissionMentionEveryone:    "Mention Everyone",
		PermissionAddReactions:       "Add Reactions",
		PermissionConnect:            "Connect to Voice",
		PermissionSpeak:              "Speak in Voice",
		PermissionMuteMembers:        "Mute Members",
		PermissionDeafenMembers:      "Deafen Members",
	}

	for perm, name := range permissionMap {
		if HasPermission(permissions, perm) {
			names = append(names, name)
		}
	}

	return names
}

// IsValidColor checks if a color value is valid (0x000000 to 0xFFFFFF)
func IsValidColor(color int32) bool {
	return color >= 0 && color <= 0xFFFFFF
}
