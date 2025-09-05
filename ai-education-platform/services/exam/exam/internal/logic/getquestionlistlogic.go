package logic

import (
	"context"
	"ai-education-platform/services/exam/exam/internal/models"
	"ai-education-platform/services/exam/exam/internal/svc"
	"ai-education-platform/services/exam/exam/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetQuestionListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetQuestionListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetQuestionListLogic {
	return &GetQuestionListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetQuestionListLogic) GetQuestionList(req *types.QuestionListRequest) (*types.QuestionListResponse, error) {
	// 构建查询条件
	query := &models.QuestionListQuery{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Subject:    req.Subject,
		Grade:      req.Grade,
		Type:       req.Type,
		Difficulty: req.Difficulty,
		Keyword:    req.Keyword,
		Tags:       req.Tags,
		CreatorId:  req.CreatorId,
		Status:     req.Status,
	}
	
	// 查询题目列表
	questions, total, err := l.svcCtx.ExamModel.GetQuestionList(l.ctx, query)
	if err != nil {
		l.Logger.Errorf("获取题目列表失败: %v", err)
		return nil, err
	}
	
	// 转换为API响应格式
	var questionList []types.QuestionInfo
	for _, q := range questions {
		questionInfo := l.convertQuestionToInfo(q)
		questionList = append(questionList, *questionInfo)
	}
	
	l.Logger.Infof("获取题目列表成功，总数: %d", total)
	
	return &types.QuestionListResponse{
		Total: total,
		List:  questionList,
	}, nil
}

// 转换题目数据为API格式
func (l *GetQuestionListLogic) convertQuestionToInfo(q *models.Question) *types.QuestionInfo {
	info := &types.QuestionInfo{
		ID:            q.Id,
		Title:         q.Title,
		Type:          q.Type,
		Subject:       q.Subject,
		Grade:         q.Grade,
		Difficulty:    q.Difficulty,
		Content:       q.Content,
		CorrectAnswer: q.CorrectAnswer,
		Analysis:      q.Analysis,
		CreatorId:     q.CreatorId,
		Status:        q.Status,
		CreatedAt:     q.CreatedAt,
		UpdatedAt:     q.UpdatedAt,
		UsageCount:    q.UsageCount,
		CorrectRate:   q.CorrectRate,
	}
	
	// 转换JSON字段
	if q.Options != "" {
		options, _ := models.JSONToStringSlice(q.Options)
		info.Options = options
	}
	
	if q.KnowledgePoints != "" {
		knowledgePoints, _ := models.JSONToInt64Slice(q.KnowledgePoints)
		info.KnowledgePoints = knowledgePoints
	}
	
	if q.Tags != "" {
		tags, _ := models.JSONToStringSlice(q.Tags)
		info.Tags = tags
	}
	
	return info
}