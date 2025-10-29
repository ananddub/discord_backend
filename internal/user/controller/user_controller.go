package controller

import (
	"context"
	"time"

	"discord/gen/proto/schema"
	userPb "discord/gen/proto/service/user"
	commonErrors "discord/internal/common/errors"
	userService "discord/internal/user/service"
)

type UserController struct {
	userPb.UnimplementedUserServiceServer
	userService *userService.UserService
}

func NewUserController(userService *userService.UserService) *userPb.UserServiceServer {
	controller := &UserController{
		userService: userService,
	}
	var grpcController userPb.UserServiceServer = controller
	return &grpcController
}

// GetUser retrieves user information
func (c *UserController) GetUser(ctx context.Context, req *userPb.GetUserRequest) (*userPb.GetUserResponse, error) {
	if req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	user, err := c.userService.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbUser := &schema.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}

	// All optional string fields
	if user.FullName.Valid {
		pbUser.FullName = user.FullName.String
	}
	if user.ProfilePic.Valid {
		pbUser.ProfilePic = user.ProfilePic.String
	}
	if user.Bio.Valid {
		pbUser.Bio = user.Bio.String
	}
	if user.ColorCode.Valid {
		pbUser.ColorCode = user.ColorCode.String
	}
	if user.BackgroundColor.Valid {
		pbUser.BackgroundColor = user.BackgroundColor.String
	}
	if user.BackgroundPic.Valid {
		pbUser.BackgroundPic = user.BackgroundPic.String
	}
	if user.CustomStatus.Valid {
		pbUser.CustomStatus = user.CustomStatus.String
	}

	// Boolean fields
	if user.IsBot.Valid {
		pbUser.IsBot = user.IsBot.Bool
	}
	if user.IsVerified.Valid {
		pbUser.IsVerified = user.IsVerified.Bool
	}

	// Timestamps
	if user.CreatedAt.Valid {
		pbUser.CreatedAt = user.CreatedAt.Time.Unix()
	}
	if user.UpdatedAt.Valid {
		pbUser.UpdatedAt = user.UpdatedAt.Time.Unix()
	}

	return &userPb.GetUserResponse{
		User:    pbUser,
		Success: true,
	}, nil
}

// UpdateUser updates user profile information
func (c *UserController) UpdateUser(ctx context.Context, req *userPb.UpdateUserRequest) (*userPb.UpdateUserResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetUser() == nil {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	var fullName, profilePic, backgroundColor, colorCode, backgroundPic, bio *string
	if req.User.GetFullName() != "" {
		fn := req.User.GetFullName()
		fullName = &fn
	}
	if req.User.GetProfilePic() != "" {
		pp := req.User.GetProfilePic()
		profilePic = &pp
	}
	if req.User.GetBio() != "" {
		b := req.User.GetBio()
		bio = &b
	}
	if req.User.GetColorCode() == "" {
		cc := req.User.GetColorCode()
		colorCode = &cc
	}
	if req.User.GetBackgroundColor() == "" {
		bc := req.User.GetBackgroundColor()
		backgroundColor = &bc
	}
	if req.User.GetBackgroundPic() == "" {
		bp := req.User.GetBackgroundPic()
		backgroundPic = &bp
	}
	user, err := c.userService.UpdateUser(ctx, userID, fullName,
		profilePic,
		bio,
		colorCode,
		backgroundColor,
		backgroundPic,
	)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbUser := &schema.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}

	// All optional string fields
	if user.FullName.Valid {
		pbUser.FullName = user.FullName.String
	}
	if user.ProfilePic.Valid {
		pbUser.ProfilePic = user.ProfilePic.String
	}
	if user.Bio.Valid {
		pbUser.Bio = user.Bio.String
	}
	if user.ColorCode.Valid {
		pbUser.ColorCode = user.ColorCode.String
	}
	if user.BackgroundColor.Valid {
		pbUser.BackgroundColor = user.BackgroundColor.String
	}
	if user.BackgroundPic.Valid {
		pbUser.BackgroundPic = user.BackgroundPic.String
	}
	if user.CustomStatus.Valid {
		pbUser.CustomStatus = user.CustomStatus.String
	}

	// Boolean fields
	if user.IsBot.Valid {
		pbUser.IsBot = user.IsBot.Bool
	}
	if user.IsVerified.Valid {
		pbUser.IsVerified = user.IsVerified.Bool
	}

	// Timestamps
	if user.CreatedAt.Valid {
		pbUser.CreatedAt = user.CreatedAt.Time.Unix()
	}
	if user.UpdatedAt.Valid {
		pbUser.UpdatedAt = user.UpdatedAt.Time.Unix()
	}

	return &userPb.UpdateUserResponse{
		User:    pbUser,
		Success: true,
	}, nil
}

// GetUserProfile retrieves detailed user profile
func (c *UserController) GetUserProfile(ctx context.Context, req *userPb.GetUserProfileRequest) (*userPb.GetUserProfileResponse, error) {
	if req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	user, err := c.userService.GetUser(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	pbUser := &schema.User{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}

	// All optional string fields
	if user.FullName.Valid {
		pbUser.FullName = user.FullName.String
	}
	if user.ProfilePic.Valid {
		pbUser.ProfilePic = user.ProfilePic.String
	}
	if user.Bio.Valid {
		pbUser.Bio = user.Bio.String
	}
	if user.ColorCode.Valid {
		pbUser.ColorCode = user.ColorCode.String
	}
	if user.BackgroundColor.Valid {
		pbUser.BackgroundColor = user.BackgroundColor.String
	}
	if user.BackgroundPic.Valid {
		pbUser.BackgroundPic = user.BackgroundPic.String
	}
	if user.CustomStatus.Valid {
		pbUser.CustomStatus = user.CustomStatus.String
	}

	// Boolean fields
	if user.IsBot.Valid {
		pbUser.IsBot = user.IsBot.Bool
	}
	if user.IsVerified.Valid {
		pbUser.IsVerified = user.IsVerified.Bool
	}

	// Timestamps
	if user.CreatedAt.Valid {
		pbUser.CreatedAt = user.CreatedAt.Time.Unix()
	}
	if user.UpdatedAt.Valid {
		pbUser.UpdatedAt = user.UpdatedAt.Time.Unix()
	}

	return &userPb.GetUserProfileResponse{
		User: pbUser,
	}, nil
}

// DeleteUser deletes a user account
func (c *UserController) DeleteUser(ctx context.Context, req *userPb.DeleteUserRequest) (*userPb.DeleteUserResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	// Can only delete your own account
	if req.GetUserId() != userID {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrPermissionDenied)
	}

	err := c.userService.DeleteUser(ctx, userID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.DeleteUserResponse{
		Success: true,
	}, nil
}

// UpdateStatus updates user status and presence
func (c *UserController) UpdateStatus(ctx context.Context, req *userPb.UpdateStatusRequest) (*userPb.UpdateStatusResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	var customStatus, activity *string
	if req.GetCustomStatus() != "" {
		cs := req.GetCustomStatus()
		customStatus = &cs
	}
	if req.GetActivity() != "" {
		a := req.GetActivity()
		activity = &a
	}

	err := c.userService.UpdateStatus(ctx, userID, req.GetStatus(), customStatus, activity)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.UpdateStatusResponse{
		Success: true,
	}, nil
}

// GetUserPresence retrieves user presence information
func (c *UserController) GetUserPresence(ctx context.Context, req *userPb.GetUserPresenceRequest) (*userPb.GetUserPresenceResponse, error) {
	if req.GetUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	presence, err := c.userService.GetUserPresence(ctx, req.GetUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	response := &userPb.GetUserPresenceResponse{}

	if presence.Status.Valid {
		response.Status = presence.Status.String
	}
	if presence.CustomStatus.Valid {
		response.CustomStatus = presence.CustomStatus.String
	}
	if presence.Activity.Valid {
		response.Activity = presence.Activity.String
	}
	if presence.LastSeen.Valid {
		response.LastSeen = presence.LastSeen.Time.Unix()
	}

	return response, nil
}

// SetCustomStatus sets user custom status
func (c *UserController) SetCustomStatus(ctx context.Context, req *userPb.SetCustomStatusRequest) (*userPb.SetCustomStatusResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	var customStatus, emoji *string
	var expiresAt *time.Time

	if req.GetCustomStatus() != "" {
		cs := req.GetCustomStatus()
		customStatus = &cs
	}
	if req.GetEmoji() != "" {
		e := req.GetEmoji()
		emoji = &e
	}
	if req.GetExpiresAt() > 0 {
		t := time.Unix(req.GetExpiresAt(), 0)
		expiresAt = &t
	}

	err := c.userService.SetCustomStatus(ctx, userID, customStatus, emoji, expiresAt)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.SetCustomStatusResponse{
		Success: true,
	}, nil
}

// UpdateUserSettings updates user settings
func (c *UserController) UpdateUserSettings(ctx context.Context, req *userPb.UpdateUserSettingsRequest) (*userPb.UpdateUserSettingsResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	var showActivity, allowDMs, enableNotif *bool
	var theme, language *string

	// Handle boolean fields properly
	if req.GetShowCurrentActivity() {
		val := req.GetShowCurrentActivity()
		showActivity = &val
	}
	if req.GetAllowDms() {
		val := req.GetAllowDms()
		allowDMs = &val
	}
	if req.GetEnableNotifications() {
		val := req.GetEnableNotifications()
		enableNotif = &val
	}
	if req.GetTheme() != "" {
		t := req.GetTheme()
		theme = &t
	}
	if req.GetLanguage() != "" {
		l := req.GetLanguage()
		language = &l
	}

	err := c.userService.UpdateUserSettings(ctx, userID, showActivity, allowDMs, enableNotif, theme, language)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.UpdateUserSettingsResponse{
		Success: true,
	}, nil
}

// GetUserSettings retrieves user settings
func (c *UserController) GetUserSettings(ctx context.Context, req *userPb.GetUserSettingsRequest) (*userPb.GetUserSettingsResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	settings, err := c.userService.GetUserSettings(ctx, userID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.GetUserSettingsResponse{
		ShowCurrentActivity: settings.ShowCurrentActivity,
		AllowDms:            settings.AllowDMs,
		EnableNotifications: settings.EnableNotifications,
		Theme:               settings.Theme,
		Language:            settings.Language,
	}, nil
}

// BlockUser blocks a user
func (c *UserController) BlockUser(ctx context.Context, req *userPb.BlockUserRequest) (*userPb.BlockUserResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetBlockedUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.userService.BlockUser(ctx, userID, req.GetBlockedUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.BlockUserResponse{
		Success: true,
	}, nil
}

// UnblockUser unblocks a user
func (c *UserController) UnblockUser(ctx context.Context, req *userPb.UnblockUserRequest) (*userPb.UnblockUserResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	if req.GetBlockedUserId() == 0 {
		return nil, commonErrors.ToGRPCError(commonErrors.ErrInvalidInput)
	}

	err := c.userService.UnblockUser(ctx, userID, req.GetBlockedUserId())
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.UnblockUserResponse{
		Success: true,
	}, nil
}

// GetBlockedUsers retrieves list of blocked users
func (c *UserController) GetBlockedUsers(ctx context.Context, req *userPb.GetBlockedUsersRequest) (*userPb.GetBlockedUsersResponse, error) {
	// Get user ID from context
	userID := ctx.Value("user_id").(int32)

	blockedIDs, err := c.userService.GetBlockedUsers(ctx, userID)
	if err != nil {
		return nil, commonErrors.ToGRPCError(err)
	}

	return &userPb.GetBlockedUsersResponse{
		BlockedUserIds: blockedIDs,
	}, nil
}
