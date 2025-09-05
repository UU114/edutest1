package logic

import (
	"context"

	"ai-education-platform/services/user/user/internal/models"
	"ai-education-platform/services/user/user/internal/svc"
	"ai-education-platform/services/user/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) (*types.CommonResponse, error) {
	// 从上下文获取用户ID
	userId := l.ctx.Value("user_id").(int64)
	
	// 查询当前用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户信息失败: %v", err)
		return nil, err
	}

	// 验证旧密码
	if !models.CheckPassword(req.OldPassword, user.Password) {
		l.Logger.Errorf("旧密码错误: %d", userId)
		return &types.CommonResponse{
			Success: false,
			Message: "旧密码错误",
		}, nil
	}

	// 加密新密码
	hashedPassword, err := models.HashPassword(req.NewPassword)
	if err != nil {
		l.Logger.Errorf("密码加密失败: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "密码修改失败",
		}, err
	}

	// 更新密码
	err = l.svcCtx.UserModel.UpdatePassword(l.ctx, userId, hashedPassword)
	if err != nil {
		l.Logger.Errorf("更新密码失败: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "密码修改失败",
		}, err
	}

	l.Logger.Infof("密码修改成功: %d", userId)
	return &types.CommonResponse{
		Success: true,
		Message: "密码修改成功",
	}, nil
}