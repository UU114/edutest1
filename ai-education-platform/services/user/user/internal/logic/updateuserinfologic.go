package logic

import (
	"context"
	"time"

	"ai-education-platform/services/user/user/internal/models"
	"ai-education-platform/services/user/user/internal/svc"
	"ai-education-platform/services/user/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserRequest) (*types.CommonResponse, error) {
	// 从上下文获取用户ID
	userId := l.ctx.Value("user_id").(int64)
	
	// 查询当前用户信息
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, userId)
	if err != nil {
		l.Logger.Errorf("查询用户信息失败: %v", err)
		return nil, err
	}

	// 更新用户信息
	now := time.Now().Unix()
	updatedUser := &models.User{
		ID:        userId,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Status:    user.Status,
		RealName:  user.RealName,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		School:    user.School,
		Grade:     user.Grade,
		Class:     user.Class,
		ParentId:  user.ParentId,
		InstitutionId: user.InstitutionId,
		Bio:       user.Bio,
		UpdatedAt: now,
	}

	// 应用更新
	if req.Nickname != "" {
		updatedUser.Nickname = req.Nickname
	}

	if req.Avatar != "" {
		updatedUser.Avatar.String = req.Avatar
		updatedUser.Avatar.Valid = true
	}

	if req.RealName != "" {
		updatedUser.RealName.String = req.RealName
		updatedUser.RealName.Valid = true
	}

	if req.Gender != 0 {
		updatedUser.Gender.Int64 = int64(req.Gender)
		updatedUser.Gender.Valid = true
	}

	if req.Birthday != "" {
		// 这里可以添加生日格式验证
		updatedUser.Birthday.String = req.Birthday
		updatedUser.Birthday.Valid = true
	}

	if req.School != "" {
		updatedUser.School.String = req.School
		updatedUser.School.Valid = true
	}

	if req.Grade != "" {
		updatedUser.Grade.String = req.Grade
		updatedUser.Grade.Valid = true
	}

	if req.Class != "" {
		updatedUser.Class.String = req.Class
		updatedUser.Class.Valid = true
	}

	if req.Bio != "" {
		updatedUser.Bio.String = req.Bio
		updatedUser.Bio.Valid = true
	}

	// 更新数据库
	err = l.svcCtx.UserModel.Update(l.ctx, updatedUser)
	if err != nil {
		l.Logger.Errorf("更新用户信息失败: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "更新失败",
		}, err
	}

	l.Logger.Infof("用户信息更新成功: %d", userId)
	return &types.CommonResponse{
		Success: true,
		Message: "更新成功",
	}, nil
}