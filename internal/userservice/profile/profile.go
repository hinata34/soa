package profile

import (
	"context"
	"errors"
	model "promo/internal/userservice/models"
	"promo/internal/userservice/proto/pb"

	"golang.org/x/crypto/bcrypt"
)

type profileServer struct {
	pb.UnimplementedProfileServer
	db ProfileRepo
}

func NewProfileServer(db ProfileRepo) *profileServer {
	return &profileServer{db: db}
}

func (p *profileServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	if req.Login == nil {
		return nil, errors.New("empty login")
	}

	user, err := p.db.GetByLogin(ctx, *req.Login)
	if err != nil {
		return nil, err
	}

	response := &pb.GetProfileResponse{Login: &user.Login, Email: &user.Email, Name: &user.Name, Surname: &user.Surname, MobileNumber: &user.MobileNumber}
	return response, nil
}

func (p *profileServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	if req.Login == nil {
		return nil, errors.New("empty login")
	} else if req.Name == nil || req.Surname == nil || req.MobileNumber == nil {
		return nil, errors.New("some of required parameters are empty")
	}

	user := &model.User{
		Login:        req.GetLogin(),
		Password:     "",
		Name:         req.GetName(),
		Surname:      req.GetSurname(),
		MobileNumber: req.GetMobileNumber(),
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 10)
		if err != nil {
			return nil, err
		}
		user.Password = string(hash)
	}

	err := p.db.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProfileResponse{}, nil
}
