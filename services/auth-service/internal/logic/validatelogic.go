package logic

import (
	"context"
	"github.com/golang-jwt/jwt/v5"

	"github.com/aicong/mine-dispatch/services/auth-service/internal/svc"
)

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *struct {
		UserId   uint64 `json:"user_id"`
		Username string `json:"username"`
		Role     int32  `json:"role"`
		MineId   uint64 `json:"mine_id"`
	} `json:"data,omitempty"`
}

type ValidateLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewValidateLogic(ctx context.Context, svc *svc.ServiceContext) *ValidateLogic {
	return &ValidateLogic{ctx: ctx, svc: svc}
}

func (l *ValidateLogic) Validate(req *ValidateRequest) (*ValidateResponse, error) {
	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(l.svc.Config.JwtSecret), nil
	})
	if err != nil || !token.Valid {
		return &ValidateResponse{Code: 401, Message: "invalid token"}, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &ValidateResponse{Code: 401, Message: "invalid token claims"}, nil
	}

	return &ValidateResponse{
		Code:    0,
		Message: "success",
		Data: &struct {
			UserId   uint64 `json:"user_id"`
			Username string `json:"username"`
			Role     int32  `json:"role"`
			MineId   uint64 `json:"mine_id"`
		}{
			UserId:   uint64(claims["user_id"].(float64)),
			Username: claims["username"].(string),
			Role:     int32(claims["role"].(float64)),
			MineId:   uint64(claims["mine_id"].(float64)),
		},
	}, nil
}
