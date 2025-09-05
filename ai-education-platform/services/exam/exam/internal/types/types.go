package types

// 题目类型枚举
type QuestionType struct {
	Type string `json:"type" options:"[multiple_choice,single_choice,fill_blank,essay,calculation,judgment]"`
}

// 学科枚举
type Subject struct {
	Subject string `json:"subject" options:"[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
}

// 题目信息
type QuestionInfo struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Type          string    `json:"type" options:"[multiple_choice,single_choice,fill_blank,essay,calculation,judgment]"`
	Subject       string    `json:"subject" options:"[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade         string    `json:"grade"`
	Difficulty    int       `json:"difficulty"` // 1-5 难度等级
	Content       string    `json:"content"` // 题目内容
	Options       []string  `json:"options,omitempty"` // 选择题选项
	CorrectAnswer string    `json:"correct_answer"` // 正确答案
	Analysis      string    `json:"analysis"` // 解析
	KnowledgePoints []int64  `json:"knowledge_points"` // 关联知识点ID
	Tags          []string  `json:"tags"` // 题目标签
	CreatorId     int64     `json:"creator_id"`
	Status        int       `json:"status"` // 0:草稿 1:发布 2:禁用
	CreatedAt     int64     `json:"created_at"`
	UpdatedAt     int64     `json:"updated_at"`
	UsageCount    int       `json:"usage_count"` // 使用次数
	CorrectRate   float64   `json:"correct_rate"` // 正确率
}

// 试卷信息
type ExamInfo struct {
	ID            int64           `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	Subject       string          `json:"subject"`
	Grade         string          `json:"grade"`
	Duration      int             `json:"duration"` // 考试时长(分钟)
	TotalScore    float64         `json:"total_score"` // 总分
	PassScore     float64         `json:"pass_score"` // 及格分
	QuestionCount int             `json:"question_count"` // 题目数量
	Questions     []ExamQuestion  `json:"questions"` // 试卷题目
	CreatorId     int64           `json:"creator_id"`
	Status        int             `json:"status"` // 0:草稿 1:发布 2:禁用
	CreatedAt     int64           `json:"created_at"`
	UpdatedAt     int64           `json:"updated_at"`
}

// 试卷题目
type ExamQuestion struct {
	QuestionId    int64    `json:"question_id"`
	QuestionType  string   `json:"question_type"`
	QuestionTitle string  `json:"question_title"`
	Score         float64  `json:"score"` // 分值
	Order         int      `json:"order"` // 题目顺序
	Content       string   `json:"content"` // 题目内容
	Options       []string `json:"options,omitempty"` // 选项
}

// 学生考试记录
type ExamRecord struct {
	ID            int64        `json:"id"`
	ExamId        int64        `json:"exam_id"`
	UserId        int64        `json:"user_id"`
	Score         float64      `json:"score"` // 得分
	TotalScore    float64      `json:"total_score"` // 总分
	Status        string       `json:"status" options:"[in_progress,completed,timeout,submitted]"` // 考试状态
	StartTime     int64        `json:"start_time"` // 开始时间
	EndTime       int64        `json:"end_time"` // 结束时间
	TimeUsed      int          `json:"time_used"` // 用时(秒)
	Answers       []ExamAnswer `json:"answers"` // 答题记录
	CreatedAt     int64        `json:"created_at"`
	UpdatedAt     int64        `json:"updated_at"`
}

// 考试答案
type ExamAnswer struct {
	QuestionId int64   `json:"question_id"`
	Answer     string  `json:"answer"` // 学生答案
	IsCorrect  bool    `json:"is_correct"` // 是否正确
	Score      float64 `json:"score"` // 得分
	TimeSpent  int     `json:"time_spent"` // 答题用时(秒)
}

// 错题记录
type WrongQuestion struct {
	ID            int64  `json:"id"`
	UserId        int64  `json:"user_id"`
	QuestionId    int64  `json:"question_id"`
	QuestionTitle string `json:"question_title"`
	StudentAnswer string `json:"student_answer"` // 学生答案
	CorrectAnswer string `json:"correct_answer"` // 正确答案
	Subject       string `json:"subject"`
	Grade         string `json:"grade"`
	WrongCount    int    `json:"wrong_count"` // 错误次数
	LastWrongTime int64  `json:"last_wrong_time"` // 最后错误时间
	Mastered      bool   `json:"mastered"` // 是否已掌握
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}

// 练习记录
type PracticeRecord struct {
	ID            int64   `json:"id"`
	UserId        int64   `json:"user_id"`
	Subject       string  `json:"subject"`
	QuestionCount int     `json:"question_count"` // 题目数量
	CorrectCount  int     `json:"correct_count"` // 正确数量
	Score         float64 `json:"score"` // 得分
	TimeUsed      int     `json:"time_used"` // 用时(秒)
	Difficulty    string  `json:"difficulty"` // 难度
	CreatedAt     int64   `json:"created_at"`
}

// 创建题目请求
type CreateQuestionRequest struct {
	Title         string   `json:"title" validate:"required,min=1,max=200"`
	Type          string   `json:"type" validate:"required,options=[multiple_choice,single_choice,fill_blank,essay,calculation,judgment]"`
	Subject       string   `json:"subject" validate:"required,options=[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade         string   `json:"grade" validate:"required"`
	Difficulty    int      `json:"difficulty" validate:"required,min=1,max=5"`
	Content       string   `json:"content" validate:"required,min=1"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer string   `json:"correct_answer" validate:"required"`
	Analysis      string   `json:"analysis"`
	KnowledgePoints []int64 `json:"knowledge_points"`
	Tags          []string `json:"tags"`
}

// 更新题目请求
type UpdateQuestionRequest struct {
	Title         string   `json:"title,omitempty"`
	Type          string   `json:"type,omitempty"`
	Subject       string   `json:"subject,omitempty"`
	Grade         string   `json:"grade,omitempty"`
	Difficulty    int      `json:"difficulty,omitempty"`
	Content       string   `json:"content,omitempty"`
	Options       []string `json:"options,omitempty"`
	CorrectAnswer string   `json:"correct_answer,omitempty"`
	Analysis      string   `json:"analysis,omitempty"`
	KnowledgePoints []int64 `json:"knowledge_points,omitempty"`
	Tags          []string `json:"tags,omitempty"`
}

// 题目列表查询请求
type QuestionListRequest struct {
	Page        int    `json:"page" validate:"required,min=1"`
	PageSize    int    `json:"page_size" validate:"required,min=1,max=100"`
	Subject     string `json:"subject,omitempty"`
	Grade       string `json:"grade,omitempty"`
	Type        string `json:"type,omitempty"`
	Difficulty  int    `json:"difficulty,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	Tags        string `json:"tags,omitempty"`
	CreatorId   int64  `json:"creator_id,omitempty"`
	Status      int    `json:"status,omitempty"`
}

// 题目列表响应
type QuestionListResponse struct {
	Total int64          `json:"total"`
	List  []QuestionInfo `json:"list"`
}

// 创建试卷请求
type CreateExamRequest struct {
	Title         string          `json:"title" validate:"required,min=1,max=100"`
	Description   string          `json:"description" validate:"required,min=1,max=500"`
	Subject       string          `json:"subject" validate:"required"`
	Grade         string          `json:"grade" validate:"required"`
	Duration      int             `json:"duration" validate:"required,min=1"`
	TotalScore    float64         `json:"total_score" validate:"required,min=1"`
	PassScore     float64         `json:"pass_score" validate:"required,min=0"`
	Questions     []ExamQuestion  `json:"questions" validate:"required,min=1"`
}

// 智能出题请求
type SmartGenerateRequest struct {
	Subject       string   `json:"subject" validate:"required"`
	Grade         string   `json:"grade" validate:"required"`
	QuestionTypes []string `json:"question_types,omitempty"` // 题目类型
	Difficulty    int      `json:"difficulty,omitempty"` // 难度
	KnowledgePoints []int64 `json:"knowledge_points,omitempty"` // 知识点
	QuestionCount int      `json:"question_count" validate:"required,min=1,max=100"`
	Tags          []string `json:"tags,omitempty"`
	ExcludeIds    []int64  `json:"exclude_ids,omitempty"` // 排除的题目ID
}

// 智能出题响应
type SmartGenerateResponse struct {
	Questions   []QuestionInfo `json:"questions"`
	GeneratedAt int64           `json:"generated_at"`
}

// 开始考试请求
type StartExamRequest struct {
	ExamId int64 `json:"exam_id" validate:"required"`
}

// 提交答案请求
type SubmitAnswerRequest struct {
	ExamId     int64  `json:"exam_id" validate:"required"`
	QuestionId int64  `json:"question_id" validate:"required"`
	Answer     string `json:"answer" validate:"required"`
	TimeSpent  int    `json:"time_spent"` // 答题用时(秒)
}

// 完成考试请求
type FinishExamRequest struct {
	ExamId int64 `json:"exam_id" validate:"required"`
}

// 获取错题本请求
type WrongBookRequest struct {
	Subject  string `json:"subject,omitempty"`
	Grade    string `json:"grade,omitempty"`
	Page     int    `json:"page" validate:"required,min=1"`
	PageSize int    `json:"page_size" validate:"required,min=1,max=100"`
	Mastered bool   `json:"mastered,omitempty"` // 是否只显示未掌握
}

// 错题本响应
type WrongBookResponse struct {
	Total int64           `json:"total"`
	List  []WrongQuestion `json:"list"`
}

// 练习请求
type PracticeRequest struct {
	Subject       string   `json:"subject" validate:"required"`
	Grade         string   `json:"grade" validate:"required"`
	QuestionTypes []string `json:"question_types,omitempty"`
	Difficulty    int      `json:"difficulty,omitempty"`
	KnowledgePoints []int64 `json:"knowledge_points,omitempty"`
	QuestionCount int      `json:"question_count" validate:"required,min=1,max=50"`
	Mode          string   `json:"mode" validate:"required,options=[practice,test]"` // 练习模式
	TimeLimit     int      `json:"time_limit,omitempty"` // 时间限制(秒)
}

// 练习响应
type PracticeResponse struct {
	PracticeId string         `json:"practice_id"`
	Questions  []QuestionInfo `json:"questions"`
	TimeLimit  int            `json:"time_limit"`
	StartedAt  int64          `json:"started_at"`
}

// 通用响应
type CommonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}