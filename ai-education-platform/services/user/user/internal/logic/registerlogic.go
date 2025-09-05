package logic

import (
	"context"
	"time"

	"ai-education-platform/services/user/user/internal/models"
	"ai-education-platform/services/user/user/internal/svc"
	"ai-education-platform/services/user/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (*types.CommonResponse, error) {
	// 检查用户名是否已存在
	_, err := l.svcCtx.UserModel.FindByUsername(l.ctx, req.Username)
	if err == nil {
		return &types.CommonResponse{
			Success: false,
			Message: "用户名已存在",
		}, nil
	}

	// 检查邮箱是否已存在
	_, err = l.svcCtx.UserModel.FindByEmail(l.ctx, req.Email)
	if err == nil {
		return &types.CommonResponse{
			Success: false,
			Message: "邮箱已存在",
		}, nil
	}

	// 密码加密
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		l.Logger.Errorf("密码加密失败: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "注册失败",
		}, err
	}

	// 创建用户
	now := time.Now().Unix()
	user := &models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Nickname:  req.Nickname,
		Role:      req.Role,
		Status:    1, // 正常状态
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 设置可选字段
	if req.Phone != "" {
		user.Phone.String = req.Phone
		user.Phone.Valid = true
	}

	if req.RealName != "" {
		user.RealName.String = req.RealName
		user.RealName.Valid = true
	}

	if req.School != "" {
		user.School.String = req.School
		user.School.Valid = true
	}

	if req.Grade != "" {
		user.Grade.String = req.Grade
		user.Grade.Valid = true
	}

	if req.Class != "" {
		user.Class.String = req.Class
		user.Class.Valid = true
	}

	if req.ParentId > 0 {
		user.ParentId.Int64 = req.ParentId
		user.ParentId.Valid = true
	}

	if req.InstitutionId > 0 {
		user.InstitutionId.Int64 = req.InstitutionId
		user.InstitutionId.Valid = true
	}

	// 插入数据库
	_, err = l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		l.Logger.Errorf("用户注册失败: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "注册失败",
		}, err
	}

	l.Logger.Infof("用户注册成功: %s", req.Username)
	return &types.CommonResponse{
		Success: true,
		Message: "注册成功",
	}, nil
}