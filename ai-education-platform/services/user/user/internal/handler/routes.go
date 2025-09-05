package handler

import (
	"net/http"

	"ai-education-platform/services/user/user/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 用户相关接口
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/user/register",
				Handler: RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/user/login",
				Handler: LoginHandler(serverCtx),
			},
		},
	)

	// 需要认证的接口
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/user/info",
				Handler: GetUserInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/v1/user/info",
				Handler: UpdateUserInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/v1/user/password",
				Handler: ChangePasswordHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	// 管理员接口
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/admin/users",
				Handler: GetUserListHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithPrefix("/api/v1/admin"),
	)
}