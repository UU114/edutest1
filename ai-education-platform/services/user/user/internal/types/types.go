package types

// 用户角色枚举
type UserRole struct {
	Role string `json:"role" options="[student,teacher,parent,institution]"`
}

// 用户基础信息
type UserInfo struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	Nickname     string    `json:"nickname"`
	Avatar       string    `json:"avatar,omitempty"`
	Role         string    `json:"role" options="[student,teacher,parent,institution]"`
	Status       int       `json:"status"` // 0:禁用 1:正常 2:待审核
	CreatedAt    int64     `json:"created_at"`
	UpdatedAt    int64     `json:"updated_at"`
	LastLoginAt  int64     `json:"last_login_at,omitempty"`
}

// 用户详细信息
type UserProfile struct {
	UserInfo
	RealName     string    `json:"real_name,omitempty"`
	Gender       int       `json:"gender,omitempty"` // 0:未知 1:男 2:女
	Birthday     string    `json:"birthday,omitempty"`
	School       string    `json:"school,omitempty"`
	Grade        string    `json:"grade,omitempty"` // 年级
	Class        string    `json:"class,omitempty"` // 班级
	ParentId     int64     `json:"parent_id,omitempty"` // 家长ID
	InstitutionId int64     `json:"institution_id,omitempty"` // 机构ID
	Bio          string    `json:"bio,omitempty"` // 个人简介
}

// 注册请求
type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=20"`
	Password    string `json:"password" validate:"required,min=6,max=20"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone,omitempty"`
	Nickname    string `json:"nickname" validate:"required,min=1,max=30"`
	Role        string `json:"role" validate:"required,options=[student,teacher,parent,institution]"`
	RealName    string `json:"real_name,omitempty"`
	School      string `json:"school,omitempty"`
	Grade       string `json:"grade,omitempty"`
	Class       string `json:"class,omitempty"`
	ParentId    int64  `json:"parent_id,omitempty"`
	InstitutionId int64  `json:"institution_id,omitempty"`
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	UserInfo
}

// 更新用户信息请求
type UpdateUserRequest struct {
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	RealName     string `json:"real_name,omitempty"`
	Gender       int    `json:"gender,omitempty"`
	Birthday     string `json:"birthday,omitempty"`
	School       string `json:"school,omitempty"`
	Grade        string `json:"grade,omitempty"`
	Class        string `json:"class,omitempty"`
	Bio          string `json:"bio,omitempty"`
}

// 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=20"`
}

// 用户列表查询请求
type UserListRequest struct {
	Page     int    `json:"page" validate:"required,min=1"`
	PageSize int    `json:"page_size" validate:"required,min=1,max=100"`
	Role     string `json:"role,omitempty" options="[student,teacher,parent,institution]"`
	Keyword  string `json:"keyword,omitempty"`
	Status   int    `json:"status,omitempty"`
}

// 用户列表响应
type UserListResponse struct {
	Total int64        `json:"total"`
	List  []UserInfo   `json:"list"`
}

// 家长孩子列表响应
type ParentChildrenResponse struct {
	Children []UserInfo `json:"children"`
}

// 机构教师列表响应
type InstitutionTeachersResponse struct {
	Teachers []UserInfo `json:"teachers"`
}

// 通用响应
type CommonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}