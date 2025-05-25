package auth

import (
	"context"
	"log"
	model "promo/internal/userservice/models"
	"promo/internal/userservice/proto/pb"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
)

type authServer struct {
	pb.UnimplementedAuthServer
	db AuthRepo
}

func NewAuthServer(db AuthRepo) *authServer {
	return &authServer{db: db}
}

func (a *authServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var res pb.RegisterResponse
	if req.Login == nil || req.Password == nil || req.Email == nil {
		// error
		// set error everywhere
		return &res, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), 10)
	if err != nil {
		// handle error
		return &res, nil
	}

	user := model.User{ // check for existance
		Login:        req.GetLogin(),
		Password:     string(hash),
		Name:         req.GetName(),
		Surname:      req.GetSurname(),
		Birthday:     req.GetBirthday().AsTime(),
		Email:        req.GetEmail(),
		MobileNumber: req.GetMobileNumber(),
	}

	_, err = a.db.Add(ctx, &user)
	if err != nil {

		return &res, nil
	}

	res.Status = proto.Uint64(200)
	return &res, nil
}

func (a *authServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var res pb.LoginResponse
	if req.Password == nil {
		// error
		return &res, nil
	} else if req.Login == nil && req.Password == nil {
		// error
		return &res, nil
	}

	var user *model.User
	var err error
	if req.Login != nil {
		user, err = a.db.GetByLogin(ctx, req.GetLogin())
		if err != nil {
			// handle error
			return &res, err
		}
	} else if req.Email != nil {
		user, err = a.db.GetByEmail(ctx, req.GetEmail())
		if err != nil {
			// handle error
			return &res, err
		}
	}

	jwt, err := a.generateJWT(ctx, user)
	if err != nil {
		// handle error
		return &res, err
	}

	res.Status = proto.Uint64(200)
	res.Jwt = proto.String(jwt)
	return &res, nil
}

func (a *authServer) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	var res pb.ValidateResponse
	if req.Jwt == nil {
		// handle error
	}

	token, err := jwt.Parse(req.GetJwt(), func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil {
		// handle error
	}

	if !token.Valid {
		// handle error
	}

	res.Valid = proto.Bool(true)
	return &res, nil
}

func (a *authServer) checkIfUserExistsById(ctx context.Context, id uint64) (bool, error) {
	user, err := a.db.GetById(ctx, id)
	if err != nil {
		log.Printf("error: %v", err) // check for error
		return false, err
	}
	return user == nil, err
}

func (a *authServer) checkIfUserExistsByLogin(ctx context.Context, login string) (bool, error) {
	user, err := a.db.GetByLogin(ctx, login)
	if err != nil {
		log.Printf("error: %v", err) // check for error
		return false, err
	}
	return user == nil, err
}

func (a *authServer) generateJWT(ctx context.Context, user *model.User) (string, error) {
	secretKey := []byte("secret_key")

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": username,                         // Subject (user identifier)
		"iss": "todo-app",                       // Issuer
		"aud": getRole(username),                // Audience (user role)
		"exp": time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat": time.Now().Unix(),                // Issued at
	})

	jwt, err := claims.SignedString(secretKey)
	if err != nil {
		// handle error
	}

	return jwt, err
}
