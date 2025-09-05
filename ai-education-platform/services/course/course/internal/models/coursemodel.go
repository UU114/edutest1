package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"ai-education-platform/services/course/course/internal/types"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CourseModel struct {
	conn sqlx.SqlConn
}

func NewCourseModel(conn sqlx.SqlConn) *CourseModel {
	return &CourseModel{
		conn: conn,
	}
}

// 课程数据结构
type Course struct {
	ID            int64          `db:"id"`
	Title         string         `db:"title"`
	Description   string         `db:"description"`
	Subject       string         `db:"subject"`
	Grade         string         `db:"grade"`
	Difficulty    int            `db:"difficulty"`
	CoverImage    sql.NullString `db:"cover_image"`
	Price         float64        `db:"price"`
	Duration      int            `db:"duration"`
	TeacherId     int64          `db:"teacher_id"`
	InstitutionId int64          `db:"institution_id"`
	Status        int            `db:"status"`
	StudentCount  int            `db:"student_count"`
	Rating        float64        `db:"rating"`
	Objectives    sql.NullString `db:"objectives"`
	Prerequisites sql.NullString `db:"prerequisites"`
	CreatedAt     int64          `db:"created_at"`
	UpdatedAt     int64          `db:"updated_at"`
}

// 章节数据结构
type Chapter struct {
	ID          int64  `db:"id"`
	CourseId    int64  `db:"course_id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Order       int    `db:"order"`
	Duration    int    `db:"duration"`
	CreatedAt   int64  `db:"created_at"`
	UpdatedAt   int64  `db:"updated_at"`
}

// 课时数据结构
type Lesson struct {
	ID         int64  `db:"id"`
	ChapterId  int64  `db:"chapter_id"`
	Title      string `db:"title"`
	Type       string `db:"type"`
	Content    string `db:"content"`
	Duration   int    `db:"duration"`
	Order      int    `db:"order"`
	IsFree     bool   `db:"is_free"`
	CreatedAt  int64  `db:"created_at"`
	UpdatedAt  int64  `db:"updated_at"`
}

// 知识点数据结构
type KnowledgePoint struct {
	ID            int64          `db:"id"`
	CourseId      int64          `db:"course_id"`
	ChapterId     sql.NullInt64  `db:"chapter_id"`
	Title         string         `db:"title"`
	Description   string         `db:"description"`
	Subject       string         `db:"subject"`
	Grade         string         `db:"grade"`
	Difficulty    int            `db:"difficulty"`
	Keywords      sql.NullString `db:"keywords"`
	Prerequisites sql.NullString `db:"prerequisites"`
	CreatedAt     int64          `db:"created_at"`
}

// 学习进度数据结构
type StudentProgress struct {
	ID              int64   `db:"id"`
	CourseId        int64   `db:"course_id"`
	UserId          int64   `db:"user_id"`
	CompletedLessons int     `db:"completed_lessons"`
	TotalLessons    int     `db:"total_lessons"`
	Progress        float64 `db:"progress"`
	LastLessonId    int64   `db:"last_lesson_id"`
	StudyTime       int     `db:"study_time"`
	StartedAt       int64   `db:"started_at"`
	CompletedAt     int64   `db:"completed_at"`
	CreatedAt       int64   `db:"created_at"`
	UpdatedAt       int64   `db:"updated_at"`
}

// 课程报名数据结构
type CourseEnrollment struct {
	ID          int64     `db:"id"`
	CourseId    int64     `db:"course_id"`
	UserId      int64     `db:"user_id"`
	Status      string    `db:"status"`
	EnrolledAt  int64     `db:"enrolled_at"`
	CompletedAt int64     `db:"completed_at"`
	CreatedAt   int64     `db:"created_at"`
}

const (
	courseTableName        = "courses"
	chapterTableName       = "chapters"
	lessonTableName       = "lessons"
	knowledgePointTableName = "knowledge_points"
	studentProgressTableName = "student_progress"
	courseEnrollmentTableName = "course_enrollments"
)

// 课程相关操作
func (m *CourseModel) InsertCourse(ctx context.Context, course *Course) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, description, subject, grade, difficulty, cover_image, price, duration, teacher_id, institution_id, status, student_count, rating, objectives, prerequisites, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", courseTableName)
	
	result, err := m.conn.ExecCtx(ctx, query,
		course.Title, course.Description, course.Subject, course.Grade, course.Difficulty,
		course.CoverImage, course.Price, course.Duration, course.TeacherId, course.InstitutionId,
		course.Status, course.StudentCount, course.Rating, course.Objectives, course.Prerequisites,
		course.CreatedAt, course.UpdatedAt)
	
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

func (m *CourseModel) FindCourseById(ctx context.Context, id int64) (*Course, error) {
	query := fmt.Sprintf("SELECT id, title, description, subject, grade, difficulty, cover_image, price, duration, teacher_id, institution_id, status, student_count, rating, objectives, prerequisites, created_at, updated_at FROM %s WHERE id = ?", courseTableName)
	
	var course Course
	err := m.conn.QueryRowCtx(ctx, &course, query, id)
	if err != nil {
		return nil, err
	}
	
	return &course, nil
}

func (m *CourseModel) UpdateCourse(ctx context.Context, course *Course) error {
	query := fmt.Sprintf("UPDATE %s SET title = ?, description = ?, subject = ?, grade = ?, difficulty = ?, cover_image = ?, price = ?, duration = ?, status = ?, student_count = ?, rating = ?, objectives = ?, prerequisites = ?, updated_at = ? WHERE id = ?", courseTableName)
	
	_, err := m.conn.ExecCtx(ctx, query,
		course.Title, course.Description, course.Subject, course.Grade, course.Difficulty,
		course.CoverImage, course.Price, course.Duration, course.Status, course.StudentCount,
		course.Rating, course.Objectives, course.Prerequisites, course.UpdatedAt, course.ID)
	
	return err
}

func (m *CourseModel) DeleteCourse(ctx context.Context, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", courseTableName)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *CourseModel) FindCourseList(ctx context.Context, req *types.CourseListRequest) ([]*Course, int64, error) {
	var whereClause []string
	var args []interface{}
	
	if req.Subject != "" {
		whereClause = append(whereClause, "subject = ?")
		args = append(args, req.Subject)
	}
	
	if req.Grade != "" {
		whereClause = append(whereClause, "grade = ?")
		args = append(args, req.Grade)
	}
	
	if req.Difficulty != 0 {
		whereClause = append(whereClause, "difficulty = ?")
		args = append(args, req.Difficulty)
	}
	
	if req.TeacherId != 0 {
		whereClause = append(whereClause, "teacher_id = ?")
		args = append(args, req.TeacherId)
	}
	
	if req.InstitutionId != 0 {
		whereClause = append(whereClause, "institution_id = ?")
		args = append(args, req.InstitutionId)
	}
	
	if req.Status != 0 {
		whereClause = append(whereClause, "status = ?")
		args = append(args, req.Status)
	}
	
	if req.Keyword != "" {
		whereClause = append(whereClause, "(title LIKE ? OR description LIKE ?)")
		keyword := "%" + req.Keyword + "%"
		args = append(args, keyword, keyword)
	}
	
	whereSQL := ""
	if len(whereClause) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClause, " AND ")
	}
	
	// 排序
	orderBy := "created_at DESC"
	if req.SortBy != "" {
		sortOrder := "DESC"
		if req.SortOrder == "asc" {
			sortOrder = "ASC"
		}
		orderBy = req.SortBy + " " + sortOrder
	}
	
	// 查询总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", courseTableName, whereSQL)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	// 查询列表
	offset := (req.Page - 1) * req.PageSize
	listQuery := fmt.Sprintf("SELECT id, title, description, subject, grade, difficulty, cover_image, price, duration, teacher_id, institution_id, status, student_count, rating, objectives, prerequisites, created_at, updated_at FROM %s %s ORDER BY %s LIMIT ? OFFSET ?", courseTableName, whereSQL, orderBy)
	
	args = append(args, req.PageSize, offset)
	rows, err := m.conn.QueryCtx(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	
	var courses []*Course
	for rows.Next() {
		var course Course
		err := rows.Scan(
			&course.ID, &course.Title, &course.Description, &course.Subject, &course.Grade, &course.Difficulty,
			&course.CoverImage, &course.Price, &course.Duration, &course.TeacherId, &course.InstitutionId,
			&course.Status, &course.StudentCount, &course.Rating, &course.Objectives, &course.Prerequisites,
			&course.CreatedAt, &course.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		courses = append(courses, &course)
	}
	
	return courses, total, nil
}

// 章节相关操作
func (m *CourseModel) InsertChapter(ctx context.Context, chapter *Chapter) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (course_id, title, description, `order`, duration, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)", chapterTableName)
	
	result, err := m.conn.ExecCtx(ctx, query,
		chapter.CourseId, chapter.Title, chapter.Description, chapter.Order, chapter.Duration,
		chapter.CreatedAt, chapter.UpdatedAt)
	
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

func (m *CourseModel) FindChaptersByCourseId(ctx context.Context, courseId int64) ([]*Chapter, error) {
	query := fmt.Sprintf("SELECT id, course_id, title, description, `order`, duration, created_at, updated_at FROM %s WHERE course_id = ? ORDER BY `order`", chapterTableName)
	
	rows, err := m.conn.QueryCtx(ctx, query, courseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var chapters []*Chapter
	for rows.Next() {
		var chapter Chapter
		err := rows.Scan(&chapter.ID, &chapter.CourseId, &chapter.Title, &chapter.Description, &chapter.Order, &chapter.Duration, &chapter.CreatedAt, &chapter.UpdatedAt)
		if err != nil {
			return nil, err
		}
		chapters = append(chapters, &chapter)
	}
	
	return chapters, nil
}

// 课时相关操作
func (m *CourseModel) FindLessonsByChapterId(ctx context.Context, chapterId int64) ([]*Lesson, error) {
	query := fmt.Sprintf("SELECT id, chapter_id, title, type, content, duration, `order`, is_free, created_at, updated_at FROM %s WHERE chapter_id = ? ORDER BY `order`", lessonTableName)
	
	rows, err := m.conn.QueryCtx(ctx, query, chapterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var lessons []*Lesson
	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(&lesson.ID, &lesson.ChapterId, &lesson.Title, &lesson.Type, &lesson.Content, &lesson.Duration, &lesson.Order, &lesson.IsFree, &lesson.CreatedAt, &lesson.UpdatedAt)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, &lesson)
	}
	
	return lessons, nil
}

// 学习进度相关操作
func (m *CourseModel) FindStudentProgress(ctx context.Context, courseId, userId int64) (*StudentProgress, error) {
	query := fmt.Sprintf("SELECT id, course_id, user_id, completed_lessons, total_lessons, progress, last_lesson_id, study_time, started_at, completed_at, created_at, updated_at FROM %s WHERE course_id = ? AND user_id = ?", studentProgressTableName)
	
	var progress StudentProgress
	err := m.conn.QueryRowCtx(ctx, &progress, query, courseId, userId)
	if err != nil {
		return nil, err
	}
	
	return &progress, nil
}

func (m *CourseModel) UpdateStudentProgress(ctx context.Context, progress *StudentProgress) error {
	query := fmt.Sprintf("UPDATE %s SET completed_lessons = ?, total_lessons = ?, progress = ?, last_lesson_id = ?, study_time = ?, completed_at = ?, updated_at = ? WHERE course_id = ? AND user_id = ?", studentProgressTableName)
	
	_, err := m.conn.ExecCtx(ctx, query,
		progress.CompletedLessons, progress.TotalLessons, progress.Progress, progress.LastLessonId,
		progress.StudyTime, progress.CompletedAt, progress.UpdatedAt, progress.CourseId, progress.UserId)
	
	return err
}

func (m *CourseModel) InsertStudentProgress(ctx context.Context, progress *StudentProgress) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (course_id, user_id, completed_lessons, total_lessons, progress, last_lesson_id, study_time, started_at, completed_at, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", studentProgressTableName)
	
	result, err := m.conn.ExecCtx(ctx, query,
		progress.CourseId, progress.UserId, progress.CompletedLessons, progress.TotalLessons, progress.Progress,
		progress.LastLessonId, progress.StudyTime, progress.StartedAt, progress.CompletedAt, progress.CreatedAt, progress.UpdatedAt)
	
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

// 课程报名相关操作
func (m *CourseModel) EnrollCourse(ctx context.Context, enrollment *CourseEnrollment) (int64, error) {
	query := fmt.Sprintf("INSERT INTO %s (course_id, user_id, status, enrolled_at, completed_at, created_at) VALUES (?, ?, ?, ?, ?, ?)", courseEnrollmentTableName)
	
	result, err := m.conn.ExecCtx(ctx, query,
		enrollment.CourseId, enrollment.UserId, enrollment.Status, enrollment.EnrolledAt,
		enrollment.CompletedAt, enrollment.CreatedAt)
	
	if err != nil {
		return 0, err
	}
	
	return result.LastInsertId()
}

func (m *CourseModel) FindUserEnrollments(ctx context.Context, userId int64, req *types.CourseListRequest) ([]*Course, int64, error) {
	// 实现用户报名的课程查询
	// 这里简化实现，实际应该关联查询
	return nil, 0, nil
}

// 数据转换函数
func (c *Course) ToCourseInfo() *types.CourseInfo {
	return &types.CourseInfo{
		ID:            c.ID,
		Title:         c.Title,
		Description:   c.Description,
		Subject:       c.Subject,
		Grade:         c.Grade,
		Difficulty:    c.Difficulty,
		CoverImage:    c.CoverImage.String,
		Price:         c.Price,
		Duration:      c.Duration,
		TeacherId:     c.TeacherId,
		InstitutionId: c.InstitutionId,
		Status:        c.Status,
		StudentCount:  c.StudentCount,
		Rating:        c.Rating,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

func (c *Course) ToCourseDetail(chapters []*types.Chapter) *types.CourseDetail {
	return &types.CourseDetail{
		CourseInfo: *c.ToCourseInfo(),
		Chapters:    chapters,
		Objectives:  c.Objectives.String,
		Prerequisites: c.Prerequisites.String,
	}
}

func (c *Chapter) ToChapter() *types.Chapter {
	return &types.Chapter{
		ID:          c.ID,
		CourseId:    c.CourseId,
		Title:       c.Title,
		Description: c.Description,
		Order:       c.Order,
		Duration:    c.Duration,
	}
}

func (l *Lesson) ToLesson() *types.Lesson {
	return &types.Lesson{
		ID:        l.ID,
		ChapterId: l.ChapterId,
		Title:     l.Title,
		Type:      l.Type,
		Content:   l.Content,
		Duration:  l.Duration,
		Order:     l.Order,
		IsFree:    l.IsFree,
	}
}