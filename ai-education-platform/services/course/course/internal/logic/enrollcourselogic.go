package logic

import (
	"context"
	"time"

	"ai-education-platform/services/course/course/internal/models"
	"ai-education-platform/services/course/course/internal/svc"
	"ai-education-platform/services/course/course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnrollCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEnrollCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnrollCourseLogic {
	return &EnrollCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *EnrollCourseLogic) EnrollCourse(req *types.EnrollCourseRequest) (*types.CommonResponse, error) {
	// 从上下文获取用户信息
	userId := l.ctx.Value("user_id").(int64)
	
	// 验证课程是否存在
	course, err := l.svcCtx.CourseModel.FindCourseById(l.ctx, req.CourseId)
	if err != nil {
		l.Logger.Errorf("课程不存在: %d", req.CourseId)
		return &types.CommonResponse{
			Success: false,
			Message: "课程不存在",
		}, nil
	}
	
	// 验证课程状态
	if course.Status != 1 {
		return &types.CommonResponse{
			Success: false,
			Message: "课程未发布",
		}, nil
	}
	
	// 检查是否已经报名
	// 这里简化处理，实际应该查询报名记录
	
	// 创建报名记录
	now := time.Now().Unix()
	enrollment := &models.CourseEnrollment{
		CourseId:   req.CourseId,
		UserId:     userId,
		Status:     "active",
		EnrolledAt: now,
		CreatedAt:  now,
	}
	
	_, err = l.svcCtx.CourseModel.EnrollCourse(l.ctx, enrollment)
	if err != nil {
		l.Logger.Errorf("课程报名失败: %v", err)
		return &types.CommonResponse{
			Success: false,
			Message: "报名失败",
		}, err
	}
	
	// 创建学习进度记录
	progress := &models.StudentProgress{
		CourseId:        req.CourseId,
		UserId:          userId,
		CompletedLessons: 0,
		TotalLessons:    0, // 后续根据课程实际课时数更新
		Progress:        0.0,
		StudyTime:       0,
		StartedAt:       now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	
	_, err = l.svcCtx.CourseModel.InsertStudentProgress(l.ctx, progress)
	if err != nil {
		l.Logger.Errorf("创建学习进度失败: %v", err)
		// 不影响报名结果
	}
	
	// 更新课程学生数量
	course.StudentCount++
	course.UpdatedAt = now
	err = l.svcCtx.CourseModel.UpdateCourse(l.ctx, course)
	if err != nil {
		l.Logger.Errorf("更新课程学生数量失败: %v", err)
		// 不影响报名结果
	}
	
	l.Logger.Infof("用户 %d 报名课程 %d 成功", userId, req.CourseId)
	return &types.CommonResponse{
		Success: true,
		Message: "报名成功",
	}, nil
}