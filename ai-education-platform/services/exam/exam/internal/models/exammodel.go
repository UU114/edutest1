package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ExamModel = (*customExamModel)(nil)

type (
	// ExamModel is an interface to be customized, add more methods here,
	// and implement the added methods in customExamModel.
	ExamModel interface {
		InsertQuestion(ctx context.Context, data *Question) (sql.Result, error)
		FindOneQuestion(ctx context.Context, id int64) (*Question, error)
		UpdateQuestion(ctx context.Context, data *Question) error
		DeleteQuestion(ctx context.Context, id int64) error
		GetQuestionList(ctx context.Context, req *QuestionListQuery) ([]*Question, int64, error)
		
		InsertExam(ctx context.Context, data *Exam) (sql.Result, error)
		FindOneExam(ctx context.Context, id int64) (*Exam, error)
		UpdateExam(ctx context.Context, data *Exam) error
		DeleteExam(ctx context.Context, id int64) error
		GetExamList(ctx context.Context, req *ExamListQuery) ([]*Exam, int64, error)
		
		InsertExamRecord(ctx context.Context, data *ExamRecord) (sql.Result, error)
		FindOneExamRecord(ctx context.Context, id int64) (*ExamRecord, error)
		UpdateExamRecord(ctx context.Context, data *ExamRecord) error
		GetUserExamRecords(ctx context.Context, userId int64) ([]*ExamRecord, error)
		
		InsertWrongQuestion(ctx context.Context, data *WrongQuestion) (sql.Result, error)
		FindOneWrongQuestion(ctx context.Context, id int64) (*WrongQuestion, error)
		UpdateWrongQuestion(ctx context.Context, data *WrongQuestion) error
		DeleteWrongQuestion(ctx context.Context, id int64) error
		GetUserWrongQuestions(ctx context.Context, req *WrongQuestionQuery) ([]*WrongQuestion, int64, error)
		
		InsertPracticeRecord(ctx context.Context, data *PracticeRecord) (sql.Result, error)
		GetUserPracticeRecords(ctx context.Context, userId int64) ([]*PracticeRecord, error)
		
		SmartGenerateQuestions(ctx context.Context, req *SmartGenerateQuery) ([]*Question, error)
	}

	customExamModel struct {
		conn  sqlx.SqlConn
		table string
	}
	
	// 题目查询条件
	QuestionListQuery struct {
		Page       int
		PageSize   int
		Subject    string
		Grade      string
		Type       string
		Difficulty int
		Keyword    string
		Tags       string
		CreatorId  int64
		Status     int
	}
	
	// 试卷查询条件
	ExamListQuery struct {
		Page     int
		PageSize int
		Subject  string
		Grade    string
		CreatorId int64
		Status   int
	}
	
	// 错题查询条件
	WrongQuestionQuery struct {
		Page     int
		PageSize int
		Subject  string
		Grade    string
		UserId   int64
		Mastered bool
	}
	
	// 智能出题查询
	SmartGenerateQuery struct {
		Subject       string
		Grade         string
		QuestionTypes []string
		Difficulty    int
		KnowledgePoints []int64
		QuestionCount int
		Tags          []string
		ExcludeIds    []int64
	}
	
	// 题目
	Question struct {
		Id             int64     `db:"id"`
		Title          string    `db:"title"`
		Type           string    `db:"type"`
		Subject        string    `db:"subject"`
		Grade          string    `db:"grade"`
		Difficulty     int       `db:"difficulty"`
		Content        string    `db:"content"`
		Options        string    `db:"options"`        // JSON格式
		CorrectAnswer  string    `db:"correct_answer"`
		Analysis       string    `db:"analysis"`
		KnowledgePoints string   `db:"knowledge_points"` // JSON格式
		Tags           string    `db:"tags"`           // JSON格式
		CreatorId      int64     `db:"creator_id"`
		Status         int       `db:"status"`
		CreatedAt      int64     `db:"created_at"`
		UpdatedAt      int64     `db:"updated_at"`
		UsageCount     int       `db:"usage_count"`
		CorrectRate    float64   `db:"correct_rate"`
	}
	
	// 试卷
	Exam struct {
		Id            int64     `db:"id"`
		Title         string    `db:"title"`
		Description   string    `db:"description"`
		Subject       string    `db:"subject"`
		Grade         string    `db:"grade"`
		Duration      int       `db:"duration"`
		TotalScore    float64   `db:"total_score"`
		PassScore     float64   `db:"pass_score"`
		Questions     string    `db:"questions"`     // JSON格式
		CreatorId     int64     `db:"creator_id"`
		Status        int       `db:"status"`
		CreatedAt     int64     `db:"created_at"`
		UpdatedAt     int64     `db:"updated_at"`
	}
	
	// 考试记录
	ExamRecord struct {
		Id         int64     `db:"id"`
		ExamId     int64     `db:"exam_id"`
		UserId     int64     `db:"user_id"`
		Score      float64   `db:"score"`
		TotalScore float64   `db:"total_score"`
		Status     string    `db:"status"`
		StartTime  int64     `db:"start_time"`
		EndTime    int64     `db:"end_time"`
		TimeUsed   int       `db:"time_used"`
		Answers    string    `db:"answers"`    // JSON格式
		CreatedAt  int64     `db:"created_at"`
		UpdatedAt  int64     `db:"updated_at"`
	}
	
	// 错题记录
	WrongQuestion struct {
		Id            int64  `db:"id"`
		UserId        int64  `db:"user_id"`
		QuestionId    int64  `db:"question_id"`
		QuestionTitle string `db:"question_title"`
		StudentAnswer string `db:"student_answer"`
		CorrectAnswer string `db:"correct_answer"`
		Subject       string `db:"subject"`
		Grade         string `db:"grade"`
		WrongCount    int    `db:"wrong_count"`
		LastWrongTime int64  `db:"last_wrong_time"`
		Mastered      bool   `db:"mastered"`
		CreatedAt     int64  `db:"created_at"`
		UpdatedAt     int64  `db:"updated_at"`
	}
	
	// 练习记录
	PracticeRecord struct {
		Id            int64   `db:"id"`
		UserId        int64   `db:"user_id"`
		Subject       string  `db:"subject"`
		QuestionCount int     `db:"question_count"`
		CorrectCount  int     `db:"correct_count"`
		Score         float64 `db:"score"`
		TimeUsed      int     `db:"time_used"`
		Difficulty    string  `db:"difficulty"`
		CreatedAt     int64   `db:"created_at"`
	}
)

// NewExamModel returns a model for the database table.
func NewExamModel(conn sqlx.SqlConn) ExamModel {
	return &customExamModel{
		conn:  conn,
		table: "questions", // 默认表名
	}
}

// Question methods
func (m *customExamModel) InsertQuestion(ctx context.Context, data *Question) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (title, type, subject, grade, difficulty, content, options, correct_answer, analysis, knowledge_points, tags, creator_id, status, created_at, updated_at, usage_count, correct_rate) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.Title, data.Type, data.Subject, data.Grade, data.Difficulty, data.Content, data.Options, data.CorrectAnswer, data.Analysis, data.KnowledgePoints, data.Tags, data.CreatorId, data.Status, data.CreatedAt, data.UpdatedAt, data.UsageCount, data.CorrectRate)
	return ret, err
}

func (m *customExamModel) FindOneQuestion(ctx context.Context, id int64) (*Question, error) {
	query := fmt.Sprintf("select id, title, type, subject, grade, difficulty, content, options, correct_answer, analysis, knowledge_points, tags, creator_id, status, created_at, updated_at, usage_count, correct_rate from %s where id = ?", m.table)
	var resp Question
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customExamModel) UpdateQuestion(ctx context.Context, data *Question) error {
	query := fmt.Sprintf("update %s set title = ?, type = ?, subject = ?, grade = ?, difficulty = ?, content = ?, options = ?, correct_answer = ?, analysis = ?, knowledge_points = ?, tags = ?, creator_id = ?, status = ?, updated_at = ?, usage_count = ?, correct_rate = ? where id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.Title, data.Type, data.Subject, data.Grade, data.Difficulty, data.Content, data.Options, data.CorrectAnswer, data.Analysis, data.KnowledgePoints, data.Tags, data.CreatorId, data.Status, data.UpdatedAt, data.UsageCount, data.CorrectRate, data.Id)
	return err
}

func (m *customExamModel) DeleteQuestion(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *customExamModel) GetQuestionList(ctx context.Context, req *QuestionListQuery) ([]*Question, int64, error) {
	whereClause, args := m.buildQuestionWhereClause(req)
	
	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	query := fmt.Sprintf("select id, title, type, subject, grade, difficulty, content, options, correct_answer, analysis, knowledge_points, tags, creator_id, status, created_at, updated_at, usage_count, correct_rate from %s %s order by created_at desc limit ? offset ?", m.table, whereClause)
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)
	
	var resp []*Question
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	return resp, total, err
}

// Exam methods
func (m *customExamModel) InsertExam(ctx context.Context, data *Exam) (sql.Result, error) {
	query := "insert into exams (title, description, subject, grade, duration, total_score, pass_score, questions, creator_id, status, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.Title, data.Description, data.Subject, data.Grade, data.Duration, data.TotalScore, data.PassScore, data.Questions, data.CreatorId, data.Status, data.CreatedAt, data.UpdatedAt)
	return ret, err
}

func (m *customExamModel) FindOneExam(ctx context.Context, id int64) (*Exam, error) {
	query := "select id, title, description, subject, grade, duration, total_score, pass_score, questions, creator_id, status, created_at, updated_at from exams where id = ?"
	var resp Exam
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customExamModel) UpdateExam(ctx context.Context, data *Exam) error {
	query := "update exams set title = ?, description = ?, subject = ?, grade = ?, duration = ?, total_score = ?, pass_score = ?, questions = ?, creator_id = ?, status = ?, updated_at = ? where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.Title, data.Description, data.Subject, data.Grade, data.Duration, data.TotalScore, data.PassScore, data.Questions, data.CreatorId, data.Status, data.UpdatedAt, data.Id)
	return err
}

func (m *customExamModel) DeleteExam(ctx context.Context, id int64) error {
	query := "delete from exams where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *customExamModel) GetExamList(ctx context.Context, req *ExamListQuery) ([]*Exam, int64, error) {
	whereClause, args := m.buildExamWhereClause(req)
	
	countQuery := fmt.Sprintf("select count(*) from exams %s", whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	query := fmt.Sprintf("select id, title, description, subject, grade, duration, total_score, pass_score, questions, creator_id, status, created_at, updated_at from exams %s order by created_at desc limit ? offset ?", whereClause)
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)
	
	var resp []*Exam
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	return resp, total, err
}

// ExamRecord methods
func (m *customExamModel) InsertExamRecord(ctx context.Context, data *ExamRecord) (sql.Result, error) {
	query := "insert into exam_records (exam_id, user_id, score, total_score, status, start_time, end_time, time_used, answers, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.ExamId, data.UserId, data.Score, data.TotalScore, data.Status, data.StartTime, data.EndTime, data.TimeUsed, data.Answers, data.CreatedAt, data.UpdatedAt)
	return ret, err
}

func (m *customExamModel) FindOneExamRecord(ctx context.Context, id int64) (*ExamRecord, error) {
	query := "select id, exam_id, user_id, score, total_score, status, start_time, end_time, time_used, answers, created_at, updated_at from exam_records where id = ?"
	var resp ExamRecord
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customExamModel) UpdateExamRecord(ctx context.Context, data *ExamRecord) error {
	query := "update exam_records set exam_id = ?, user_id = ?, score = ?, total_score = ?, status = ?, start_time = ?, end_time = ?, time_used = ?, answers = ?, updated_at = ? where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.ExamId, data.UserId, data.Score, data.TotalScore, data.Status, data.StartTime, data.EndTime, data.TimeUsed, data.Answers, data.UpdatedAt, data.Id)
	return err
}

func (m *customExamModel) GetUserExamRecords(ctx context.Context, userId int64) ([]*ExamRecord, error) {
	query := "select id, exam_id, user_id, score, total_score, status, start_time, end_time, time_used, answers, created_at, updated_at from exam_records where user_id = ? order by created_at desc"
	var resp []*ExamRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	return resp, err
}

// WrongQuestion methods
func (m *customExamModel) InsertWrongQuestion(ctx context.Context, data *WrongQuestion) (sql.Result, error) {
	query := "insert into wrong_questions (user_id, question_id, question_title, student_answer, correct_answer, subject, grade, wrong_count, last_wrong_time, mastered, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.QuestionId, data.QuestionTitle, data.StudentAnswer, data.CorrectAnswer, data.Subject, data.Grade, data.WrongCount, data.LastWrongTime, data.Mastered, data.CreatedAt, data.UpdatedAt)
	return ret, err
}

func (m *customExamModel) FindOneWrongQuestion(ctx context.Context, id int64) (*WrongQuestion, error) {
	query := "select id, user_id, question_id, question_title, student_answer, correct_answer, subject, grade, wrong_count, last_wrong_time, mastered, created_at, updated_at from wrong_questions where id = ?"
	var resp WrongQuestion
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *customExamModel) UpdateWrongQuestion(ctx context.Context, data *WrongQuestion) error {
	query := "update wrong_questions set user_id = ?, question_id = ?, question_title = ?, student_answer = ?, correct_answer = ?, subject = ?, grade = ?, wrong_count = ?, last_wrong_time = ?, mastered = ?, updated_at = ? where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.QuestionId, data.QuestionTitle, data.StudentAnswer, data.CorrectAnswer, data.Subject, data.Grade, data.WrongCount, data.LastWrongTime, data.Mastered, data.UpdatedAt, data.Id)
	return err
}

func (m *customExamModel) DeleteWrongQuestion(ctx context.Context, id int64) error {
	query := "delete from wrong_questions where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *customExamModel) GetUserWrongQuestions(ctx context.Context, req *WrongQuestionQuery) ([]*WrongQuestion, int64, error) {
	whereClause, args := m.buildWrongQuestionWhereClause(req)
	
	countQuery := fmt.Sprintf("select count(*) from wrong_questions %s", whereClause)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	
	query := fmt.Sprintf("select id, user_id, question_id, question_title, student_answer, correct_answer, subject, grade, wrong_count, last_wrong_time, mastered, created_at, updated_at from wrong_questions %s order by last_wrong_time desc limit ? offset ?", whereClause)
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)
	
	var resp []*WrongQuestion
	err = m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	return resp, total, err
}

// PracticeRecord methods
func (m *customExamModel) InsertPracticeRecord(ctx context.Context, data *PracticeRecord) (sql.Result, error) {
	query := "insert into practice_records (user_id, subject, question_count, correct_count, score, time_used, difficulty, created_at) values (?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.Subject, data.QuestionCount, data.CorrectCount, data.Score, data.TimeUsed, data.Difficulty, data.CreatedAt)
	return ret, err
}

func (m *customExamModel) GetUserPracticeRecords(ctx context.Context, userId int64) ([]*PracticeRecord, error) {
	query := "select id, user_id, subject, question_count, correct_count, score, time_used, difficulty, created_at from practice_records where user_id = ? order by created_at desc"
	var resp []*PracticeRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	return resp, err
}

// Smart generate questions
func (m *customExamModel) SmartGenerateQuestions(ctx context.Context, req *SmartGenerateQuery) ([]*Question, error) {
	whereClause, args := m.buildSmartGenerateWhereClause(req)
	
	query := fmt.Sprintf("select id, title, type, subject, grade, difficulty, content, options, correct_answer, analysis, knowledge_points, tags, creator_id, status, created_at, updated_at, usage_count, correct_rate from questions %s order by rand() limit ?", whereClause)
	args = append(args, req.QuestionCount)
	
	var resp []*Question
	err := m.conn.QueryRowsCtx(ctx, &resp, query, args...)
	return resp, err
}

// Helper methods
func (m *customExamModel) buildQuestionWhereClause(req *QuestionListQuery) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	
	if req.Subject != "" {
		conditions = append(conditions, "subject = ?")
		args = append(args, req.Subject)
	}
	
	if req.Grade != "" {
		conditions = append(conditions, "grade = ?")
		args = append(args, req.Grade)
	}
	
	if req.Type != "" {
		conditions = append(conditions, "type = ?")
		args = append(args, req.Type)
	}
	
	if req.Difficulty > 0 {
		conditions = append(conditions, "difficulty = ?")
		args = append(args, req.Difficulty)
	}
	
	if req.Keyword != "" {
		conditions = append(conditions, "(title like ? or content like ?)")
		keyword := "%" + req.Keyword + "%"
		args = append(args, keyword, keyword)
	}
	
	if req.Tags != "" {
		conditions = append(conditions, "tags like ?")
		args = append(args, "%"+req.Tags+"%")
	}
	
	if req.CreatorId > 0 {
		conditions = append(conditions, "creator_id = ?")
		args = append(args, req.CreatorId)
	}
	
	if req.Status >= 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, req.Status)
	}
	
	if len(conditions) == 0 {
		return "", args
	}
	
	return "where " + strings.Join(conditions, " and "), args
}

func (m *customExamModel) buildExamWhereClause(req *ExamListQuery) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	
	if req.Subject != "" {
		conditions = append(conditions, "subject = ?")
		args = append(args, req.Subject)
	}
	
	if req.Grade != "" {
		conditions = append(conditions, "grade = ?")
		args = append(args, req.Grade)
	}
	
	if req.CreatorId > 0 {
		conditions = append(conditions, "creator_id = ?")
		args = append(args, req.CreatorId)
	}
	
	if req.Status >= 0 {
		conditions = append(conditions, "status = ?")
		args = append(args, req.Status)
	}
	
	if len(conditions) == 0 {
		return "", args
	}
	
	return "where " + strings.Join(conditions, " and "), args
}

func (m *customExamModel) buildWrongQuestionWhereClause(req *WrongQuestionQuery) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	
	conditions = append(conditions, "user_id = ?")
	args = append(args, req.UserId)
	
	if req.Subject != "" {
		conditions = append(conditions, "subject = ?")
		args = append(args, req.Subject)
	}
	
	if req.Grade != "" {
		conditions = append(conditions, "grade = ?")
		args = append(args, req.Grade)
	}
	
	if req.Mastered {
		conditions = append(conditions, "mastered = ?")
		args = append(args, req.Mastered)
	}
	
	return "where " + strings.Join(conditions, " and "), args
}

func (m *customExamModel) buildSmartGenerateWhereClause(req *SmartGenerateQuery) (string, []interface{}) {
	var conditions []string
	var args []interface{}
	
	conditions = append(conditions, "subject = ?")
	args = append(args, req.Subject)
	
	conditions = append(conditions, "grade = ?")
	args = append(args, req.Grade)
	
	conditions = append(conditions, "status = ?")
	args = append(args, 1) // 只选择已发布的题目
	
	if req.Difficulty > 0 {
		conditions = append(conditions, "difficulty = ?")
		args = append(args, req.Difficulty)
	}
	
	if len(req.QuestionTypes) > 0 {
		placeholders := strings.Repeat("?,", len(req.QuestionTypes))
		placeholders = placeholders[:len(placeholders)-1]
		conditions = append(conditions, fmt.Sprintf("type in (%s)", placeholders))
		for _, qt := range req.QuestionTypes {
			args = append(args, qt)
		}
	}
	
	if len(req.ExcludeIds) > 0 {
		placeholders := strings.Repeat("?,", len(req.ExcludeIds))
		placeholders = placeholders[:len(placeholders)-1]
		conditions = append(conditions, fmt.Sprintf("id not in (%s)", placeholders))
		for _, id := range req.ExcludeIds {
			args = append(args, id)
		}
	}
	
	return "where " + strings.Join(conditions, " and "), args
}

// Helper functions for JSON serialization
func StringSliceToJSON(slice []string) (string, error) {
	if slice == nil {
		return "[]", nil
	}
	data, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func JSONToStringSlice(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return []string{}, nil
	}
	var slice []string
	err := json.Unmarshal([]byte(jsonStr), &slice)
	return slice, err
}

func Int64SliceToJSON(slice []int64) (string, error) {
	if slice == nil {
		return "[]", nil
	}
	data, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func JSONToInt64Slice(jsonStr string) ([]int64, error) {
	if jsonStr == "" {
		return []int64{}, nil
	}
	var slice []int64
	err := json.Unmarshal([]byte(jsonStr), &slice)
	return slice, err
}

var ErrNotFound = sqlx.ErrNotFound