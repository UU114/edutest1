package logic

import (
	"context"
	"fmt"
	"time"
	"ai-education-platform/services/exam/exam/internal/models"
	"ai-education-platform/services/exam/exam/internal/svc"
	"ai-education-platform/services/exam/exam/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SmartGenerateQuestionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSmartGenerateQuestionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmartGenerateQuestionsLogic {
	return &SmartGenerateQuestionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SmartGenerateQuestionsLogic) SmartGenerateQuestions(req *types.SmartGenerateRequest) (*types.SmartGenerateResponse, error) {
	// 构建智能出题查询条件
	query := &models.SmartGenerateQuery{
		Subject:       req.Subject,
		Grade:         req.Grade,
		QuestionTypes: req.QuestionTypes,
		Difficulty:    req.Difficulty,
		KnowledgePoints: req.KnowledgePoints,
		QuestionCount: req.QuestionCount,
		Tags:          req.Tags,
		ExcludeIds:    req.ExcludeIds,
	}
	
	// 智能生成题目
	questions, err := l.svcCtx.ExamModel.SmartGenerateQuestions(l.ctx, query)
	if err != nil {
		l.Logger.Errorf("智能生成题目失败: %v", err)
		return nil, err
	}
	
	// 转换为API响应格式
	var questionList []types.QuestionInfo
	for _, q := range questions {
		questionInfo := l.convertQuestionToInfo(q)
		questionList = append(questionList, *questionInfo)
	}
	
	l.Logger.Infof("智能生成题目成功，生成数量: %d", len(questionList))
	
	return &types.SmartGenerateResponse{
		Questions:   questionList,
		GeneratedAt: time.Now().Unix(),
	}, nil
}

// 转换题目数据为API格式
func (l *SmartGenerateQuestionsLogic) convertQuestionToInfo(q *models.Question) *types.QuestionInfo {
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