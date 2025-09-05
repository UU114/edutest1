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

type AIChatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAIChatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AIChatLogic {
	return &AIChatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AIChatLogic) AIChat(req *types.AIChatRequest) (*types.AIChatResponse, error) {
	// 从上下文获取用户信息
	userId := l.ctx.Value("user_id").(int64)
	
	// 生成会话ID（如果未提供）
	sessionId := req.SessionId
	if sessionId == "" {
		sessionId = fmt.Sprintf("session_%d_%d", userId, time.Now().Unix())
	}
	
	// 调用AI服务获取回答
	aiResponse, err := l.callAIChatService(req)
	if err != nil {
		l.Logger.Errorf("AI聊天服务调用失败: %v", err)
		return nil, fmt.Errorf("AI服务暂时不可用")
	}
	
	// 保存聊天记录
	now := time.Now().Unix()
	answerId := fmt.Sprintf("answer_%d_%d", userId, time.Now().Unix())
	
	chat := &models.AIChat{
		AnswerId:      answerId,
		UserId:        userId,
		SessionId:     sessionId,
		Question:      req.Question,
		Answer:        aiResponse.Answer,
		Subject:       req.Subject,
		Grade:         req.Grade,
		Context:       req.Context,
		Confidence:    aiResponse.Confidence,
		CreatedAt:     now,
	}
	
	// 转换JSON字段
	suggestionsJSON, _ := json.Marshal(aiResponse.Suggestions)
	chat.Suggestions = string(suggestionsJSON)
	
	topicsJSON, _ := json.Marshal(aiResponse.RelatedTopics)
	chat.RelatedTopics = string(topicsJSON)
	
	historyJSON, _ := json.Marshal(req.History)
	chat.History = string(historyJSON)
	
	_, err = l.svcCtx.AIModel.InsertAIChat(l.ctx, chat)
	if err != nil {
		l.Logger.Errorf("保存AI聊天记录失败: %v", err)
		// 不影响返回结果
	}
	
	l.Logger.Infof("用户 %d AI聊天: %s", userId, req.Question)
	
	return &types.AIChatResponse{
		AnswerId:      answerId,
		Answer:        aiResponse.Answer,
		Confidence:    aiResponse.Confidence,
		Suggestions:   aiResponse.Suggestions,
		RelatedTopics: aiResponse.RelatedTopics,
		SessionId:     sessionId,
		CreatedAt:     now,
	}, nil
}

// AI聊天服务响应结构
type AIChatServiceResponse struct {
	Answer        string   `json:"answer"`
	Confidence    float64  `json:"confidence"`
	Suggestions   []string `json:"suggestions"`
	RelatedTopics []string `json:"related_topics"`
}

// 调用AI聊天服务
func (l *AIChatLogic) callAIChatService(req *types.AIChatRequest) (*AIChatServiceResponse, error) {
	// 构建聊天提示词
	prompt := l.buildChatPrompt(req)
	
	// 这里应该调用实际的AI服务
	// 为了演示，我们返回模拟数据
	return l.generateMockChatResponse(req), nil
}

// 构建聊天提示词
func (l *AIChatLogic) buildChatPrompt(req *types.AIChatRequest) string {
	prompt := fmt.Sprintf("学生提问：%s\n", req.Question)
	
	if req.Subject != "" {
		prompt += fmt.Sprintf("学科：%s\n", req.Subject)
	}
	
	if req.Grade != "" {
		prompt += fmt.Sprintf("年级：%s\n", req.Grade)
	}
	
	if req.Context != "" {
		prompt += fmt.Sprintf("背景信息：%s\n", req.Context)
	}
	
	if len(req.History) > 0 {
		prompt += "\n对话历史：\n"
		for _, msg := range req.History {
			prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
		}
	}
	
	prompt += "\n请以教师的身份回答这个问题，要求：\n"
	prompt += "1. 回答准确、易懂\n"
	prompt += "2. 适合学生理解水平\n"
	prompt += "3. 提供详细的解释和思路\n"
	prompt += "4. 必要时给出相关知识点和扩展建议"
	
	return prompt
}

// 生成模拟聊天响应
func (l *AIChatLogic) generateMockChatResponse(req *types.AIChatRequest) *AIChatServiceResponse {
	question := req.Question
	
	// 根据问题内容生成不同的回答
	var answer string
	var suggestions []string
	var relatedTopics []string
	
	if contains(question, []string{"公式", "定理", "证明"}) {
		answer = fmt.Sprintf("关于'%s'的回答：\n\n", question)
		answer += "这是一个很好的问题！让我来详细解释一下：\n\n"
		answer += "1. **基本概念**：首先理解相关的定义和基本概念\n"
		answer += "2. **推导过程**：了解公式或定理的推导思路和方法\n"
		answer += "3. **应用条件**：明确适用范围和限制条件\n"
		answer += "4. **实例分析**：通过具体例子加深理解\n\n"
		answer += "建议多做相关练习，加深对这个知识点的掌握。"
		
		suggestions = []string{
			"你能给我一个具体的例子吗？",
			"这个知识点在实际中如何应用？",
			"有没有相关的练习题推荐？",
		}
		
		relatedTopics = []string{"基础概念", "推导方法", "应用技巧", "相关定理"}
		
	} else if contains(question, []string{"理解", "意思", "解释"}) {
		answer = fmt.Sprintf("对于'%s'的理解：\n\n", question)
		answer += "这个问题的核心在于理解概念的本质和内涵。\n\n"
		answer += "**理解要点**：\n"
		answer += "- 从定义出发，把握核心特征\n"
		answer += "- 联系已知知识，建立知识网络\n"
		answer += "- 通过对比和类比，加深理解\n"
		answer += "- 结合实际应用，体会其价值\n\n"
		answer += "如果还有不理解的地方，欢迎继续提问！"
		
		suggestions = []string{
			"能举个简单的例子吗？",
			"这个和之前学的有什么关系？",
			"如何判断自己是否真正理解了？",
		}
		
		relatedTopics = []string{"概念理解", "知识联系", "应用实践", "学习方法"}
		
	} else {
		answer = fmt.Sprintf("感谢你的提问：'%s'\n\n", question)
		answer += "这是一个很有价值的问题。让我来为你分析：\n\n"
		answer += "**分析思路**：\n"
		answer += "1. 明确问题的关键点\n"
		answer += "2. 分析涉及的知识点\n"
		answer += "3. 提供解决方案\n"
		answer += "4. 总结方法和技巧\n\n"
		answer += "通过这样的思考方式，可以帮助你更好地理解和解决问题。"
		
		suggestions = []string{
			"这个问题还有其他解法吗？",
			"常见的错误有哪些？",
			"如何避免类似的错误？",
		}
		
		relatedTopics = []string{"问题分析", "解题方法", "错误分析", "技巧总结"}
	}
	
	return &AIChatServiceResponse{
		Answer:        answer,
		Confidence:    0.85 + (float64(time.Now().UnixNano()%1000))/1000*0.1, // 0.85-0.95之间的随机值
		Suggestions:   suggestions,
		RelatedTopics: relatedTopics,
	}
}

// 检查字符串是否包含关键词
func contains(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if len(text) >= len(keyword) {
			// 简单的包含检查
			for i := 0; i <= len(text)-len(keyword); i++ {
				if text[i:i+len(keyword)] == keyword {
					return true
				}
			}
		}
	}
	return false
}