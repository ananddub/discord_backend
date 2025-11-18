package reactive

import (
	"discord/pkg/mypg"
	"strconv"
	"strings"
)

func ReactiveEvents(pg mypg.QueryData) {
	pgtype := strings.ToLower(pg.Table)
	switch pgtype {
	case "users":
		UserReactive{}.ReactUser(pg)
	case "friends":
		FriendReactive{}.ReactFriend(pg)
	default:
		// handle unknown table changes
	}
}

func atoi32(s string) int32 {
	val, _ := strconv.ParseInt(s, 10, 32)
	return int32(val)
}

func atoi64(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

func atoiBool(s string) bool {
	return s == "true" || s == "1"
}
