package logic

import (
	"context"

	"github.com/aicong/mine-dispatch/services/ai-service/internal/svc"
)

type AiLogic struct {
	ctx context.Context
	svc *svc.ServiceContext
}

func NewAiLogic(ctx context.Context, svc *svc.ServiceContext) *AiLogic {
	return &AiLogic{ctx: ctx, svc: svc}
}
