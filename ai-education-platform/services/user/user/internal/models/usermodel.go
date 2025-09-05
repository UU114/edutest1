package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"ai-education-platform/services/user/user/internal/types"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	conn sqlx.SqlConn
}

func NewUserModel(conn sqlx.SqlConn) *UserModel {
	return &UserModel{
		conn: conn,
	}
}

// 用户数据结构
type User struct {
	ID            int64          `db:"id"`
	Username      string         `db:"username"`
	Password      string         `db:"password"`
	Email         string         `db:"email"`
	Phone         sql.NullString `db:"phone"`
	Nickname      string         `db:"nickname"`
	Avatar        sql.NullString `db:"avatar"`
	Role          string         `db:"role"`
	Status        int            `db:"status"`
	RealName      sql.NullString `db:"real_name"`
	Gender        sql.NullInt64  `db:"gender"`
	Birthday      sql.NullTime   `db:"birthday"`
	School        sql.NullString `db:"school"`
	Grade         sql.NullString `db:"grade"`
	Class         sql.NullString `db:"class"`
	ParentId      sql.NullInt64  `db:"parent_id"`
	InstitutionId sql.NullInt64  `db:"institution_id"`
	Bio           sql.NullString `db:"bio"`
	LastLoginAt   sql.NullInt64  `db:"last_login_at"`
	CreatedAt     int64          `db:"created_at"`
	UpdatedAt     int64          `db:"updated_at"`
}

// 用户表名
const userTableName = "users"

// 插入用户
func (m *UserModel) Insert(ctx context.Context, user *User) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (username, password, email, phone, nickname, avatar, role, status, real_name, gender, birthday, school, grade, class, parent_id, institution_id, bio, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", userTableName)
	
	result, err := m.conn.ExecCtx(ctx, query,
		user.Username, user.Password, user.Email, user.Phone, user.Nickname, user.Avatar, user.Role, user.Status,
		user.RealName, user.Gender, user.Birthday, user.School, user.Grade, user.Class, user.ParentId, user.InstitutionId,
		user.Bio, user.CreatedAt, user.UpdatedAt)
	
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

// 根据ID查询用户
func (m *UserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	query := fmt.Sprintf("SELECT id, username, password, email, phone, nickname, avatar, role, status, real_name, gender, birthday, school, grade, class, parent_id, institution_id, bio, last_login_at, created_at, updated_at FROM %s WHERE id = ?", userTableName)
	
	var user User
	err := m.conn.QueryRowCtx(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// 根据用户名查询用户
func (m *UserModel) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := fmt.Sprintf("SELECT id, username, password, email, phone, nickname, avatar, role, status, real_name, gender, birthday, school, grade, class, parent_id, institution_id, bio, last_login_at, created_at, updated_at FROM %s WHERE username = ?", userTableName)
	
	var user User
	err := m.conn.QueryRowCtx(ctx, &user, query, username)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// 根据邮箱查询用户
func (m *UserModel) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := fmt.Sprintf("SELECT id, username, password, email, phone, nickname, avatar, role, status, real_name, gender, birthday, school, grade, class, parent_id, institution_id, bio, last_login_at, created_at, updated_at FROM %s WHERE email = ?", userTableName)
	
	var user User
	err := m.conn.QueryRowCtx(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// 更新用户信息
func (m *UserModel) Update(ctx context.Context, user *User) error {
	query := fmt.Sprintf("UPDATE %s SET nickname = ?, avatar = ?, real_name = ?, gender = ?, birthday = ?, school = ?, grade = ?, class = ?, bio = ?, updated_at = ? WHERE id = ?", userTableName)
	
	_, err := m.conn.ExecCtx(ctx, query,
		user.Nickname, user.Avatar, user.RealName, user.Gender, user.Birthday,
		user.School, user.Grade, user.Class, user.Bio, user.UpdatedAt, user.ID)
	
	return err
}

// 更新用户密码
func (m *UserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password = ?, updated_at = ? WHERE id = ?", userTableName)
	
	_, err := m.conn.ExecCtx(ctx, query, password, time.Now().Unix(), id)
	return err
}

// 更新用户状态
func (m *UserModel) UpdateStatus(ctx context.Context, id int64, status int) error {
	query := fmt.Sprintf("UPDATE %s SET status = ?, updated_at = ? WHERE id = ?", userTableName)
	
	_, err := m.conn.ExecCtx(ctx, query, status, time.Now().Unix(), id)
	return err
}

// 更新最后登录时间
func (m *UserModel) UpdateLastLogin(ctx context.Context, id int64) error {
	query := fmt.Sprintf("UPDATE %s SET last_login_at = ? WHERE id = ?", userTableName)
	
	_, err := m.conn.ExecCtx(ctx, query, time.Now().Unix(), id)
	return err
}

// 删除用户
func (m *UserModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", userTableName)
	
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// 查询用户列表
func (m *UserModel) FindList(ctx context.Context, req *types.UserListRequest) ([]*User, int64, error) {
	var whereClause []string
	var args []interface{}
	
	if req.Role != "" {
		whereClause = append(whereClause, "role = ?")
		args = append(args, req.Role)
	}
	
	if req.Status != 0 {
		whereClause = append(whereClause, "status = ?")
		args = append(args, req.Status)
	}
	
	if req.Keyword != "" {
		whereClause = append(whereClause, "(username LIKE ? OR nickname LIKE ? OR email LIKE ?)")
		keyword := "%" + req.Keyword + "%"
		args = append(args, keyword, keyword, keyword)
	}
	
	whereSQL := ""
	if len(whereClause) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClause, " AND ")
	}
	
	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", userTableName, whereSQL)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	// 查询列表
	offset := (req.Page - 1) * req.PageSize
	listQuery := fmt.Sprintf("SELECT id, username, password, email, phone, nickname, avatar, role, status, real_name, gender, birthday, school, grade, class, parent_id, institution_id, bio, last_login_at, created_at, updated_at FROM %s %s ORDER BY created_at DESC LIMIT ? OFFSET ?", userTableName, whereSQL)
	
	args = append(args, req.PageSize, offset)
	rows, err := m.conn.QueryCtx(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone, &user.Nickname, &user.Avatar,
			&user.Role, &user.Status, &user.RealName, &user.Gender, &user.Birthday, &user.School, &user.Grade,
			&user.Class, &user.ParentId, &user.InstitutionId, &user.Bio, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}
	
	return users, total, nil
}

// 查询家长的孩子列表
func (m *UserModel) FindChildrenByParentId(ctx context.Context, parentId int64) ([]*User, error) {
	query := fmt.Sprintf("SELECT id, username, password, email, phone, nickname, avatar, role, status, real_name, gender, birthday, school, grade, class, parent_id, institution_id, bio, last_login_at, created_at, updated_at FROM %s WHERE parent_id = ? AND role = 'student' AND status = 1", userTableName)
	
	rows, err := m.conn.QueryCtx(ctx, query, parentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var children []*User
	for rows.Next() {
		var child User
		err := rows.Scan(
			&child.ID, &child.Username, &child.Password, &child.Email, &child.Phone, &child.Nickname, &child.Avatar,
			&child.Role, &child.Status, &child.RealName, &child.Gender, &child.Birthday, &child.School, &child.Grade,
			&child.Class, &child.ParentId, &child.InstitutionId, &child.Bio, &child.LastLoginAt, &child.CreatedAt, &child.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		children = append(children, &child)
	}
	
	return children, nil
}

// 查询机构的教师列表
func (m *UserModel) FindTeachersByInstitutionId(ctx context.Context, institutionId int64) ([]*User, error) {
	query := fmt.Sprintf("SELECT id, username, password, email, phone, nickname, avatar, role, status, real_name, gender, birthday, school, grade, class, parent_id, institution_id, bio, last_login_at, created_at, updated_at FROM %s WHERE institution_id = ? AND role = 'teacher' AND status = 1", userTableName)
	
	rows, err := m.conn.QueryCtx(ctx, query, institutionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var teachers []*User
	for rows.Next() {
		var teacher User
		err := rows.Scan(
			&teacher.ID, &teacher.Username, &teacher.Password, &teacher.Email, &teacher.Phone, &teacher.Nickname, &teacher.Avatar,
			&teacher.Role, &teacher.Status, &teacher.RealName, &teacher.Gender, &teacher.Birthday, &teacher.School, &teacher.Grade,
			&teacher.Class, &teacher.ParentId, &teacher.InstitutionId, &teacher.Bio, &teacher.LastLoginAt, &teacher.CreatedAt, &teacher.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, &teacher)
	}
	
	return teachers, nil
}

// 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 转换为UserInfo类型
func (u *User) ToUserInfo() *types.UserInfo {
	return &types.UserInfo{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		Phone:       u.Phone.String,
		Nickname:    u.Nickname,
		Avatar:      u.Avatar.String,
		Role:        u.Role,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt.Int64,
	}
}

// 转换为UserProfile类型
func (u *User) ToUserProfile() *types.UserProfile {
	profile := &types.UserProfile{
		UserInfo:     *u.ToUserInfo(),
		RealName:     u.RealName.String,
		Gender:       int(u.Gender.Int64),
		School:       u.School.String,
		Grade:        u.Grade.String,
		Class:        u.Class.String,
		Bio:          u.Bio.String,
	}
	
	if u.ParentId.Valid {
		profile.ParentId = u.ParentId.Int64
	}
	
	if u.InstitutionId.Valid {
		profile.InstitutionId = u.InstitutionId.Int64
	}
	
	if u.Birthday.Valid {
		profile.Birthday = u.Birthday.Time.Format("2006-01-02")
	}
	
	return profile
}