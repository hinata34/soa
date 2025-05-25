package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"promo/internal/userservice/auth"
	postgresqlAuth "promo/internal/userservice/auth/postgresql"
	"promo/internal/userservice/db"
	"promo/internal/userservice/profile"
	postgresqlProfile "promo/internal/userservice/profile/postgresql"
	"promo/internal/userservice/proto/pb"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	database := db.NewDB(ctx, "")
	authDatabase := postgresqlAuth.NewAuthRepoPostgresql(database)
	profileDatabase := postgresqlProfile.NewProfileRepoPostgresql(database)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 123))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAuthServer(grpcServer, auth.NewAuthServer(authDatabase)) // maybe make RegisterService function
	pb.RegisterProfileServer(grpcServer, profile.NewProfileServer(profileDatabase))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}
