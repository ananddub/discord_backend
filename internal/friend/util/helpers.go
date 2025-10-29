package util

import (
	"discord/gen/proto/schema"
	"discord/gen/repo"
)

// ConvertFriendToProto converts a repo.Friend to proto.Friend
func ConvertFriendToProto(friend repo.Friend) *schema.Friend {
	pbFriend := &schema.Friend{
		Id:         friend.ID,
		UserId:     friend.UserID,
		FriendId:   friend.FriendID,
		Status:     friend.Status,
		IsAccepted: friend.Status == "accepted",
		IsBlocked:  friend.Status == "blocked",
		IsPending:  friend.Status == "pending",
		IsRejected: friend.Status == "rejected",
		IsFavorite: friend.IsFavorite.Bool,
		CreatedAt:  friend.CreatedAt.Time.Unix(),
		UpdatedAt:  friend.UpdatedAt.Time.Unix(),
	}

	if friend.AliasName.Valid {
		pbFriend.AliasName = &friend.AliasName.String
	}

	return pbFriend
}

// ValidateFriendshipStatus checks if the friendship status is valid
func ValidateFriendshipStatus(status string) bool {
	validStatuses := []string{"pending", "accepted", "blocked", "rejected"}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsMutualFriend checks if two users are mutual friends
func IsMutualFriend(user1Friends, user2Friends []repo.Friend) bool {
	friendMap := make(map[int32]bool)

	for _, friend := range user1Friends {
		if friend.Status == "accepted" {
			friendMap[friend.FriendID] = true
		}
	}

	for _, friend := range user2Friends {
		if friend.Status == "accepted" && friendMap[friend.FriendID] {
			return true
		}
	}

	return false
}

// GetMutualFriends returns the list of mutual friends between two users
func GetMutualFriends(user1Friends, user2Friends []repo.Friend) []int32 {
	friendMap := make(map[int32]bool)
	var mutualFriends []int32

	for _, friend := range user1Friends {
		if friend.Status == "accepted" {
			friendMap[friend.FriendID] = true
		}
	}

	for _, friend := range user2Friends {
		if friend.Status == "accepted" && friendMap[friend.FriendID] {
			mutualFriends = append(mutualFriends, friend.FriendID)
		}
	}

	return mutualFriends
}

// FilterFriendsByStatus filters friends by their status
func FilterFriendsByStatus(friends []repo.Friend, status string) []repo.Friend {
	var filtered []repo.Friend
	for _, friend := range friends {
		if friend.Status == status {
			filtered = append(filtered, friend)
		}
	}
	return filtered
}

// GetFavoriteFriends returns only favorite friends
func GetFavoriteFriends(friends []repo.Friend) []repo.Friend {
	var favorites []repo.Friend
	for _, friend := range friends {
		if friend.IsFavorite.Bool && friend.Status == "accepted" {
			favorites = append(favorites, friend)
		}
	}
	return favorites
}
