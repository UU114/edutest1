package logic

import (
	"context"
	"time"

	"ai-education-platform/services/course/course/internal/models"
	"ai-education-platform/services/course/course/internal/svc"
	"ai-education-platform/services/course/course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCourseLogic {
	return &CreateCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCourseLogic) CreateCourse(req *types.CreateCourseRequest) (*types.CommonResponse, error) {
	// 从上下文获取用户信息
	userId := l.ctx.Value("user_id").(int64)
	userRole := l.ctx.Value("role").(string)
	
	// 验证用户权限（只有教师和机构可以创建课程）
	if userRole != "teacher" && userRole != "institution" {
		return &types.CommonResponse{
			Success: false,
			Message: "只有教师和机构可以创建课程",
		}, nil
	}
	
	// 创建课程对象
	now := time.Now().Unix()
	course := &models.Course{
		Title:         req.Title,
		Description:   req.Description,
		Subject:       req.Subject,
		Grade:         req.Grade,
		Difficulty:    req.Difficulty,
		Price:         req.Price,
		Duration:      req.Duration,
		TeacherId:     userId,
		InstitutionId: userId, // 简化处理，实际应该根据角色区分
		Status:        0, // 草稿状态
		StudentCount:  0,
		Rating:        0.0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	// 设置可选字段
	if req.CoverImage != "" {
		course.CoverImage.String = req.CoverImage
		course.CoverImage.Valid = true
	}
	
	if len(req.Objectives) > 0 {
		course.Objectives.String = req.Objectives
		course.Objectives.Valid = true
	}
	
	if len(req.Prerequisites) > 0 {
		course.Prerequisites.String = req.Prerequisites
		course.Prerequisites.Valid = true
	}
	
	// 插入数据库
	courseId, err := l.svcCtx.CourseModel.InsertCourse(l.ctx, course)
	if err != nil {
		l.Logger.Errorf("创建课程失败: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "创建课程失败",
		}, err
	}
	
	l.Logger.Infof("课程创建成功: %s (ID: %d)", req.Title, courseId)
	return &types.CommonResponse{
		Success: true,
		Message: "课程创建成功",
	}, nil
}