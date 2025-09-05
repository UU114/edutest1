package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/models"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/svc"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AIExplainLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAIExplainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIExplainLogic {
	return &AIExplainLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AIExplainLogic) AIExplain(req *types.AIExplainRequest) (*types.AIExplainResponse, error) {
	// 从上下文获取用户信息
	userId := l.ctx.Value("user_id").(int64)
	
	// 生成唯一ID
	explanationId := fmt.Sprintf("exp_%d_%d", userId, time.Now().Unix())
	
	// 调用AI服务获取解释
	aiResponse, err := l.callAIService(req)
	if err != nil {
		l.Logger.Errorf("AI服务调用失败: %v", err)
		return nil, fmt.Errorf("AI服务暂时不可用")
	}
	
	// 保存解释记录
	now := time.Now().Unix()
	explanation := &models.AIExplanation{
		ExplanationId: explanationId,
		UserId:        userId,
		Subject:       req.Subject,
		Grade:         req.Grade,
		Topic:         req.Topic,
		Difficulty:    req.Difficulty,
		Language:      req.Language,
		Style:         req.Style,
		Content:       aiResponse.Content,
		Summary:       aiResponse.Summary,
		EstimatedTime: aiResponse.EstimatedTime,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	
	// 转换JSON字段
	keyPointsJSON, _ := json.Marshal(aiResponse.KeyPoints)
	explanation.KeyPoints = string(keyPointsJSON)
	
	examplesJSON, _ := json.Marshal(aiResponse.Examples)
	explanation.Examples = string(examplesJSON)
	
	resourcesJSON, _ := json.Marshal(aiResponse.Resources)
	explanation.Resources = string(resourcesJSON)
	
	_, err = l.svcCtx.AIModel.InsertAIExplanation(l.ctx, explanation)
	if err != nil {
		l.Logger.Errorf("保存AI解释记录失败: %v", err)
		// 不影响返回结果
	}
	
	l.Logger.Infof("用户 %d 请求AI解释: %s - %s", userId, req.Subject, req.Topic)
	
	return &types.AIExplainResponse{
		ExplanationId: explanationId,
		Content:       aiResponse.Content,
		Summary:       aiResponse.Summary,
		KeyPoints:     aiResponse.KeyPoints,
		Examples:      aiResponse.Examples,
		Resources:     aiResponse.Resources,
		EstimatedTime: aiResponse.EstimatedTime,
		CreatedAt:     now,
	}, nil
}

// AI服务响应结构
type AIExplainServiceResponse struct {
	Content       string            `json:"content"`
	Summary       string            `json:"summary"`
	KeyPoints     []string          `json:"key_points"`
	Examples      []types.Example   `json:"examples"`
	Resources     []types.Resource  `json:"resources"`
	EstimatedTime int               `json:"estimated_time"`
}

// 调用AI服务
func (l *AIExplainLogic) callAIService(req *types.AIExplainRequest) (*AIExplainServiceResponse, error) {
	// 构建AI请求提示词
	prompt := l.buildExplanationPrompt(req)
	
	// 这里应该调用实际的AI服务（如OpenAI、Claude等）
	// 为了演示，我们返回模拟数据
	return l.generateMockResponse(req), nil
}

// 构建解释提示词
func (l *AIExplainLogic) buildExplanationPrompt(req *types.AIExplainRequest) string {
	styleMap := map[string]string{
		"simple":      "简单易懂",
		"detailed":    "详细深入",
		"interactive": "互动式",
	}
	
	difficultyMap := map[string]int{
		1: "简单",
		2: "中等",
		3: "困难",
	}
	
	return fmt.Sprintf(`请为%s年级的学生讲解%s学科的"%s"知识点。
要求：
- 难度等级：%s
- 讲解风格：%s
- 使用语言：%s
- %s

请提供：
1. 详细的内容解释
2. 内容摘要
3. 关键知识点列表
4. 相关例题和解析
5. 推荐学习资源
6. 预估学习时间（分钟）`,
		req.Grade,
		req.Subject,
		req.Topic,
		difficultyMap[req.Difficulty],
		styleMap[req.Style],
		req.Language,
		req.Context)
}

// 生成模拟响应
func (l *AIExplainLogic) generateMockResponse(req *types.AIExplainRequest) *AIExplainServiceResponse {
	// 根据不同学科和难度生成不同的模拟内容
	content := fmt.Sprintf("关于%s的详细解释：\n\n", req.Topic)
	
	switch req.Subject {
	case "math":
		content += "这是一个数学概念，涉及到数理逻辑和计算方法。\n\n"
		content += "基本原理：\n"
		content += "1. 概念定义和基本性质\n"
		content += "2. 相关公式和定理\n"
		content += "3. 应用场景和实例分析\n\n"
		content += "学习建议：\n"
		content += "- 理解基本概念的重要性\n"
		content += "- 多做练习巩固知识点\n"
		content += "- 结合实际应用加深理解"
		
	case "chinese":
		content += "这是一个语文知识点，涉及语言文字的理解和运用。\n\n"
		content += "主要内容：\n"
		content += "1. 字词解析和用法\n"
		content += "2. 语法结构和特点\n"
		content += "3. 文化和历史背景\n\n"
		content += "学习方法：\n"
		content += "- 注重积累和记忆\n"
		content += "- 多读多写多练习\n"
		content += "- 理解文化内涵"
		
	default:
		content += "这是一个重要的知识点，需要系统性地学习和理解。\n\n"
		content += "学习要点：\n"
		content += "1. 基本概念和定义\n"
		content += "2. 核心原理和规律\n"
		content += "3. 实际应用和价值\n\n"
		content += "建议通过多种方式学习，包括理论学习、实践练习和讨论交流。"
	}
	
	return &AIExplainServiceResponse{
		Content: content,
		Summary: fmt.Sprintf("关于%s的核心知识点总结，包含了基本概念、原理和应用方法。", req.Topic),
		KeyPoints: []string{
			"基本概念理解",
			"核心原理掌握",
			"实际应用能力",
			"相关知识点联系",
		},
		Examples: []types.Example{
			{
				Title:       "基础例题",
				Description: "帮助理解基本概念",
				Solution:    "通过步骤化解析，掌握解题方法",
				Difficulty:  1,
			},
			{
				Title:       "进阶练习",
				Description: "巩固知识点应用",
				Solution:    "综合运用多个知识点解决问题",
				Difficulty:  2,
			},
		},
		Resources: []types.Resource{
			{
				Type:        "article",
				Title:       "相关学习资料",
				Url:         "https://example.com/resource1",
				Description: "详细的学习指导资料",
				Duration:    30,
			},
			{
				Type:        "video",
				Title:       "视频讲解",
				Url:         "https://example.com/video1",
				Description: "生动形象的视频教程",
				Duration:    15,
			},
		},
		EstimatedTime: 30 + req.Difficulty*10,
	}
}