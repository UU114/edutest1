package logic

import (
	"context"

	"ai-education-platform/services/course/course/internal/models"
	"ai-education-platform/services/course/course/internal/svc"
	"ai-education-platform/services/course/course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCourseListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseListLogic {
	return &GetCourseListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseListLogic) GetCourseList(req *types.CourseListRequest) (*types.CourseListResponse, error) {
	// 查询课程列表
	courses, total, err := l.svcCtx.CourseModel.FindCourseList(l.ctx, req)
	if err != nil {
		l.Logger.Errorf("查询课程列表失败: %v", err)
		return nil, err
	}
	
	// 转换为CourseInfo格式
	courseInfos := make([]*types.CourseInfo, 0, len(courses))
	for _, course := range courses {
		courseInfos = append(courseInfos, course.ToCourseInfo())
	}
	
	return &types.CourseListResponse{
		Total: total,
		List:  courseInfos,
	}, nil
}