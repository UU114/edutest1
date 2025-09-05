package middleware

import (
	"context"
	"net/http"
	"strings"

	"ai-education-platform/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// JWT认证中间件
type JWTMiddleware struct {
	secret string
	logx.Logger
}

func NewJWTMiddleware(secret string) *JWTMiddleware {
	return &JWTMiddleware{
		secret: secret,
		Logger: logx.WithContext(context.Background()),
	}
}

func (m *JWTMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取Authorization头
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.Error(w, &utils.Response{
				Success: false,
				Message: "缺少Authorization头",
				Code:    401,
			})
			return
		}

		// 解析Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			httpx.Error(w, &utils.Response{
				Success: false,
				Message: "无效的token格式",
				Code:    401,
			})
			return
		}

		// 验证token
		jwtUtil := utils.NewJWTUtil(m.secret)
		claims, err := jwtUtil.ParseToken(tokenString)
		if err != nil {
			m.Logger.Errorf("JWT验证失败: %v", err)
			httpx.Error(w, &utils.Response{
				Success: false,
				Message: "无效的token",
				Code:    401,
			})
			return
		}

		// 将用户信息添加到上下文
		ctx := context.WithValue(r.Context(), "user_id", (*claims)["user_id"])
		ctx = context.WithValue(ctx, "username", (*claims)["username"])
		ctx = context.WithValue(ctx, "role", (*claims)["role"])

		// 继续处理请求
		next(w, r.WithContext(ctx))
	}
}

// 角色权限中间件
type RoleMiddleware struct {
	allowedRoles []string
	logx.Logger
}

func NewRoleMiddleware(allowedRoles []string) *RoleMiddleware {
	return &RoleMiddleware{
		allowedRoles: allowedRoles,
		Logger:      logx.WithContext(context.Background()),
	}
}

func (m *RoleMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从上下文获取用户角色
		role, ok := r.Context().Value("role").(string)
		if !ok {
			httpx.Error(w, &utils.Response{
				Success: false,
				Message: "无法获取用户角色",
				Code:    403,
			})
			return
		}

		// 检查角色权限
		if !utils.Contains(m.allowedRoles, role) {
			m.Logger.Errorf("用户角色无权限访问: %s", role)
			httpx.Error(w, &utils.Response{
				Success: false,
				Message: "权限不足",
				Code:    403,
			})
			return
		}

		// 继续处理请求
		next(w, r)
	}
}

// CORS中间件
func CORSMiddleware() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		http.DefaultServeMux.ServeHTTP(w, r)
	}
}

// 日志中间件
type LoggingMiddleware struct {
	logx.Logger
}

func NewLoggingMiddleware() *LoggingMiddleware {
	return &LoggingMiddleware{
		Logger: logx.WithContext(context.Background()),
	}
}

func (m *LoggingMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := utils.GetCurrentTimestamp()
		
		// 记录请求开始
		m.Logger.Infof("开始处理请求: %s %s", r.Method, r.URL.Path)
		
		// 创建响应写入器来捕获状态码
		rw := &responseWriter{ResponseWriter: w}
		
		// 调用下一个处理器
		next(rw, r)
		
		// 记录请求完成
		duration := utils.GetCurrentTimestamp() - start
		m.Logger.Infof("请求处理完成: %s %s - 状态码: %d, 耗时: %dms", 
			r.Method, r.URL.Path, rw.statusCode, duration)
	}
}

// 响应写入器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// 限流中间件
type RateLimitMiddleware struct {
	rate   int
	bucket int
	logx.Logger
}

func NewRateLimitMiddleware(rate, bucket int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		rate:   rate,
		bucket: bucket,
		Logger: logx.WithContext(context.Background()),
	}
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 这里可以实现基于Redis的分布式限流
		// 为了简化，这里只是示例
		// 实际项目中应该使用Redis或内存限流算法
		
		// 获取客户端IP
		ip := r.RemoteAddr
		
		// 简单的限流逻辑（实际应该使用更复杂的算法）
		// 这里只是示例，实际应该使用Redis或内存缓存
		
		m.Logger.Debugf("限流检查: IP=%s", ip)
		
		// 继续处理请求
		next(w, r)
	}
}

// 错误处理中间件
type ErrorHandlerMiddleware struct {
	logx.Logger
}

func NewErrorHandlerMiddleware() *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		Logger: logx.WithContext(context.Background()),
	}
}

func (m *ErrorHandlerMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.Logger.Errorf("请求处理异常: %v", err)
				httpx.Error(w, &utils.Response{
					Success: false,
					Message: "服务器内部错误",
					Code:    500,
				})
			}
		}()

		next(w, r)
	}
}

// 中间件链
func MiddlewareChain(next http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}