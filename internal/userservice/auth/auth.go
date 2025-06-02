package auth

import (
	"context"
	"errors"
	model "promo/internal/userservice/models"
	"promo/internal/userservice/proto/pb"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

type authServer struct {
	pb.UnimplementedAuthServer
	db           AuthRepo
	hashFunction Hash
}

func NewAuthServer(db AuthRepo, hashFunction Hash) *authServer {
	return &authServer{db: db, hashFunction: hashFunction}
}

func (a *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Login == nil || req.Password == nil || req.Email == nil {
		return nil, errors.New("empty login or password or email")
	}

	hash, err := a.hashFunction.GenerateFromPassword([]byte(req.GetPassword()), 10)
	if err != nil {
		return nil, err
	}

	isExist, _ := a.checkIfUserExistsByLogin(ctx, *req.Login)
	if isExist {
		return nil, errors.New("login already exists")
	}

	user := &model.User{
		Login:        req.GetLogin(),
		Password:     string(hash),
		Name:         req.GetName(),
		Surname:      req.GetSurname(),
		Birthday:     req.GetBirthday().AsTime(),
		Email:        req.GetEmail(),
		MobileNumber: req.GetMobileNumber(),
	}

	_, err = a.db.Add(ctx, user)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{}, nil
}

func (a *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Password == nil {
		return nil, errors.New("empty password")
	} else if req.Login == nil && req.Email == nil {
		return nil, errors.New("empty login and email")
	}

	var user *model.User
	var err error
	var hashedPassword string
	if req.Login != nil {
		user, err = a.db.GetByLogin(ctx, req.GetLogin())
		if err != nil {
			return nil, err
		}
		hashedPassword = user.Password
	} else if req.Email != nil {
		user, err = a.db.GetByEmail(ctx, req.GetEmail())
		if err != nil {
			return nil, err
		}
		hashedPassword = user.Password
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(*req.Password)); err != nil {
		return nil, err
	}

	jwt, err := a.generateJWT(ctx, user)
	if err != nil {
		return nil, err
	}

	var response pb.LoginResponse
	response.Jwt = proto.String(jwt)
	return &response, nil
}

func (a *authServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	if req.Jwt == nil {
		return nil, errors.New("empty jwt")
	}

	token, err := jwt.Parse(req.GetJwt(), func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("jwt token is invalid")
	}

	login, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	var res pb.ValidateResponse
	res.Valid = proto.Bool(true)
	res.Login = proto.String(login)
	return &res, nil
}

func (a *authServer) checkIfUserExistsById(ctx context.Context, id uint64) (bool, error) {
	user, err := a.db.GetById(ctx, id)
	if err != nil {
		return false, err
	}
	return user == nil, err
}

func (a *authServer) checkIfUserExistsByLogin(ctx context.Context, login string) (bool, error) {
	user, err := a.db.GetByLogin(ctx, login) // rewrite with exec instead of query
	return user != nil, err
}

func (a *authServer) generateJWT(ctx context.Context, user *model.User) (string, error) {
	secretKey := []byte("secret_key")
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Login,          // Subject (user identifier)
		"iss": "promocode-service", // Issuer
		// "aud": getRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	jwt, err := claims.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
