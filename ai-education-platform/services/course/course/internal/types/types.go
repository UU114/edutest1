package types

// 学科枚举
type Subject struct {
	Subject string `json:"subject" options="[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
}

// 课程信息
type CourseInfo struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Subject       string    `json:"subject" options="[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade         string    `json:"grade"` // 年级
	Difficulty    int       `json:"difficulty"` // 1:简单 2:中等 3:困难
	CoverImage    string    `json:"cover_image,omitempty"`
	Price         float64   `json:"price"` // 价格
	Duration      int       `json:"duration"` // 课程时长(分钟)
	TeacherId     int64     `json:"teacher_id"`
	TeacherName   string    `json:"teacher_name"`
	InstitutionId int64     `json:"institution_id"`
	Status        int       `json:"status"` // 0:草稿 1:发布 2:下架
	StudentCount  int       `json:"student_count"`
	Rating        float64   `json:"rating"` // 评分
	CreatedAt     int64     `json:"created_at"`
	UpdatedAt     int64     `json:"updated_at"`
}

// 课程详细信息
type CourseDetail struct {
	CourseInfo
	Chapters     []Chapter `json:"chapters"`
	Tags         []string  `json:"tags"`
	Objectives   string    `json:"objectives"` // 学习目标
	Prerequisites string   `json:"prerequisites"` // 先修要求
}

// 章节信息
type Chapter struct {
	ID          int64        `json:"id"`
	CourseId    int64        `json:"course_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Order       int          `json:"order"`
	Duration    int          `json:"duration"` // 章节时长(分钟)
	Lessons     []Lesson     `json:"lessons"`
}

// 课时信息
type Lesson struct {
	ID        int64  `json:"id"`
	ChapterId int64  `json:"chapter_id"`
	Title     string `json:"title"`
	Type      string `json:"type" options="[video,audio,document,exercise]"` // 课时类型
	Content   string `json:"content"` // 内容URL或文本
	Duration  int    `json:"duration"` // 课时时长(分钟)
	Order     int    `json:"order"`
	IsFree    bool   `json:"is_free"` // 是否免费
}

// 知识点信息
type KnowledgePoint struct {
	ID            int64     `json:"id"`
	CourseId      int64     `json:"course_id"`
	ChapterId     int64     `json:"chapter_id,omitempty"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Subject       string    `json:"subject"`
	Grade         string    `json:"grade"`
	Difficulty    int       `json:"difficulty"`
	Keywords      []string  `json:"keywords"`
	Prerequisites []int64   `json:"prerequisites"` // 前置知识点ID
	CreatedAt     int64     `json:"created_at"`
}

// 创建课程请求
type CreateCourseRequest struct {
	Title         string   `json:"title" validate:"required,min=1,max=100"`
	Description   string   `json:"description" validate:"required,min=1,max=1000"`
	Subject       string   `json:"subject" validate:"required,options=[math,chinese,english,physics,chemistry,biology,history,geography,politics]"`
	Grade         string   `json:"grade" validate:"required"`
	Difficulty    int      `json:"difficulty" validate:"required,min=1,max=3"`
	CoverImage    string   `json:"cover_image,omitempty"`
	Price         float64  `json:"price" validate:"min=0"`
	Duration      int      `json:"duration" validate:"required,min=1"`
	Tags          []string `json:"tags"`
	Objectives    string   `json:"objectives"`
	Prerequisites string   `json:"prerequisites"`
}

// 更新课程请求
type UpdateCourseRequest struct {
	Title         string   `json:"title,omitempty"`
	Description   string   `json:"description,omitempty"`
	Subject       string   `json:"subject,omitempty"`
	Grade         string   `json:"grade,omitempty"`
	Difficulty    int      `json:"difficulty,omitempty"`
	CoverImage    string   `json:"cover_image,omitempty"`
	Price         float64  `json:"price,omitempty"`
	Duration      int      `json:"duration,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Objectives    string   `json:"objectives,omitempty"`
	Prerequisites string   `json:"prerequisites,omitempty"`
}

// 课程列表查询请求
type CourseListRequest struct {
	Page        int    `json:"page" validate:"required,min=1"`
	PageSize    int    `json:"page_size" validate:"required,min=1,max=100"`
	Subject     string `json:"subject,omitempty"`
	Grade       string `json:"grade,omitempty"`
	Difficulty  int    `json:"difficulty,omitempty"`
	TeacherId   int64  `json:"teacher_id,omitempty"`
	InstitutionId int64 `json:"institution_id,omitempty"`
	Keyword     string `json:"keyword,omitempty"`
	Status      int    `json:"status,omitempty"`
	SortBy      string `json:"sort_by,omitempty" options="[created_at,rating,student_count,price]"`
	SortOrder   string `json:"sort_order,omitempty" options="[asc,desc]"`
}

// 课程列表响应
type CourseListResponse struct {
	Total int64        `json:"total"`
	List  []CourseInfo `json:"list"`
}

// 学习进度
type StudyProgress struct {
	CourseId        int64   `json:"course_id"`
	UserId          int64   `json:"user_id"`
	CompletedLessons int     `json:"completed_lessons"`
	TotalLessons    int     `json:"total_lessons"`
	Progress        float64 `json:"progress"` // 0-100
	LastLessonId    int64   `json:"last_lesson_id"`
	StudyTime       int     `json:"study_time"` // 学习时长(分钟)
	StartedAt       int64   `json:"started_at"`
	CompletedAt     int64   `json:"completed_at,omitempty"`
}

// 课程报名请求
type EnrollCourseRequest struct {
	CourseId int64 `json:"course_id" validate:"required"`
}

// 通用响应
type CommonResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}