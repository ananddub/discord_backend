package controller

import (
	"context"
	"discord/config"
	dmPb "discord/gen/proto/service/dm"
	"discord/gen/repo"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func load(t *testing.T) *DMController {
	err := config.InitDB()
	if err != nil {
		assert.Error(t, err)
	}
	return NewDMController(config.DB)
}

func createrUser(t *testing.T, ctx context.Context, name string, email string, password string) int {
	userRepo := repo.New(config.DB)
	userRepo.DeleteByUsername(ctx, name)
	v, err := userRepo.CreateUser(ctx, repo.CreateUserParams{
		Username: name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		assert.Error(t, err)
	}
	return int(v.ID)
}

func DeleteUser(t *testing.T, ctx context.Context, id int) {
	userRepo := repo.New(config.DB)
	userRepo.DeleteUserById(ctx, int32(id))
}
func DeleteFriendship(t *testing.T, ctx context.Context, uid, fid int) {
	friendRepo := repo.New(config.DB)
	friendRepo.HardDeleteFriendship(ctx, repo.HardDeleteFriendshipParams{
		UserID:   int32(uid),
		FriendID: int32(fid),
	})
}
func TestCreateDM(t *testing.T) {
	ctx := context.Background()
	f := load(t)
	assert.NotNil(t, f)
	id1 := createrUser(t, ctx, "test", "test", "test")
	defer DeleteUser(t, ctx, id1)
	id2 := createrUser(t, ctx, "test2", "test2", "test2")
	defer DeleteUser(t, ctx, id2)
	_, e := f.SendMessage(ctx, &dmPb.SendMessageRequest{
		SenderId:         int32(id1),
		ReceiverId:       int32(id2),
		ReplyToMessageId: nil,
		Content:          "Hello World",
	})
	if e != nil {
		assert.Error(t, e)
	}
	chat := repo.New(config.DB)
	msg, err := chat.GetChatMessages(ctx, repo.GetChatMessagesParams{
		SenderID:   int32(id1),
		ReceiverID: pgtype.Int4{Valid: true, Int32: int32(id2)},
		Limit:      10,
		Offset:     0,
	})
	defer func() {
		DeleteUser(t, ctx, id1)
		DeleteUser(t, ctx, id2)
		chat.HardDeleteChatMessages(ctx, repo.HardDeleteChatMessagesParams{
			SenderID:   int32(id1),
			ReceiverID: pgtype.Int4{Valid: true, Int32: int32(id2)},
		})
	}()
	assert.ErrorIs(t, err, nil)
	assert.NotNil(t, msg)
	assert.Equal(t, "Hello World", msg[0].Content)
}
