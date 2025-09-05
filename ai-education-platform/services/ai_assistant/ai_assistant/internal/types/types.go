package types

// 学科枚举
type Subject struct {
	Subject string `json:"subject" options:"[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
}

// AI讲解请求
type AIExplainRequest struct {
	Subject      string `json:"subject" validate:"required,options=[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade        string `json:"grade" validate:"required"`
	Topic        string `json:"topic" validate:"required,min=1,max=200"`
	Difficulty   int    `json:"difficulty" validate:"min=1,max=3"` // 1:简单 2:中等 3:困难
	Language     string `json:"language" validate:"required,options=[zh,en]"` // 语言
	Context      string `json:"context,omitempty"` // 背景信息
	Style        string `json:"style" validate:"required,options=[simple,detailed,interactive]"` // 讲解风格
}

// AI讲解响应
type AIExplainResponse struct {
	ExplanationId string     `json:"explanation_id"`
	Content       string     `json:"content"` // 讲解内容
	Summary       string     `json:"summary"` // 内容摘要
	KeyPoints     []string   `json:"key_points"` // 关键点
	Examples      []Example  `json:"examples"` // 示例
	Resources     []Resource `json:"resources"` // 相关资源
	EstimatedTime int        `json:"estimated_time"` // 预估学习时间(分钟)
	CreatedAt     int64      `json:"created_at"`
}

// 示例
type Example struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Solution    string `json:"solution"`
	Difficulty  int    `json:"difficulty"`
}

// 学习资源
type Resource struct {
	Type        string `json:"type" options:"[video,article,exercise,book]"` // 资源类型
	Title       string `json:"title"`
	Url         string `json:"url"` // 资源链接
	Description string `json:"description"`
	Duration    int    `json:"duration,omitempty"` // 时长(分钟)
}

// AI问答请求
type AIChatRequest struct {
	Question      string        `json:"question" validate:"required,min=1,max=500"`
	Subject       string        `json:"subject,omitempty" options:"[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade         string        `json:"grade,omitempty"`
	Context       string        `json:"context,omitempty"` // 上下文信息
	SessionId     string        `json:"session_id,omitempty"` // 会话ID
	History       []ChatMessage `json:"history,omitempty"` // 历史对话
}

// AI问答响应
type AIChatResponse struct {
	AnswerId      string   `json:"answer_id"`
	Answer        string   `json:"answer"` // 回答内容
	Confidence    float64  `json:"confidence"` // 置信度 0-1
	Suggestions   []string `json:"suggestions"` // 建议的后续问题
	RelatedTopics []string `json:"related_topics"` // 相关知识点
	SessionId     string   `json:"session_id"`
	CreatedAt     int64    `json:"created_at"`
}

// 聊天消息
type ChatMessage struct {
	Role    string `json:"role" options:"[user,assistant]"` // 角色
	Content string `json:"content"` // 消息内容
	Time    int64  `json:"time"` // 时间戳
}

// 个性化学习路径请求
type LearningPathRequest struct {
	UserId       int64    `json:"user_id" validate:"required"`
	Subject      string   `json:"subject" validate:"required,options=[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade        string   `json:"grade" validate:"required"`
	Goal         string   `json:"goal" validate:"required"` // 学习目标
	CurrentLevel int      `json:"current_level" validate:"min=1,max=10"` // 当前水平
	WeakPoints   []string `json:"weak_points"` // 薄弱知识点
	StrongPoints []string `json:"strong_points"` // 强项知识点
	StudyTime    int      `json:"study_time" validate:"min=1"` // 每日学习时间(分钟)
}

// 个性化学习路径响应
type LearningPathResponse struct {
	PathId         string            `json:"path_id"`
	Title          string            `json:"title"` // 路径标题
	Description    string            `json:"description"` // 路径描述
	Duration       int               `json:"duration"` // 预估完成时间(天)
	Steps          []LearningStep    `json:"steps"` // 学习步骤
	Milestones     []Milestone       `json:"milestones"` // 里程碑
	Recommendations []string         `json:"recommendations"` // 学习建议
	CreatedAt      int64             `json:"created_at"`
}

// 学习步骤
type LearningStep struct {
	StepId       string     `json:"step_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Type         string     `json:"type" options:"[learn,practice,review,test]"` // 步骤类型
	Content      string     `json:"content"` // 学习内容
	Duration     int        `json:"duration"` // 预估时长(分钟)
	Difficulty   int        `json:"difficulty"`
	Prerequisites []string   `json:"prerequisites"` // 前置条件
	Resources    []Resource `json:"resources"` // 相关资源
}

// 里程碑
type Milestone struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetDate  int64  `json:"target_date"` // 目标完成时间
	Criteria    string `json:"criteria"` // 完成标准
}

// 学习建议请求
type StudyRecommendationRequest struct {
	UserId           int64                    `json:"user_id" validate:"required"`
	Subject          string                   `json:"subject,omitempty"`
	RecentActivity   []string                 `json:"recent_activity,omitempty"` // 最近学习活动
	PerformanceData  map[string]interface{}   `json:"performance_data,omitempty"` // 学习表现数据
	Preferences      map[string]interface{}   `json:"preferences,omitempty"` // 学习偏好
}

// 学习建议响应
type StudyRecommendationResponse struct {
	Recommendations []Recommendation `json:"recommendations"`
	GeneratedAt     int64            `json:"generated_at"`
}

// 建议
type Recommendation struct {
	Type        string                 `json:"type" options:"[course,knowledge_point,exercise,study_method]"` // 建议类型
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Priority    int                    `json:"priority"` // 优先级 1-5
	Reason      string                 `json:"reason"` // 推荐理由
	Action      string                 `json:"action"` // 建议操作
	Metadata    map[string]interface{} `json:"metadata"` // 额外信息
}

// 作业批改请求
type HomeworkGradingRequest struct {
	Question        string `json:"question" validate:"required"`
	StudentAnswer   string `json:"student_answer" validate:"required"`
	Subject         string `json:"subject" validate:"required,options=[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade           string `json:"grade"`
	QuestionType    string `json:"question_type" validate:"required,options=[multiple_choice,fill_blank,essay,calculation]"` // 题目类型
	ExpectedAnswer  string `json:"expected_answer,omitempty"` // 期望答案
	GradingCriteria string `json:"grading_criteria,omitempty"` // 批改标准
}

// 作业批改响应
type HomeworkGradingResponse struct {
	GradeId        string   `json:"grade_id"`
	Score          float64  `json:"score"` // 得分
	MaxScore       float64  `json:"max_score"` // 满分
	Feedback       string   `json:"feedback"` // 反馈意见
	Strengths      []string `json:"strengths"` // 优点
	Weaknesses     []string `json:"weaknesses"` // 不足
	Suggestions    []string `json:"suggestions"` // 改进建议
	CorrectAnswer  string   `json:"correct_answer,omitempty"` // 正确答案
	Explanation    string   `json:"explanation"` // 解析
	Confidence     float64  `json:"confidence"` // 批改置信度
	CreatedAt      int64    `json:"created_at"`
}

// 语音识别请求
type SpeechRecognitionRequest struct {
	AudioData     string `json:"audio_data" validate:"required"` // Base64编码的音频数据
	Language      string `json:"language" validate:"required,options=[zh,en]"` // 语言
	Subject       string `json:"subject,omitempty"` // 学科
	ExpectedText  string `json:"expected_text,omitempty"` // 期望文本(用于口语练习)
}

// 语音识别响应
type SpeechRecognitionResponse struct {
	RecognitionId      string  `json:"recognition_id"`
	Text               string  `json:"text"` // 识别文本
	Confidence         float64 `json:"confidence"` // 识别置信度
	PronunciationScore float64 `json:"pronunciation_score,omitempty"` // 发音评分(0-100)
	FluencyScore       float64 `json:"fluency_score,omitempty"` // 流利度评分(0-100)
	Feedback           string  `json:"feedback"` // 反馈
	CreatedAt          int64   `json:"created_at"`
}

// 通用响应
type CommonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}