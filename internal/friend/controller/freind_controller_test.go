package controller_test

import (
	"context"
	"discord/config"
	friendPb "discord/gen/proto/service/friend"
	"discord/gen/repo"
	"discord/internal/friend/controller"
	"discord/internal/friend/repository"
	"discord/internal/friend/service"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func load(t *testing.T) *controller.FriendController {
	err := config.InitDB()
	if err != nil {
		assert.Error(t, err)
	}
	friendRepo := repository.NewFriendRepository(config.DB)
	return controller.NewFriendControllerInstance(service.NewFriendService(friendRepo))
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

func TestAddFriend(t *testing.T) {
	ctx := context.Background()
	f := load(t)
	assert.NotNil(t, f)
	id1 := createrUser(t, ctx, "test", "test", "test")
	defer DeleteUser(t, ctx, id1)
	id2 := createrUser(t, ctx, "test2", "test2", "test2")
	defer DeleteUser(t, ctx, id2)
	v, e := f.SendFriendRequest(ctx, &friendPb.SendFriendRequestRequest{
		UserId:   int32(id1),
		FriendId: int32(id2),
	})
	fmt.Println("value is ", v)
	if e != nil {
		assert.Error(t, e)
	}
	p, err := f.GetPendingRequests(ctx, &friendPb.GetPendingRequestsRequest{
		UserId: int32(id1),
	})

	defer DeleteFriendship(t, ctx, id1, id2)
	fmt.Println("pending requests are ", p.PendingRequests, len(p.PendingRequests))
	assert.NoError(t, err)
	assert.Equal(t, len(p.PendingRequests), 1)
	assert.Equal(t, p.PendingRequests[0].FriendId, int32(id1))
	assert.Equal(t, v.Success, true)
}

func TestAcceptFriendRequest(t *testing.T) {
	ctx := context.Background()
	f := load(t)
	assert.NotNil(t, f)
	id1 := createrUser(t, ctx, "test", "test", "test")
	defer DeleteUser(t, ctx, id1)
	id2 := createrUser(t, ctx, "test2", "test2", "test2")
	defer DeleteUser(t, ctx, id2)
	v, e := f.SendFriendRequest(ctx, &friendPb.SendFriendRequestRequest{
		UserId:   int32(id1),
		FriendId: int32(id2),
	})
	fmt.Println("value is ", v)
	if e != nil {
		assert.Error(t, e)
	}
	ctx = context.WithValue(ctx, "user_id", int32(id2))
	p, err := f.AcceptFriendRequest(ctx, &friendPb.AcceptFriendRequestRequest{
		UserId:      int32(id2),
		RequesterId: int32(id1),
	})
	fmt.Println("value is ", p)
	if e != nil {
		assert.Error(t, e)
	}
	pf, err := f.GetFriends(ctx, &friendPb.GetFriendsRequest{
		UserId: int32(id1),
	})
	fmt.Println("friends are ", pf.Friends, len(pf.Friends))
	assert.NoError(t, err)
	assert.Equal(t, len(pf.Friends), 1)
	assert.Equal(t, v.Success, true)
}

func TestRejectFriendRequest(t *testing.T) {
	ctx := context.Background()
	f := load(t)
	assert.NotNil(t, f)
	id1 := createrUser(t, ctx, "test", "test", "test")
	defer DeleteUser(t, ctx, id1)
	id2 := createrUser(t, ctx, "test2", "test2", "test2")
	defer DeleteUser(t, ctx, id2)
	v, e := f.SendFriendRequest(ctx, &friendPb.SendFriendRequestRequest{
		UserId:   int32(id1),
		FriendId: int32(id2),
	})
	fmt.Println("value is ", v)
	if e != nil {
		assert.Error(t, e)
	}
	ctx = context.WithValue(ctx, "user_id", int32(id2))
	p, err := f.RejectFriendRequest(ctx, &friendPb.RejectFriendRequestRequest{
		UserId:      int32(id2),
		RequesterId: int32(id1),
	})
	fmt.Println("value is ", p)
	if e != nil {
		assert.Error(t, e)
	}
	pf, err := f.GetFriends(ctx, &friendPb.GetFriendsRequest{
		UserId: int32(id1),
	})
	fmt.Println("friends are ", pf.Friends, len(pf.Friends))
	assert.NoError(t, err)
	assert.Equal(t, len(pf.Friends), 0)
	assert.Equal(t, v.Success, true)
}

func TestBlockFriends(t *testing.T) {
	ctx := context.Background()
	f := load(t)
	assert.NotNil(t, f)
	id1 := createrUser(t, ctx, "test", "test", "test")
	defer DeleteUser(t, ctx, id1)
	id2 := createrUser(t, ctx, "test2", "test2", "test2")
	defer DeleteUser(t, ctx, id2)
	v, e := f.SendFriendRequest(ctx, &friendPb.SendFriendRequestRequest{
		UserId:   int32(id1),
		FriendId: int32(id2),
	})
	fmt.Println("value is ", v)
	if e != nil {
		assert.Error(t, e)
	}
	ctx = context.WithValue(ctx, "user_id", int32(id2))
	p, err := f.AcceptFriendRequest(ctx, &friendPb.AcceptFriendRequestRequest{
		UserId:      int32(id2),
		RequesterId: int32(id1),
	})
	fmt.Println("value is ", p)
	if e != nil {
		assert.Error(t, e)
	}
	pf, err := f.GetFriends(ctx, &friendPb.GetFriendsRequest{
		UserId: int32(id1),
	})
	fmt.Println("friends are ", pf.Friends, len(pf.Friends))
	assert.NoError(t, err)
	assert.Equal(t, len(pf.Friends), 1)
	assert.Equal(t, v.Success, true)

}
