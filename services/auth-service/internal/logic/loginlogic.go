package logic

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/aicong/mine-dispatch/proto/user/v1"
	"github.com/aicong/mine-dispatch/services/auth-service/internal/svc"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		Token    string       `json:"token"`
		ExpireAt int64        `json:"expire_at"`
		User     *userv1.User `json:"user"`
	} `json:"data,omitempty"`
}

type LoginLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svc *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{ctx: ctx, svc: svc}
}

func (l *LoginLogic) Login(req *LoginRequest) (*LoginResponse, error) {
	var userModel svc.UserModel
	if err := l.svc.DB.Where("username = ?", req.Username).First(&userModel).Error; err != nil {
		return &LoginResponse{Code: 401, Message: "invalid credentials"}, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(req.Password)); err != nil {
		return &LoginResponse{Code: 401, Message: "invalid credentials"}, nil
	}

	userResp, _ := l.svc.UserRpc.GetUser(l.ctx, &userv1.GetUserRequest{Id: userModel.ID})
	user := userResp.GetData()

	expireAt := time.Now().Add(time.Duration(l.svc.Config.JwtExpire) * time.Second).Unix()
	claims := jwt.MapClaims{
		"user_id":  user.GetId(),
		"username": user.GetUsername(),
		"role":     user.GetRole(),
		"mine_id":  user.GetMineId(),
		"exp":      expireAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(l.svc.Config.JwtSecret))
	if err != nil {
		return &LoginResponse{Code: 500, Message: "token generation failed"}, nil
	}

	return &LoginResponse{
		Code:    0,
		Message: "success",
		Data: &struct {
			Token    string       `json:"token"`
			ExpireAt int64        `json:"expire_at"`
			User     *userv1.User `json:"user"`
		}{
			Token:    tokenStr,
			ExpireAt: expireAt,
			User:     user,
		},
	}, nil
}
