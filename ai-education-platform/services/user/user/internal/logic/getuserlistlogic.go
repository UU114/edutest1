package logic

import (
	"context"

	"ai-education-platform/services/user/user/internal/svc"
	"ai-education-platform/services/user/user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(req *types.UserListRequest) (*types.UserListResponse, error) {
	// 查询用户列表
	users, total, err := l.svcCtx.UserModel.FindList(l.ctx, req)
	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	// 转换为UserInfo格式
	userInfos := make([]*types.UserInfo, 0, len(users))
	for _, user := range users {
		userInfos = append(userInfos, user.ToUserInfo())
	}

	return &types.UserListResponse{
		Total: total,
		List:  userInfos,
	}, nil
}