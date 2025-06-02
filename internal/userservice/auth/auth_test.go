package auth

import (
	"context"
	mock_repository "promo/internal/userservice/auth/mocks"
	model "promo/internal/userservice/models"
	pb "promo/internal/userservice/proto/pb"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authRepoFixture struct {
	ctrl *gomock.Controller
	repo *mock_repository.MockAuthRepo
	hash *mock_repository.MockHash
	auth authServer
}

func setUp(t *testing.T) authRepoFixture {
	ctrl := gomock.NewController(t)
	repo := mock_repository.NewMockAuthRepo(ctrl)
	hash := mock_repository.NewMockHash(ctrl)
	auth := authServer{db: repo, hashFunction: hash}

	return authRepoFixture{ctrl: ctrl, repo: repo, auth: auth}
}

func (u *authRepoFixture) tearDown() {
	u.ctrl.Finish()
}

func Test_Register(t *testing.T) {
	t.Parallel()
	var (
		ctx = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		a := setUp(t)
		defer a.tearDown()

		now := timestamppb.Now()
		var registerRequest pb.RegisterRequest
		registerRequest.Login = proto.String("hinata34")
		registerRequest.Email = proto.String("hinata34@yandex.ru")
		registerRequest.Password = proto.String("lolkek")
		registerRequest.Name = proto.String("Ivan")
		registerRequest.Surname = proto.String("Kochkarev")
		registerRequest.MobileNumber = proto.String("+7199999")
		registerRequest.Birthday = now

		password := []byte("lolkek")

		hash, _ := bcrypt.GenerateFromPassword(password, 10)
		user := &model.User{
			Login:        "hinata34",
			Password:     string(hash),
			Name:         "Ivan",
			Surname:      "Kochkarev",
			Birthday:     now.AsTime(),
			Email:        "hinata34@yandex.ru",
			MobileNumber: "+7199999",
		}

		a.repo.EXPECT().GetByLogin(gomock.Any(), "hinata34").Return(nil, nil)
		a.repo.EXPECT().Add(gomock.Any(), user).Return(uint64(1), nil)
		a.hash.EXPECT().GenerateFromPassword(password, 10).Return(hash, nil)

		_, err := a.auth.Register(ctx, &registerRequest)
		require.NoError(t, err)
	})
}
