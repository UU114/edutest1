package logic

import (
	"context"

	"ai-education-platform/services/course/course/internal/models"
	"ai-education-platform/services/course/course/internal/svc"
	"ai-education-platform/services/course/course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCourseDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCourseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCourseDetailLogic {
	return &GetCourseDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCourseDetailLogic) GetCourseDetail(courseId int64) (*types.CourseDetail, error) {
	// 查询课程信息
	course, err := l.svcCtx.CourseModel.FindCourseById(l.ctx, courseId)
	if err != nil {
		l.Logger.Errorf("查询课程详情失败: %v", err)
		return nil, err
	}
	
	// 查询章节信息
	chapters, err := l.svcCtx.CourseModel.FindChaptersByCourseId(l.ctx, courseId)
	if err != nil {
		l.Logger.Errorf("查询章节信息失败: %v", err)
		return nil, err
	}
	
	// 转换章节信息并查询课时
	chapterDetails := make([]*types.Chapter, 0, len(chapters))
	for _, chapter := range chapters {
		chapterDetail := chapter.ToChapter()
		
		// 查询章节下的课时
		lessons, err := l.svcCtx.CourseModel.FindLessonsByChapterId(l.ctx, chapter.ID)
		if err != nil {
			l.Logger.Errorf("查询课时信息失败: %v", err)
			return nil, err
		}
		
		// 转换课时信息
		lessonDetails := make([]*types.Lesson, 0, len(lessons))
		for _, lesson := range lessons {
			lessonDetails = append(lessonDetails, lesson.ToLesson())
		}
		
		chapterDetail.Lessons = lessonDetails
		chapterDetails = append(chapterDetails, chapterDetail)
	}
	
	return course.ToCourseDetail(chapterDetails), nil
}