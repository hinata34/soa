package app

import (
	"context"
	"log"
	apigateway "promo/internal/apigateway/swagger"
	"promo/internal/userservice/proto/pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

type ApigatewayServer struct {
	authClient    pb.AuthClient
	profileClient pb.ProfileClient
}

func NewApigatewayServer(serverAddr string) *ApigatewayServer {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	authClient := pb.NewAuthClient(conn)
	profileClient := pb.NewProfileClient(conn)
	return &ApigatewayServer{authClient: authClient, profileClient: profileClient}
}

func (s *ApigatewayServer) Register(ctx *gin.Context) {
	var body apigateway.RegisterJSONRequestBody
	if err := ctx.ShouldBind(&body); err != nil {
		errorResponse := &apigateway.RegisterdefaultJSONResponse{StatusCode: 400, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitRegisterResponse(ctx.Writer)
		return
	}

	registerRequest := &pb.RegisterRequest{Login: &body.Login, Email: &body.Email, Password: &body.Password, Surname: body.Surname}
	_, err := s.authClient.Register(context.Background(), registerRequest)
	if err != nil {
		errorResponse := &apigateway.RegisterdefaultJSONResponse{StatusCode: 400, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitRegisterResponse(ctx.Writer)
		return
	}

	var response apigateway.Register200JSONResponse
	response.VisitRegisterResponse(ctx.Writer)
}

func (s *ApigatewayServer) Login(ctx *gin.Context) {
	var body apigateway.LoginJSONRequestBody
	if err := ctx.ShouldBind(&body); err != nil {
		errorResponse := &apigateway.LogindefaultJSONResponse{StatusCode: 400, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitLoginResponse(ctx.Writer)
		return
	}

	if body.Email == nil && body.Login == nil {
		errorResponse := &apigateway.LogindefaultJSONResponse{StatusCode: 400, Body: apigateway.Error{Message: "Email or login == nil"}}
		errorResponse.VisitLoginResponse(ctx.Writer)
		return
	}

	loginRequest := &pb.LoginRequest{Login: body.Login, Email: body.Email, Password: &body.Password}
	loginResponse, err := s.authClient.Login(context.Background(), loginRequest)
	if err != nil {
		errorResponse := &apigateway.LogindefaultJSONResponse{StatusCode: 404, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitLoginResponse(ctx.Writer)
		return
	}

	var response apigateway.Login200JSONResponse
	response.Jwt = loginResponse.Jwt
	response.VisitLoginResponse(ctx.Writer)
}

func (s *ApigatewayServer) GetProfile(ctx *gin.Context, params apigateway.GetProfileParams) {
	login, err := s.validateJwt(&params.Jwt)
	if err != nil {
		errorResponse := &apigateway.GetProfiledefaultJSONResponse{StatusCode: 403, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitGetProfileResponse(ctx.Writer)
		return
	}

	getProfileRequest := &pb.GetProfileRequest{Login: login}
	getProfileResponse, err := s.profileClient.GetProfile(context.Background(), getProfileRequest)
	if err != nil {
		errorResponse := &apigateway.GetProfiledefaultJSONResponse{StatusCode: 404, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitGetProfileResponse(ctx.Writer)
		return
	}

	var response apigateway.GetProfile200JSONResponse
	response.Login = proto.String(getProfileResponse.GetLogin())
	response.Email = proto.String(getProfileResponse.GetEmail())
	response.Name = proto.String(getProfileResponse.GetName())
	response.Surname = proto.String(getProfileResponse.GetSurname())
	response.MobileNumber = proto.String(getProfileResponse.GetMobileNumber())
	response.VisitGetProfileResponse(ctx.Writer)
}

func (s *ApigatewayServer) UpdateProfile(ctx *gin.Context, params apigateway.UpdateProfileParams) {
	login, err := s.validateJwt(&params.Jwt)
	if err != nil {
		errorResponse := &apigateway.UpdateProfiledefaultJSONResponse{StatusCode: 403, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitUpdateProfileResponse(ctx.Writer)
		return
	}

	var body apigateway.UpdateProfileJSONRequestBody
	if err := ctx.ShouldBind(&body); err != nil {
		errorResponse := &apigateway.UpdateProfiledefaultJSONResponse{StatusCode: 400, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitUpdateProfileResponse(ctx.Writer)
		return
	}
	log.Println(*login, body.Name, body.Surname, body.MobileNumber)
	updateProfileRequest := &pb.UpdateProfileRequest{Login: login, Name: &body.Name, Surname: &body.Surname, MobileNumber: &body.MobileNumber, Password: body.Password}
	_, err = s.profileClient.UpdateProfile(context.Background(), updateProfileRequest)
	if err != nil {
		errorResponse := &apigateway.UpdateProfiledefaultJSONResponse{StatusCode: 404, Body: apigateway.Error{Message: err.Error()}}
		errorResponse.VisitUpdateProfileResponse(ctx.Writer)
		return
	}

	response := apigateway.UpdateProfile200TextResponse("true")
	response.VisitUpdateProfileResponse(ctx.Writer)
}

func (s *ApigatewayServer) validateJwt(jwt *string) (*string, error) {
	validateRequest := &pb.ValidateRequest{Jwt: jwt}
	validateResponse, err := s.authClient.Validate(context.Background(), validateRequest)
	if !validateResponse.GetValid() && err != nil {
		return nil, err
	}

	return validateResponse.Login, nil
}
