package service

import (
	"context"
	"errors"
	"time"

	commonErrors "discord/internal/common/errors"
	dmRepo "discord/internal/dm/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type DMService struct {
	dmRepo *dmRepo.DMRepository
}

func NewDMService(dmRepo *dmRepo.DMRepository) *DMService {
	return &DMService{
		dmRepo: dmRepo,
	}
}

// DMChannelInfo represents DM channel information for response
type DMChannelInfo struct {
	ID             int32
	ParticipantIDs []int32
	Name           string
	Icon           string
	LastMessageID  int32
	LastMessageAt  int64
	UnreadCount    int32
	IsGroup        bool
}

// CreateDMChannel creates a 1-on-1 DM channel between two users
func (s *DMService) CreateDMChannel(ctx context.Context, userID, recipientID int32) (int32, error) {
	// Check if DM channel already exists between these users
	existingChannel, err := s.dmRepo.GetDMChannelForUsers(ctx, userID, recipientID)
	if err == nil && existingChannel != nil {
		return existingChannel.ID, nil
	}

	// Create new DM channel
	channel, err := s.dmRepo.CreateDMChannel(ctx, "", "", nil, false)
	if err != nil {
		return 0, err
	}

	// Add both users as participants
	_, err = s.dmRepo.AddDMParticipant(ctx, channel.ID, userID)
	if err != nil {
		s.dmRepo.DeleteDMChannel(ctx, channel.ID)
		return 0, err
	}

	_, err = s.dmRepo.AddDMParticipant(ctx, channel.ID, recipientID)
	if err != nil {
		s.dmRepo.DeleteDMChannel(ctx, channel.ID)
		return 0, err
	}

	return channel.ID, nil
}

// CreateGroupDM creates a group DM channel
func (s *DMService) CreateGroupDM(ctx context.Context, ownerID int32, userIDs []int32, name string) (int32, error) {
	if len(userIDs) < 2 {
		return 0, errors.New("group DM requires at least 2 users")
	}

	// Create group DM channel
	channel, err := s.dmRepo.CreateDMChannel(ctx, name, "", &ownerID, true)
	if err != nil {
		return 0, err
	}

	// Add owner as participant
	_, err = s.dmRepo.AddDMParticipant(ctx, channel.ID, ownerID)
	if err != nil {
		s.dmRepo.DeleteDMChannel(ctx, channel.ID)
		return 0, err
	}

	// Add all users as participants
	for _, userID := range userIDs {
		if userID == ownerID {
			continue // Skip owner, already added
		}
		_, err = s.dmRepo.AddDMParticipant(ctx, channel.ID, userID)
		if err != nil {
			// Continue adding others even if one fails
			continue
		}
	}

	return channel.ID, nil
}

// GetDMChannel retrieves DM channel with participants
func (s *DMService) GetDMChannel(ctx context.Context, dmChannelID int32) (*DMChannelInfo, error) {
	channel, err := s.dmRepo.GetDMChannelByID(ctx, dmChannelID)
	if err != nil {
		return nil, commonErrors.ErrNotFound
	}

	participants, err := s.dmRepo.GetDMParticipants(ctx, dmChannelID)
	if err != nil {
		return nil, err
	}

	participantIDs := make([]int32, len(participants))
	for i, p := range participants {
		participantIDs[i] = p.UserID
	}

	info := &DMChannelInfo{
		ID:             channel.ID,
		ParticipantIDs: participantIDs,
	}

	if channel.Name.Valid {
		info.Name = channel.Name.String
	}

	if channel.Icon.Valid {
		info.Icon = channel.Icon.String
	}

	if channel.LastMessageID.Valid {
		info.LastMessageID = channel.LastMessageID.Int32
	}

	if channel.LastMessageAt.Valid {
		info.LastMessageAt = channel.LastMessageAt.Time.Unix()
	}

	if channel.IsGroup.Valid {
		info.IsGroup = channel.IsGroup.Bool
	}

	return info, nil
}

// GetUserDMChannels retrieves all DM channels for a user
func (s *DMService) GetUserDMChannels(ctx context.Context, userID int32) ([]*DMChannelInfo, error) {
	channels, err := s.dmRepo.GetUserDMChannels(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*DMChannelInfo, len(channels))
	for i, channel := range channels {
		participants, err := s.dmRepo.GetDMParticipants(ctx, channel.ID)
		if err != nil {
			continue
		}

		participantIDs := make([]int32, len(participants))
		for j, p := range participants {
			participantIDs[j] = p.UserID
		}

		info := &DMChannelInfo{
			ID:             channel.ID,
			ParticipantIDs: participantIDs,
		}

		if channel.Name.Valid {
			info.Name = channel.Name.String
		}

		if channel.Icon.Valid {
			info.Icon = channel.Icon.String
		}

		if channel.LastMessageID.Valid {
			info.LastMessageID = channel.LastMessageID.Int32
		}

		if channel.LastMessageAt.Valid {
			info.LastMessageAt = channel.LastMessageAt.Time.Unix()
		}

		if channel.IsGroup.Valid {
			info.IsGroup = channel.IsGroup.Bool
		}

		result[i] = info
	}

	return result, nil
}

// CloseDMChannel closes (deletes) a DM channel for a user
func (s *DMService) CloseDMChannel(ctx context.Context, dmChannelID, userID int32) error {
	// Remove user as participant
	return s.dmRepo.RemoveDMParticipant(ctx, dmChannelID, userID)
}

// AddUserToGroupDM adds a user to a group DM
func (s *DMService) AddUserToGroupDM(ctx context.Context, dmChannelID, userID int32) error {
	// Verify it's a group DM
	channel, err := s.dmRepo.GetDMChannelByID(ctx, dmChannelID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if !channel.IsGroup.Valid || !channel.IsGroup.Bool {
		return errors.New("not a group DM")
	}

	// Add participant
	_, err = s.dmRepo.AddDMParticipant(ctx, dmChannelID, userID)
	return err
}

// RemoveUserFromGroupDM removes a user from a group DM
func (s *DMService) RemoveUserFromGroupDM(ctx context.Context, dmChannelID, userID int32) error {
	// Verify it's a group DM
	channel, err := s.dmRepo.GetDMChannelByID(ctx, dmChannelID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if !channel.IsGroup.Valid || !channel.IsGroup.Bool {
		return errors.New("not a group DM")
	}

	// Remove participant
	return s.dmRepo.RemoveDMParticipant(ctx, dmChannelID, userID)
}

// UpdateGroupDM updates group DM information
func (s *DMService) UpdateGroupDM(ctx context.Context, dmChannelID int32, name, icon *string) error {
	// Verify it's a group DM
	channel, err := s.dmRepo.GetDMChannelByID(ctx, dmChannelID)
	if err != nil {
		return commonErrors.ErrNotFound
	}

	if !channel.IsGroup.Valid || !channel.IsGroup.Bool {
		return errors.New("not a group DM")
	}

	// Update channel
	_, err = s.dmRepo.UpdateDMChannel(ctx, dmChannelID, name, icon, nil, nil)
	return err
}

// MarkAsRead marks messages as read for a user
func (s *DMService) MarkAsRead(ctx context.Context, dmChannelID, userID, messageID int32) error {
	return s.dmRepo.UpdateLastReadMessage(ctx, dmChannelID, userID, messageID)
}

// UpdateLastMessage updates the last message info for a DM channel
func (s *DMService) UpdateLastMessage(ctx context.Context, dmChannelID, messageID int32) error {
	now := pgtype.Timestamp{Time: time.Now(), Valid: true}
	_, err := s.dmRepo.UpdateDMChannel(ctx, dmChannelID, nil, nil, &messageID, &now)
	return err
}
