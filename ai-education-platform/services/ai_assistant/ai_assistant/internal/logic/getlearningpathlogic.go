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

type GetLearningPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLearningPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLearningPathLogic {
	return &GetLearningPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLearningPathLogic) GetLearningPath(req *types.LearningPathRequest) (*types.LearningPathResponse, error) {
	// 验证用户ID
	if req.UserId != l.ctx.Value("user_id").(int64) {
		return nil, fmt.Errorf("无权访问其他用户的学习路径")
	}
	
	// 生成学习路径ID
	pathId := fmt.Sprintf("path_%d_%d", req.UserId, time.Now().Unix())
	
	// 调用AI服务生成学习路径
	aiResponse, err := l.callAILearningPathService(req)
	if err != nil {
		l.Logger.Errorf("AI学习路径服务调用失败: %v", err)
		return nil, fmt.Errorf("AI服务暂时不可用")
	}
	
	// 保存学习路径记录
	now := time.Now().Unix()
	learningPath := &models.LearningPath{
		PathId:         pathId,
		UserId:         req.UserId,
		Subject:        req.Subject,
		Grade:          req.Grade,
		Goal:           req.Goal,
		CurrentLevel:   req.CurrentLevel,
		StudyTime:      req.StudyTime,
		Title:          aiResponse.Title,
		Description:    aiResponse.Description,
		Duration:       aiResponse.Duration,
		Status:         "active",
		Progress:       0.0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	
	// 转换JSON字段
	weakPointsJSON, _ := json.Marshal(req.WeakPoints)
	learningPath.WeakPoints = string(weakPointsJSON)
	
	strongPointsJSON, _ := json.Marshal(req.StrongPoints)
	learningPath.StrongPoints = string(strongPointsJSON)
	
	stepsJSON, _ := json.Marshal(aiResponse.Steps)
	learningPath.Steps = string(stepsJSON)
	
	milestonesJSON, _ := json.Marshal(aiResponse.Milestones)
	learningPath.Milestones = string(milestonesJSON)
	
	recommendationsJSON, _ := json.Marshal(aiResponse.Recommendations)
	learningPath.Recommendations = string(recommendationsJSON)
	
	_, err = l.svcCtx.AIModel.InsertLearningPath(l.ctx, learningPath)
	if err != nil {
		l.Logger.Errorf("保存学习路径记录失败: %v", err)
		// 不影响返回结果
	}
	
	l.Logger.Infof("用户 %d 生成学习路径: %s - %s", req.UserId, req.Subject, req.Goal)
	
	return &types.LearningPathResponse{
		PathId:         pathId,
		Title:          aiResponse.Title,
		Description:    aiResponse.Description,
		Duration:       aiResponse.Duration,
		Steps:          aiResponse.Steps,
		Milestones:     aiResponse.Milestones,
		Recommendations: aiResponse.Recommendations,
		CreatedAt:      now,
	}, nil
}

// AI学习路径服务响应结构
type AILearningPathServiceResponse struct {
	Title          string               `json:"title"`
	Description    string               `json:"description"`
	Duration       int                  `json:"duration"`
	Steps          []types.LearningStep `json:"steps"`
	Milestones     []types.Milestone    `json:"milestones"`
	Recommendations []string             `json:"recommendations"`
}

// 调用AI学习路径服务
func (l *GetLearningPathLogic) callAILearningPathService(req *types.LearningPathRequest) (*AILearningPathServiceResponse, error) {
	// 构建学习路径提示词
	prompt := l.buildLearningPathPrompt(req)
	
	// 这里应该调用实际的AI服务
	// 为了演示，我们返回模拟数据
	return l.generateMockLearningPathResponse(req), nil
}

// 构建学习路径提示词
func (l *GetLearningPathLogic) buildLearningPathPrompt(req *types.LearningPathRequest) string {
	levelMap := map[string]string{
		"1": "初学者",
		"2": "基础",
		"3": "入门",
		"4": "初级",
		"5": "中级",
		"6": "中高级",
		"7": "高级",
		"8": "熟练",
		"9": "精通",
		"10": "专家",
	}
	
	prompt := fmt.Sprintf(`为学生制定个性化学习路径：

学生信息：
- 年级：%s
- 学科：%s
- 当前水平：%s（1-10级）
- 学习目标：%s
- 每日学习时间：%d分钟

已知情况：`, req.Grade, req.Subject, levelMap[fmt.Sprintf("%d", req.CurrentLevel)], req.Goal, req.StudyTime)
	
	if len(req.WeakPoints) > 0 {
		prompt += fmt.Sprintf("\n- 薄弱知识点：%v", req.WeakPoints)
	}
	
	if len(req.StrongPoints) > 0 {
		prompt += fmt.Sprintf("\n- 强项知识点：%v", req.StrongPoints)
	}
	
	prompt += `

请制定一个系统性的学习路径，包括：
1. 路径标题和描述
2. 预估完成时间（天）
3. 详细的学习步骤（每个步骤包含目标、内容、时长、难度等）
4. 关键里程碑节点
5. 学习建议和注意事项

要求：
- 根据学生当前水平调整难度
- 优先弥补薄弱环节
- 发挥学生优势
- 合理安排学习进度
- 提供具体可执行的计划`
	
	return prompt
}

// 生成模拟学习路径响应
func (l *GetLearningPathLogic) generateMockLearningPathResponse(req *types.LearningPathRequest) *AILearningPathServiceResponse {
	// 根据学科和目标生成不同的学习路径
	title := fmt.Sprintf("%s%s学习路径", req.Grade, req.Subject)
	description := fmt.Sprintf("专为%s年级学生定制的%s学习计划，目标：%s", req.Grade, req.Subject, req.Goal)
	
	// 根据当前水平和目标计算预估时间
	baseDuration := 30 // 基础时间
	levelDiff := req.CurrentLevel - 5 // 以5级为基准
	if levelDiff < 0 {
		baseDuration += (-levelDiff) * 10 // 水平低需要更多时间
	} else {
		baseDuration -= levelDiff * 5 // 水平高可以缩短时间
	}
	
	if baseDuration < 15 {
		baseDuration = 15
	} else if baseDuration > 90 {
		baseDuration = 90
	}
	
	// 生成学习步骤
	steps := l.generateLearningSteps(req)
	
	// 生成里程碑
	milestones := l.generateMilestones(req, baseDuration)
	
	// 生成学习建议
	recommendations := l.generateRecommendations(req)
	
	return &AILearningPathServiceResponse{
		Title:          title,
		Description:    description,
		Duration:       baseDuration,
		Steps:          steps,
		Milestones:     milestones,
		Recommendations: recommendations,
	}
}

// 生成学习步骤
func (l *GetLearningPathLogic) generateLearningSteps(req *types.LearningPathRequest) []types.LearningStep {
	steps := []types.LearningStep{}
	
	// 基础知识学习阶段
	if req.CurrentLevel <= 3 {
		steps = append(steps, types.LearningStep{
			StepId:       "step_1",
			Title:        "基础知识学习",
			Description:  "掌握基本概念和原理",
			Type:         "learn",
			Content:      "系统学习基础知识，建立完整的知识框架",
			Duration:     120,
			Difficulty:   1,
			Prerequisites: []string{},
			Resources: []types.Resource{
				{
					Type:        "video",
					Title:       "基础知识视频讲解",
					Url:         "https://example.com/basic-video",
					Description: "生动形象的基础知识讲解",
					Duration:    60,
				},
			},
		})
	}
	
	// 薄弱环节加强
	if len(req.WeakPoints) > 0 {
		steps = append(steps, types.LearningStep{
			StepId:       "step_2",
			Title:        "薄弱环节加强",
			Description:  "针对薄弱知识点进行专项训练",
			Type:         "practice",
			Content:      fmt.Sprintf("重点加强：%v", req.WeakPoints),
			Duration:     90,
			Difficulty:   2,
			Prerequisites: []string{"step_1"},
			Resources: []types.Resource{
				{
					Type:        "exercise",
					Title:       "专项练习题",
					Url:         "https://example.com/weakness-exercises",
					Description: "针对性练习题目",
					Duration:    45,
				},
			},
		})
	}
	
	// 核心技能提升
	steps = append(steps, types.LearningStep{
		StepId:       "step_3",
		Title:        "核心技能提升",
		Description:  "掌握核心解题方法和技巧",
		Type:         "learn",
		Content:      "学习核心知识点和解题技巧",
		Duration:     150,
		Difficulty:   3,
		Prerequisites: []string{"step_1"},
		Resources: []types.Resource{
			{
				Type:        "article",
				Title:       "核心技能指南",
				Url:         "https://example.com/core-skills",
				Description: "核心技能详细指导",
				Duration:    75,
			},
		},
	})
	
	// 综合练习
	steps = append(steps, types.LearningStep{
		StepId:       "step_4",
		Title:        "综合练习",
		Description:  "通过综合题目检验学习效果",
		Type:         "practice",
		Content:      "综合运用所学知识解决问题",
		Duration:     120,
		Difficulty:   3,
		Prerequisites: []string{"step_2", "step_3"},
		Resources: []types.Resource{
			{
				Type:        "exercise",
				Title:       "综合练习题集",
				Url:         "https://example.com/comprehensive-exercises",
				Description: "综合性练习题目",
				Duration:    90,
			},
		},
	})
	
	// 模拟测试
	steps = append(steps, types.LearningStep{
		StepId:       "step_5",
		Title:        "模拟测试",
		Description:  "通过模拟测试评估学习成果",
		Type:         "test",
		Content:      "完成模拟测试，检验学习效果",
		Duration:     90,
		Difficulty:   3,
		Prerequisites: []string{"step_4"},
		Resources: []types.Resource{
			{
				Type:        "exercise",
				Title:       "模拟测试题",
				Url:         "https://example.com/mock-test",
				Description: "模拟测试题目",
				Duration:    60,
			},
		},
	})
	
	return steps
}

// 生成里程碑
func (l *GetLearningPathLogic) generateMilestones(req *types.LearningPathRequest, duration int) []types.Milestone {
	milestones := []types.Milestone{}
	
	// 第一周里程碑
	milestones = append(milestones, types.Milestone{
		Title:       "基础阶段完成",
		Description: "完成基础知识学习，掌握基本概念",
		TargetDate:  time.Now().AddDate(0, 0, duration/5).Unix(),
		Criteria:    "完成基础练习题，正确率达到80%以上",
	})
	
	// 中期里程碑
	milestones = append(milestones, types.Milestone{
		Title:       "技能提升阶段",
		Description: "掌握核心技能，能够解决中等难度问题",
		TargetDate:  time.Now().AddDate(0, 0, duration*2/5).Unix(),
		Criteria:    "完成综合练习，正确率达到75%以上",
	})
	
	// 后期里程碑
	milestones = append(milestones, types.Milestone{
		Title:       "综合应用阶段",
		Description: "能够综合运用知识解决复杂问题",
		TargetDate:  time.Now().AddDate(0, 0, duration*3/5).Unix(),
		Criteria:    "完成模拟测试，成绩达到预期目标",
	})
	
	// 最终里程碑
	milestones = append(milestones, types.Milestone{
		Title:       "学习目标达成",
		Description: "达到预期的学习目标",
		TargetDate:  time.Now().AddDate(0, 0, duration).Unix(),
		Criteria:    req.Goal,
	})
	
	return milestones
}

// 生成学习建议
func (l *GetLearningPathLogic) generateRecommendations(req *types.LearningPathRequest) []string {
	recommendations := []string{
		"坚持每日学习，保持学习连续性",
		"做好学习笔记，及时总结归纳",
		"多做练习，在实践中巩固知识",
		"遇到问题及时提问，不要积压",
		"定期复习，防止遗忘",
	}
	
	// 根据每日学习时间调整建议
	if req.StudyTime < 30 {
		recommendations = append(recommendations, "建议增加每日学习时间，至少保证30分钟")
	} else if req.StudyTime > 120 {
		recommendations = append(recommendations, "注意劳逸结合，避免过度疲劳")
	}
	
	// 根据薄弱环节给出建议
	if len(req.WeakPoints) > 0 {
		recommendations = append(recommendations, fmt.Sprintf("重点加强薄弱环节：%v", req.WeakPoints))
	}
	
	// 根据当前水平给出建议
	if req.CurrentLevel <= 3 {
		recommendations = append(recommendations, "建议从基础开始，循序渐进")
	} else if req.CurrentLevel >= 8 {
		recommendations = append(recommendations, "可以尝试更高难度的挑战")
	}
	
	return recommendations
}