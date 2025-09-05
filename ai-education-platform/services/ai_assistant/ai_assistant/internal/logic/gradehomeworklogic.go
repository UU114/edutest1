package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/models"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/svc"
	"ai-education-platform/services/ai_assistant/ai_assistant/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GradeHomeworkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGradeHomeworkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GradeHomeworkLogic {
	return &GradeHomeworkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GradeHomeworkLogic) GradeHomework(req *types.HomeworkGradingRequest) (*types.HomeworkGradingResponse, error) {
	// 从上下文获取用户信息
	userId := l.ctx.Value("user_id").(int64)
	
	// 调用AI服务进行作业批改
	aiResponse, err := l.callAIGradingService(req)
	if err != nil {
		l.Logger.Errorf("AI作业批改服务调用失败: %v", err)
		return nil, fmt.Errorf("AI批改服务暂时不可用")
	}
	
	// 保存批改记录
	now := time.Now().Unix()
	gradeId := fmt.Sprintf("grade_%d_%d", userId, time.Now().Unix())
	
	grade := &models.HomeworkGrade{
		GradeId:         gradeId,
		UserId:          userId,
		Question:        req.Question,
		StudentAnswer:   req.StudentAnswer,
		Subject:         req.Subject,
		Grade:           req.Grade,
		QuestionType:    req.QuestionType,
		ExpectedAnswer:  req.ExpectedAnswer,
		GradingCriteria: req.GradingCriteria,
		Score:           aiResponse.Score,
		MaxScore:        aiResponse.MaxScore,
		Feedback:        aiResponse.Feedback,
		CorrectAnswer:   aiResponse.CorrectAnswer,
		Explanation:     aiResponse.Explanation,
		Confidence:      aiResponse.Confidence,
		CreatedAt:       now,
	}
	
	// 转换JSON字段
	strengthsJSON, _ := json.Marshal(aiResponse.Strengths)
	grade.Strengths = string(strengthsJSON)
	
	weaknessesJSON, _ := json.Marshal(aiResponse.Weaknesses)
	grade.Weaknesses = string(weaknessesJSON)
	
	suggestionsJSON, _ := json.Marshal(aiResponse.Suggestions)
	grade.Suggestions = string(suggestionsJSON)
	
	_, err = l.svcCtx.AIModel.InsertHomeworkGrade(l.ctx, grade)
	if err != nil {
		l.Logger.Errorf("保存作业批改记录失败: %v", err)
		// 不影响返回结果
	}
	
	l.Logger.Infof("用户 %d 提交作业批改: %s - %s", userId, req.Subject, req.QuestionType)
	
	return &types.HomeworkGradingResponse{
		GradeId:       gradeId,
		Score:         aiResponse.Score,
		MaxScore:      aiResponse.MaxScore,
		Feedback:      aiResponse.Feedback,
		Strengths:     aiResponse.Strengths,
		Weaknesses:    aiResponse.Weaknesses,
		Suggestions:   aiResponse.Suggestions,
		CorrectAnswer: aiResponse.CorrectAnswer,
		Explanation:   aiResponse.Explanation,
		Confidence:    aiResponse.Confidence,
		CreatedAt:     now,
	}, nil
}

// AI作业批改服务响应结构
type AIHomeworkGradingServiceResponse struct {
	Score         float64  `json:"score"`
	MaxScore      float64  `json:"max_score"`
	Feedback      string   `json:"feedback"`
	Strengths     []string `json:"strengths"`
	Weaknesses    []string `json:"weaknesses"`
	Suggestions   []string `json:"suggestions"`
	CorrectAnswer string   `json:"correct_answer"`
	Explanation   string   `json:"explanation"`
	Confidence    float64  `json:"confidence"`
}

// 调用AI作业批改服务
func (l *GradeHomeworkLogic) callAIGradingService(req *types.HomeworkGradingRequest) (*AIHomeworkGradingServiceResponse, error) {
	// 构建批改提示词
	prompt := l.buildGradingPrompt(req)
	
	// 这里应该调用实际的AI服务
	// 为了演示，我们返回模拟数据
	return l.generateMockGradingResponse(req), nil
}

// 构建批改提示词
func (l *GradeHomeworkLogic) buildGradingPrompt(req *types.HomeworkGradingRequest) string {
	prompt := fmt.Sprintf(`请对学生作业进行智能批改：

题目：%s
学生答案：%s
学科：%s
年级：%s
题目类型：%s`, req.Question, req.StudentAnswer, req.Subject, req.Grade, req.QuestionType)
	
	if req.ExpectedAnswer != "" {
		prompt += fmt.Sprintf("\n标准答案：%s", req.ExpectedAnswer)
	}
	
	if req.GradingCriteria != "" {
		prompt += fmt.Sprintf("\n批改标准：%s", req.GradingCriteria)
	}
	
	prompt += `

请提供：
1. 得分和满分
2. 详细的反馈意见
3. 答案的优点
4. 答案的不足之处
5. 改进建议
6. 正确答案（如有）
7. 解析说明
8. 批改置信度

要求：
- 批改要客观公正
- 反馈要具体详细
- 既要指出问题，也要给予鼓励
- 提供有建设性的改进建议`
	
	return prompt
}

// 生成模拟批改响应
func (l *GradeHomeworkLogic) generateMockGradingResponse(req *types.HomeworkGradingRequest) *AIHomeworkGradingServiceResponse {
	// 根据题目类型和学生答案生成不同的批改结果
	var score float64
	var maxScore float64 = 100
	var feedback string
	var strengths []string
	var weaknesses []string
	var suggestions []string
	var correctAnswer string
	var explanation string
	var confidence float64
	
	switch req.QuestionType {
	case "multiple_choice":
		score, feedback, strengths, weaknesses, suggestions = l.gradeMultipleChoice(req)
		correctAnswer = "A" // 模拟标准答案
		explanation = "这道题考察了对基本概念的理解。正确答案是A，因为..."
		confidence = 0.95
		
	case "fill_blank":
		score, feedback, strengths, weaknesses, suggestions = l.gradeFillBlank(req)
		correctAnswer = "正确答案内容" // 模拟标准答案
		explanation = "填空题需要准确记忆和理解。正确答案是..."
		confidence = 0.90
		
	case "essay":
		score, feedback, strengths, weaknesses, suggestions = l.gradeEssay(req)
		correctAnswer = "参考答案要点"
		explanation = "作文题主要考察语言表达能力和思维逻辑。评分要点包括..."
		confidence = 0.80
		
	case "calculation":
		score, feedback, strengths, weaknesses, suggestions = l.gradeCalculation(req)
		correctAnswer = "42" // 模拟计算结果
		explanation = "计算题需要准确应用公式和步骤。正确解法是..."
		confidence = 0.92
		
	default:
		score = 75.0
		feedback = "答案基本正确，但还有提升空间。"
		strengths = []string{"基本理解题意"}
		weaknesses = []string{"某些细节处理不够准确"}
		suggestions = []string{"注意检查计算过程"}
		correctAnswer = "标准答案"
		explanation = "这道题的综合解答..."
		confidence = 0.85
	}
	
	return &AIHomeworkGradingServiceResponse{
		Score:         score,
		MaxScore:      maxScore,
		Feedback:      feedback,
		Strengths:     strengths,
		Weaknesses:    weaknesses,
		Suggestions:   suggestions,
		CorrectAnswer: correctAnswer,
		Explanation:   explanation,
		Confidence:    confidence,
	}
}

// 批改选择题
func (l *GradeHomeworkLogic) gradeMultipleChoice(req *types.HomeworkGradingRequest) (float64, string, []string, []string, []string) {
	// 模拟选择题批改
	answerLen := len(req.StudentAnswer)
	if answerLen == 0 {
		return 0.0, "未作答", []string{}, []string{"未选择答案"}, []string{"请仔细阅读题目并选择正确答案"}
	}
	
	// 简单的评分逻辑（实际应该对比标准答案）
	if strings.Contains(strings.ToUpper(req.StudentAnswer), "A") {
		return 100.0, "回答正确！很好地理解了题目要求。", 
			[]string{"答案选择正确", "理解题意准确"}, 
			[]string{}, 
			[]string{"继续保持，注意其他相关知识点"}
	}
	
	return 50.0, "答案不正确，请重新思考题目要求。", 
		[]string{"尝试作答"}, 
		[]string{"答案选择错误", "对知识点理解不够准确"}, 
		[]string{"复习相关知识点", "仔细分析题目"}
}

// 批改填空题
func (l *GradeHomeworkLogic) gradeFillBlank(req *types.HomeworkGradingRequest) (float64, string, []string, []string, []string) {
	answerLen := len(req.StudentAnswer)
	if answerLen == 0 {
		return 0.0, "未填写答案", []string{}, []string{"未填写内容"}, []string{"请根据题目要求填写答案"}
	}
	
	// 根据答案长度和质量评分
	if answerLen > 10 {
		return 85.0, "答案基本正确，内容较为完整。", 
			[]string{"内容充实", "理解题意"}, 
			[]string{"某些表述可以更准确"}, 
			[]string{"注意用词准确性", "可以更简洁表达"}
	}
	
	return 60.0, "答案过于简单，需要更加详细和准确。", 
		[]string{"基本理解要求"}, 
		[]string{"内容不够详细", "准确性有待提高"}, 
		[]string{"提供更详细的解答", "注意关键要点"}
}

// 批改作文题
func (l *GradeHomeworkLogic) gradeEssay(req *types.HomeworkGradingRequest) (float64, string, []string, []string, []string) {
	answerLen := len(req.StudentAnswer)
	if answerLen == 0 {
		return 0.0, "未提交作文", []string{}, []string{"未写作文"}, []string{"请按要求完成作文写作"}
	}
	
	// 根据作文长度和质量评分
	if answerLen > 200 {
		return 88.0, "作文内容充实，表达较为流畅。", 
			[]string{"内容丰富", "结构清晰", "语言表达较好"}, 
			[]string{"某些语句可以更精炼", "逻辑性可以加强"}, 
			[]string{"注意语言精炼", "加强逻辑训练", "多读优秀范文"}
	}
	
	return 70.0, "作文基本完成，但内容不够充实。", 
		[]string{"完成写作任务", "基本符合要求"}, 
		[]string{"内容不够详细", "语言表达有待提高"}, 
		[]string{"增加内容深度", "提高语言表达能力"}
}

// 批改计算题
func (l *GradeHomeworkLogic) gradeCalculation(req *types.HomeworkGradingRequest) (float64, string, []string, []string, []string) {
	if len(req.StudentAnswer) == 0 {
		return 0.0, "未进行计算", []string{}, []string{"未提供计算过程"}, []string{"请完成计算并写出解答过程"}
	}
	
	// 检查是否包含数字（模拟计算过程检查）
	if containsNumber(req.StudentAnswer) {
		return 92.0, "计算过程完整，结果正确。", 
			[]string{"计算步骤清晰", "结果准确", "方法正确"}, 
			[]string{}, 
			[]string{"继续保持，注意计算的准确性"}
	}
	
	return 65.0, "计算过程不完整，需要更加详细的步骤。", 
		[]string{"尝试进行计算"}, 
		[]string{"计算步骤不够完整", "结果准确性有待提高"}, 
		[]string{"写出详细计算过程", "检查每一步的计算结果"}
}

// 检查字符串是否包含数字
func containsNumber(s string) bool {
	for _, char := range s {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}