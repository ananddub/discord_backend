package reactive

import (
	"discord/pkg/mypg"
	"strings"
)

func ReactiveEvents(pg mypg.QueryData) {
	pgtype := strings.ToLower(pg.Table)
	switch pgtype {
	case "users":
		UserReactive{}.ReactUser(pg)
	case "servers":
		ServerReactive{}.ReactServer(pg)
	case "roles":
		RoleReactive{}.ReactRole(pg)
	case "server_members":
		ServerMemberReactive{}.ReactServerMember(pg)
	case "member_roles":
		MemberRoleReactive{}.ReactMemberRole(pg)
	case "channels":
		ChannelReactive{}.ReactChannel(pg)
	case "channel_permissions":
		ChannelPermissionReactive{}.ReactChannelPermission(pg)
	case "messages":
		MessageRective{}.ReactMessage(pg)
	case "message_attachments":
		MessageAttachmentReactive{}.ReactMessageAttachment(pg)
	case "message_reactions":
		MessageReactionReactive{}.ReactMessageReaction(pg)
	case "friends":
		FriendReactive{}.ReactFriend(pg)
	default:
		// handle unknown table changes
	}
}
