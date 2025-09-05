package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 生成随机字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// 验证邮箱格式
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// 验证手机号格式
func IsValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}

// 验证密码强度
func IsValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	
	// 至少包含字母和数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	return hasLetter && hasNumber
}

// 分页工具
type Pagination struct {
	Page     int `json:"page" validate:"required,min=1"`
	PageSize int `json:"page_size" validate:"required,min=1,max=100"`
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

func (p *Pagination) GetLimit() int {
	return p.PageSize
}

// 响应工具
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code,omitempty"`
}

func SuccessResponse(message string, data interface{}) *Response {
	return &Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string, code int) *Response {
	return &Response{
		Success: false,
		Message: message,
		Code:    code,
	}
}

// 分页响应
type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	TotalPages int64     `json:"total_pages"`
}

func NewPageResponse(list interface{}, total int64, page, pageSize int) *PageResponse {
	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}
	
	return &PageResponse{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// JWT工具
type JWTUtil struct {
	secretKey string
}

func NewJWTUtil(secretKey string) *JWTUtil {
	return &JWTUtil{secretKey: secretKey}
}

func (j *JWTUtil) GenerateToken(userId int64, username, role string, expire int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userId,
		"username": username,
		"role":     role,
		"exp":      expire,
		"iat":      time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTUtil) ParseToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}
	
	return nil, fmt.Errorf("invalid token")
}

// 时间工具
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

func ParseTimeToTimestamp(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

// 字符串工具
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func TruncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length] + "..."
}

// 数字工具
func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func StringToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func CalculatePercentage(value, total int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(value) / float64(total) * 100
}

// 验证工具
func ValidateRequiredFields(data map[string]interface{}, requiredFields []string) []string {
	var missingFields []string
	
	for _, field := range requiredFields {
		if val, exists := data[field]; !exists || val == "" || val == nil {
			missingFields = append(missingFields, field)
		}
	}
	
	return missingFields
}

// 安全工具
func SanitizeString(input string) string {
	// 移除潜在的XSS攻击字符
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")
	input = strings.ReplaceAll(input, "script", "")
	
	return input
}

func IsValidGrade(grade string) bool {
	validGrades := []string{
		"一年级", "二年级", "三年级", "四年级", "五年级", "六年级",
		"初一", "初二", "初三",
		"高一", "高二", "高三",
	}
	
	return Contains(validGrades, grade)
}

func IsValidSubject(subject string) bool {
	validSubjects := []string{
		"math", "chinese", "english", "physics", "chemistry", "biology", "history", "geography", "politics",
	}
	
	return Contains(validSubjects, subject)
}