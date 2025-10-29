package util

// Permission flags as constants
const (
	PermissionCreateInvite         int64 = 1 << 0  // 0x1
	PermissionKickMembers          int64 = 1 << 1  // 0x2
	PermissionBanMembers           int64 = 1 << 2  // 0x4
	PermissionAdministrator        int64 = 1 << 3  // 0x8
	PermissionManageChannels       int64 = 1 << 4  // 0x10
	PermissionManageServer         int64 = 1 << 5  // 0x20
	PermissionAddReactions         int64 = 1 << 6  // 0x40
	PermissionViewAuditLog         int64 = 1 << 7  // 0x80
	PermissionPrioritySpeaker      int64 = 1 << 8  // 0x100
	PermissionStream               int64 = 1 << 9  // 0x200
	PermissionViewChannel          int64 = 1 << 10 // 0x400
	PermissionSendMessages         int64 = 1 << 11 // 0x800
	PermissionSendTTSMessages      int64 = 1 << 12 // 0x1000
	PermissionManageMessages       int64 = 1 << 13 // 0x2000
	PermissionEmbedLinks           int64 = 1 << 14 // 0x4000
	PermissionAttachFiles          int64 = 1 << 15 // 0x8000
	PermissionReadMessageHistory   int64 = 1 << 16 // 0x10000
	PermissionMentionEveryone      int64 = 1 << 17 // 0x20000
	PermissionUseExternalEmojis    int64 = 1 << 18 // 0x40000
	PermissionViewServerInsights   int64 = 1 << 19 // 0x80000
	PermissionConnect              int64 = 1 << 20 // 0x100000
	PermissionSpeak                int64 = 1 << 21 // 0x200000
	PermissionMuteMembers          int64 = 1 << 22 // 0x400000
	PermissionDeafenMembers        int64 = 1 << 23 // 0x800000
	PermissionMoveMembers          int64 = 1 << 24 // 0x1000000
	PermissionUseVAD               int64 = 1 << 25 // 0x2000000
	PermissionChangeNickname       int64 = 1 << 26 // 0x4000000
	PermissionManageNicknames      int64 = 1 << 27 // 0x8000000
	PermissionManageRoles          int64 = 1 << 28 // 0x10000000
	PermissionManageWebhooks       int64 = 1 << 29 // 0x20000000
	PermissionManageEmojisStickers int64 = 1 << 30 // 0x40000000
)

// HasPermission checks if the permission bits contain a specific permission
func HasPermission(permissions, permission int64) bool {
	return (permissions & permission) == permission
}

// AddPermission adds a permission to the permission bits
func AddPermission(permissions, permission int64) int64 {
	return permissions | permission
}

// RemovePermission removes a permission from the permission bits
func RemovePermission(permissions, permission int64) int64 {
	return permissions &^ permission
}

// CalculatePermissions calculates effective permissions considering allow and deny overwrites
func CalculatePermissions(basePermissions, allowOverwrites, denyOverwrites int64) int64 {
	// Apply deny overwrites first
	permissions := basePermissions &^ denyOverwrites
	// Then apply allow overwrites
	permissions = permissions | allowOverwrites
	return permissions
}

// CalculateChannelPermissions calculates permissions for a user in a channel
// considering role permissions and channel overwrites
func CalculateChannelPermissions(basePermissions int64, roleOverwrites map[int32]*ChannelOverwrite, userOverwrite *ChannelOverwrite) int64 {
	permissions := basePermissions

	// Apply role overwrites
	for _, overwrite := range roleOverwrites {
		permissions = CalculatePermissions(permissions, overwrite.Allow, overwrite.Deny)
	}

	// Apply user overwrite (highest priority)
	if userOverwrite != nil {
		permissions = CalculatePermissions(permissions, userOverwrite.Allow, userOverwrite.Deny)
	}

	// Administrator has all permissions
	if HasPermission(basePermissions, PermissionAdministrator) {
		return ^int64(0) // All bits set to 1
	}

	return permissions
}

// ChannelOverwrite represents permission overwrite for a channel
type ChannelOverwrite struct {
	ID    int32
	Type  string // "role" or "member"
	Allow int64
	Deny  int64
}

// CanViewChannel checks if user can view channel
func CanViewChannel(permissions int64) bool {
	return HasPermission(permissions, PermissionViewChannel)
}

// CanSendMessages checks if user can send messages
func CanSendMessages(permissions int64) bool {
	return HasPermission(permissions, PermissionSendMessages)
}

// CanManageChannel checks if user can manage channel
func CanManageChannel(permissions int64) bool {
	return HasPermission(permissions, PermissionManageChannels)
}

// CanManageMessages checks if user can manage messages
func CanManageMessages(permissions int64) bool {
	return HasPermission(permissions, PermissionManageMessages)
}

// CanConnect checks if user can connect to voice channel
func CanConnect(permissions int64) bool {
	return HasPermission(permissions, PermissionConnect)
}

// CanSpeak checks if user can speak in voice channel
func CanSpeak(permissions int64) bool {
	return HasPermission(permissions, PermissionSpeak)
}

// IsAdministrator checks if user has administrator permission
func IsAdministrator(permissions int64) bool {
	return HasPermission(permissions, PermissionAdministrator)
}

// GetPermissionNames returns human-readable names for permission bits
func GetPermissionNames(permissions int64) []string {
	names := []string{}

	permissionMap := map[int64]string{
		PermissionCreateInvite:         "Create Invite",
		PermissionKickMembers:          "Kick Members",
		PermissionBanMembers:           "Ban Members",
		PermissionAdministrator:        "Administrator",
		PermissionManageChannels:       "Manage Channels",
		PermissionManageServer:         "Manage Server",
		PermissionAddReactions:         "Add Reactions",
		PermissionViewAuditLog:         "View Audit Log",
		PermissionPrioritySpeaker:      "Priority Speaker",
		PermissionStream:               "Stream",
		PermissionViewChannel:          "View Channel",
		PermissionSendMessages:         "Send Messages",
		PermissionSendTTSMessages:      "Send TTS Messages",
		PermissionManageMessages:       "Manage Messages",
		PermissionEmbedLinks:           "Embed Links",
		PermissionAttachFiles:          "Attach Files",
		PermissionReadMessageHistory:   "Read Message History",
		PermissionMentionEveryone:      "Mention Everyone",
		PermissionUseExternalEmojis:    "Use External Emojis",
		PermissionViewServerInsights:   "View Server Insights",
		PermissionConnect:              "Connect",
		PermissionSpeak:                "Speak",
		PermissionMuteMembers:          "Mute Members",
		PermissionDeafenMembers:        "Deafen Members",
		PermissionMoveMembers:          "Move Members",
		PermissionUseVAD:               "Use Voice Activity",
		PermissionChangeNickname:       "Change Nickname",
		PermissionManageNicknames:      "Manage Nicknames",
		PermissionManageRoles:          "Manage Roles",
		PermissionManageWebhooks:       "Manage Webhooks",
		PermissionManageEmojisStickers: "Manage Emojis and Stickers",
	}

	for perm, name := range permissionMap {
		if HasPermission(permissions, perm) {
			names = append(names, name)
		}
	}

	return names
}
