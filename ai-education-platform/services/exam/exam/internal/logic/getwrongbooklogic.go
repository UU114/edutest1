package logic

import (
	"context"
	"ai-education-platform/services/exam/exam/internal/models"
	"ai-education-platform/services/exam/exam/internal/svc"
	"ai-education-platform/services/exam/exam/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWrongBookLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWrongBookLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWrongBookLogic {
	return &GetWrongBookLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetWrongBookLogic) GetWrongBook(req *types.WrongBookRequest) (*types.WrongBookResponse, error) {
	// 从上下文获取用户信息
	userId := l.ctx.Value("user_id").(int64)
	
	// 构建错题查询条件
	query := &models.WrongQuestionQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Subject:  req.Subject,
		Grade:    req.Grade,
		UserId:   userId,
		Mastered: req.Mastered,
	}
	
	// 查询错题列表
	wrongQuestions, total, err := l.svcCtx.ExamModel.GetUserWrongQuestions(l.ctx, query)
	if err != nil {
		l.Logger.Errorf("获取错题本失败: %v", err)
		return nil, err
	}
	
	// 转换为API响应格式
	var wrongBookList []types.WrongQuestion
	for _, wq := range wrongQuestions {
		wrongQuestionInfo := l.convertWrongQuestionToInfo(wq)
		wrongBookList = append(wrongBookList, *wrongQuestionInfo)
	}
	
	l.Logger.Infof("获取错题本成功，总数: %d", total)
	
	return &types.WrongBookResponse{
		Total: total,
		List:  wrongBookList,
	}, nil
}

// 转换错题数据为API格式
func (l *GetWrongBookLogic) convertWrongQuestionToInfo(wq *models.WrongQuestion) *types.WrongQuestion {
	info := &types.WrongQuestion{
		ID:            wq.Id,
		UserId:        wq.UserId,
		QuestionId:    wq.QuestionId,
		QuestionTitle: wq.QuestionTitle,
		StudentAnswer: wq.StudentAnswer,
		CorrectAnswer: wq.CorrectAnswer,
		Subject:       wq.Subject,
		Grade:         wq.Grade,
		WrongCount:    wq.WrongCount,
		LastWrongTime: wq.LastWrongTime,
		Mastered:      wq.Mastered,
		CreatedAt:     wq.CreatedAt,
		UpdatedAt:     wq.UpdatedAt,
	}
	
	return info
}