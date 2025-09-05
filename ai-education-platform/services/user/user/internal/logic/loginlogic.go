package logic

import (
	"context"
	"time"

	"ai-education-platform/services/user/user/internal/models"
	"ai-education-platform/services/user/user/internal/svc"
	"ai-education-platform/services/user/user/internal/types"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (*types.LoginResponse, error) {
	// 查询用户
	user, err := l.svcCtx.UserModel.FindByUsername(l.ctx, req.Username)
	if err != nil {
		l.Logger.Errorf("用户不存在: %s", req.Username)
		return nil, err
	}

	// 验证密码
	if !models.CheckPassword(req.Password, user.Password) {
		l.Logger.Errorf("密码错误: %s", req.Username)
		return nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		l.Logger.Errorf("用户状态异常: %s, status: %d", req.Username, user.Status)
		return nil, err
	}

	// 生成JWT token
	now := time.Now().Unix()
	accessToken, err := l.generateToken(user.ID, user.Username, user.Role, now+l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		l.Logger.Errorf("生成token失败: %v", err)
		return nil, err
	}

	// 生成refresh token
	refreshToken, err := l.generateToken(user.ID, user.Username, user.Role, now+l.svcCtx.Config.Auth.AccessExpire*2)
	if err != nil {
		l.Logger.Errorf("生成refresh token失败: %v", err)
		return nil, err
	}

	// 更新最后登录时间
	err = l.svcCtx.UserModel.UpdateLastLogin(l.ctx, user.ID)
	if err != nil {
		l.Logger.Errorf("更新最后登录时间失败: %v", err)
	}

	l.Logger.Infof("用户登录成功: %s", req.Username)
	return &types.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
		UserInfo:     *user.ToUserInfo(),
	}, nil
}

// 生成JWT token
func (l *LoginLogic) generateToken(userId int64, username, role string, expire int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"role":     role,
		"exp":      expire,
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.Auth.AccessSecret))
}