package logic

import (
	"context"
	"time"

	"ai-education-platform/services/user/user/internal/svc"
	"ai-education-platform/services/user/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (*types.UserProfile, error) {
	// 从上下文获取用户ID
	userId := l.ctx.Value("user_id").(int64)
	
	// 查询用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户信息失败: %v", err)
		return nil, err
	}

	return user.ToUserProfile(), nil
}