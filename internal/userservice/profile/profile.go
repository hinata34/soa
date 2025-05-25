package profile

import (
	"context"
	"promo/internal/userservice/proto/pb"
)

type profileServer struct {
	pb.UnimplementedProfileServer
	db ProfileRepo
}

func NewProfileServer(db ProfileRepo) *profileServer {
	return &profileServer{db: db}
}

func (p *profileServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
}

func (p *profileServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {

}
