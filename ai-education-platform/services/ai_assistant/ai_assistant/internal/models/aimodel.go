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

var _ AIModel = (*customAIModel)(nil)

type (
	// AIModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAIModel.
	AIModel interface {
		InsertAIExplanation(ctx context.Context, data *AIExplanation) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*AIExplanation, error)
		UpdateAIExplanation(ctx context.Context, data *AIExplanation) error
		DeleteAIExplanation(ctx context.Context, id int64) error
		
		InsertAIChat(ctx context.Context, data *AIChat) (sql.Result, error)
		FindAIChatById(ctx context.Context, id int64) (*AIChat, error)
		UpdateAIChat(ctx context.Context, data *AIChat) error
		DeleteAIChat(ctx context.Context, id int64) error
		GetUserChatHistory(ctx context.Context, userId int64, limit int) ([]*AIChat, error)
		
		InsertLearningPath(ctx context.Context, data *LearningPath) (sql.Result, error)
		FindLearningPathById(ctx context.Context, id int64) (*LearningPath, error)
		UpdateLearningPath(ctx context.Context, data *LearningPath) error
		DeleteLearningPath(ctx context.Context, id int64) error
		GetUserLearningPaths(ctx context.Context, userId int64) ([]*LearningPath, error)
		
		InsertHomeworkGrade(ctx context.Context, data *HomeworkGrade) (sql.Result, error)
		FindHomeworkGradeById(ctx context.Context, id int64) (*HomeworkGrade, error)
		UpdateHomeworkGrade(ctx context.Context, data *HomeworkGrade) error
		DeleteHomeworkGrade(ctx context.Context, id int64) error
		
		InsertSpeechRecognition(ctx context.Context, data *SpeechRecognition) (sql.Result, error)
		FindSpeechRecognitionById(ctx context.Context, id int64) (*SpeechRecognition, error)
		UpdateSpeechRecognition(ctx context.Context, data *SpeechRecognition) error
		DeleteSpeechRecognition(ctx context.Context, id int64) error
	}

	customAIModel struct {
		conn  sqlx.SqlConn
		table string
	}
	
	// AI讲解记录
	AIExplanation struct {
		Id            int64     `db:"id"`
		ExplanationId string    `db:"explanation_id"`
		UserId        int64     `db:"user_id"`
		Subject       string    `db:"subject"`
		Grade         string    `db:"grade"`
		Topic         string    `db:"topic"`
		Difficulty    int       `db:"difficulty"`
		Language      string    `db:"language"`
		Style         string    `db:"style"`
		Content       string    `db:"content"`
		Summary       string    `db:"summary"`
		KeyPoints     string    `db:"key_points"` // JSON格式
		Examples      string    `db:"examples"`   // JSON格式
		Resources     string    `db:"resources"`  // JSON格式
		EstimatedTime int       `db:"estimated_time"`
		CreatedAt     int64     `db:"created_at"`
		UpdatedAt     int64     `db:"updated_at"`
	}
	
	// AI聊天记录
	AIChat struct {
		Id           int64     `db:"id"`
		AnswerId     string    `db:"answer_id"`
		UserId       int64     `db:"user_id"`
		SessionId    string    `db:"session_id"`
		Question     string    `db:"question"`
		Answer       string    `db:"answer"`
		Subject      string    `db:"subject"`
		Grade        string    `db:"grade"`
		Context      string    `db:"context"`
		Confidence   float64   `db:"confidence"`
		Suggestions  string    `db:"suggestions"`  // JSON格式
		RelatedTopics string   `db:"related_topics"` // JSON格式
		History      string    `db:"history"`     // JSON格式
		CreatedAt    int64     `db:"created_at"`
	}
	
	// 学习路径记录
	LearningPath struct {
		Id             int64     `db:"id"`
		PathId         string    `db:"path_id"`
		UserId         int64     `db:"user_id"`
		Subject        string    `db:"subject"`
		Grade          string    `db:"grade"`
		Goal           string    `db:"goal"`
		CurrentLevel   int       `db:"current_level"`
		WeakPoints     string    `db:"weak_points"`    // JSON格式
		StrongPoints   string    `db:"strong_points"`  // JSON格式
		StudyTime      int       `db:"study_time"`
		Title          string    `db:"title"`
		Description    string    `db:"description"`
		Duration       int       `db:"duration"`
		Steps          string    `db:"steps"`         // JSON格式
		Milestones     string    `db:"milestones"`    // JSON格式
		Recommendations string   `db:"recommendations"` // JSON格式
		Status         string    `db:"status"`        // active, completed, paused
		Progress       float64   `db:"progress"`      // 0-100
		CreatedAt      int64     `db:"created_at"`
		UpdatedAt      int64     `db:"updated_at"`
	}
	
	// 作业批改记录
	HomeworkGrade struct {
		Id              int64     `db:"id"`
		GradeId         string    `db:"grade_id"`
		UserId          int64     `db:"user_id"`
		Question        string    `db:"question"`
		StudentAnswer   string    `db:"student_answer"`
		Subject         string    `db:"subject"`
		Grade           string    `db:"grade"`
		QuestionType    string    `db:"question_type"`
		ExpectedAnswer  string    `db:"expected_answer"`
		GradingCriteria string    `db:"grading_criteria"`
		Score           float64   `db:"score"`
		MaxScore        float64   `db:"max_score"`
		Feedback        string    `db:"feedback"`
		Strengths       string    `db:"strengths"`   // JSON格式
		Weaknesses      string    `db:"weaknesses"`  // JSON格式
		Suggestions     string    `db:"suggestions"` // JSON格式
		CorrectAnswer   string    `db:"correct_answer"`
		Explanation     string    `db:"explanation"`
		Confidence      float64   `db:"confidence"`
		CreatedAt       int64     `db:"created_at"`
	}
	
	// 语音识别记录
	SpeechRecognition struct {
		Id                 int64     `db:"id"`
		RecognitionId      string    `db:"recognition_id"`
		UserId             int64     `db:"user_id"`
		AudioData          string    `db:"audio_data"`
		Language           string    `db:"language"`
		Subject            string    `db:"subject"`
		ExpectedText       string    `db:"expected_text"`
		Text               string    `db:"text"`
		Confidence         float64   `db:"confidence"`
		PronunciationScore float64   `db:"pronunciation_score"`
		FluencyScore       float64   `db:"fluency_score"`
		Feedback           string    `db:"feedback"`
		CreatedAt          int64     `db:"created_at"`
	}
)

// NewAIModel returns a model for the database table.
func NewAIModel(conn sqlx.SqlConn) AIModel {
	return &customAIModel{
		conn:  conn,
		table: "ai_explanations", // 默认表名
	}
}

// AI Explanation methods
func (m *customAIModel) InsertAIExplanation(ctx context.Context, data *AIExplanation) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (explanation_id, user_id, subject, grade, topic, difficulty, language, style, content, summary, key_points, examples, resources, estimated_time, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	ret, err := m.conn.ExecCtx(ctx, query, data.ExplanationId, data.UserId, data.Subject, data.Grade, data.Topic, data.Difficulty, data.Language, data.Style, data.Content, data.Summary, data.KeyPoints, data.Examples, data.Resources, data.EstimatedTime, data.CreatedAt, data.UpdatedAt)
	return ret, err
}

func (m *customAIModel) FindOne(ctx context.Context, id int64) (*AIExplanation, error) {
	query := fmt.Sprintf("select id, explanation_id, user_id, subject, grade, topic, difficulty, language, style, content, summary, key_points, examples, resources, estimated_time, created_at, updated_at from %s where id = ?", m.table)
	var resp AIExplanation
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

func (m *customAIModel) UpdateAIExplanation(ctx context.Context, data *AIExplanation) error {
	query := fmt.Sprintf("update %s set explanation_id = ?, user_id = ?, subject = ?, grade = ?, topic = ?, difficulty = ?, language = ?, style = ?, content = ?, summary = ?, key_points = ?, examples = ?, resources = ?, estimated_time = ?, updated_at = ? where id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, data.ExplanationId, data.UserId, data.Subject, data.Grade, data.Topic, data.Difficulty, data.Language, data.Style, data.Content, data.Summary, data.KeyPoints, data.Examples, data.Resources, data.EstimatedTime, data.UpdatedAt, data.Id)
	return err
}

func (m *customAIModel) DeleteAIExplanation(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// AI Chat methods
func (m *customAIModel) InsertAIChat(ctx context.Context, data *AIChat) (sql.Result, error) {
	query := "insert into ai_chats (answer_id, user_id, session_id, question, answer, subject, grade, context, confidence, suggestions, related_topics, history, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.AnswerId, data.UserId, data.SessionId, data.Question, data.Answer, data.Subject, data.Grade, data.Context, data.Confidence, data.Suggestions, data.RelatedTopics, data.History, data.CreatedAt)
	return ret, err
}

func (m *customAIModel) FindAIChatById(ctx context.Context, id int64) (*AIChat, error) {
	query := "select id, answer_id, user_id, session_id, question, answer, subject, grade, context, confidence, suggestions, related_topics, history, created_at from ai_chats where id = ?"
	var resp AIChat
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

func (m *customAIModel) UpdateAIChat(ctx context.Context, data *AIChat) error {
	query := "update ai_chats set answer_id = ?, user_id = ?, session_id = ?, question = ?, answer = ?, subject = ?, grade = ?, context = ?, confidence = ?, suggestions = ?, related_topics = ?, history = ? where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.AnswerId, data.UserId, data.SessionId, data.Question, data.Answer, data.Subject, data.Grade, data.Context, data.Confidence, data.Suggestions, data.RelatedTopics, data.History, data.Id)
	return err
}

func (m *customAIModel) DeleteAIChat(ctx context.Context, id int64) error {
	query := "delete from ai_chats where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *customAIModel) GetUserChatHistory(ctx context.Context, userId int64, limit int) ([]*AIChat, error) {
	query := "select id, answer_id, user_id, session_id, question, answer, subject, grade, context, confidence, suggestions, related_topics, history, created_at from ai_chats where user_id = ? order by created_at desc limit ?"
	var resp []*AIChat
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId, limit)
	return resp, err
}

// Learning Path methods
func (m *customAIModel) InsertLearningPath(ctx context.Context, data *LearningPath) (sql.Result, error) {
	query := "insert into learning_paths (path_id, user_id, subject, grade, goal, current_level, weak_points, strong_points, study_time, title, description, duration, steps, milestones, recommendations, status, progress, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.PathId, data.UserId, data.Subject, data.Grade, data.Goal, data.CurrentLevel, data.WeakPoints, data.StrongPoints, data.StudyTime, data.Title, data.Description, data.Duration, data.Steps, data.Milestones, data.Recommendations, data.Status, data.Progress, data.CreatedAt, data.UpdatedAt)
	return ret, err
}

func (m *customAIModel) FindLearningPathById(ctx context.Context, id int64) (*LearningPath, error) {
	query := "select id, path_id, user_id, subject, grade, goal, current_level, weak_points, strong_points, study_time, title, description, duration, steps, milestones, recommendations, status, progress, created_at, updated_at from learning_paths where id = ?"
	var resp LearningPath
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

func (m *customAIModel) UpdateLearningPath(ctx context.Context, data *LearningPath) error {
	query := "update learning_paths set path_id = ?, user_id = ?, subject = ?, grade = ?, goal = ?, current_level = ?, weak_points = ?, strong_points = ?, study_time = ?, title = ?, description = ?, duration = ?, steps = ?, milestones = ?, recommendations = ?, status = ?, progress = ?, updated_at = ? where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.PathId, data.UserId, data.Subject, data.Grade, data.Goal, data.CurrentLevel, data.WeakPoints, data.StrongPoints, data.StudyTime, data.Title, data.Description, data.Duration, data.Steps, data.Milestones, data.Recommendations, data.Status, data.Progress, data.UpdatedAt, data.Id)
	return err
}

func (m *customAIModel) DeleteLearningPath(ctx context.Context, id int64) error {
	query := "delete from learning_paths where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *customAIModel) GetUserLearningPaths(ctx context.Context, userId int64) ([]*LearningPath, error) {
	query := "select id, path_id, user_id, subject, grade, goal, current_level, weak_points, strong_points, study_time, title, description, duration, steps, milestones, recommendations, status, progress, created_at, updated_at from learning_paths where user_id = ? order by created_at desc"
	var resp []*LearningPath
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	return resp, err
}

// Homework Grade methods
func (m *customAIModel) InsertHomeworkGrade(ctx context.Context, data *HomeworkGrade) (sql.Result, error) {
	query := "insert into homework_grades (grade_id, user_id, question, student_answer, subject, grade, question_type, expected_answer, grading_criteria, score, max_score, feedback, strengths, weaknesses, suggestions, correct_answer, explanation, confidence, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.GradeId, data.UserId, data.Question, data.StudentAnswer, data.Subject, data.Grade, data.QuestionType, data.ExpectedAnswer, data.GradingCriteria, data.Score, data.MaxScore, data.Feedback, data.Strengths, data.Weaknesses, data.Suggestions, data.CorrectAnswer, data.Explanation, data.Confidence, data.CreatedAt)
	return ret, err
}

func (m *customAIModel) FindHomeworkGradeById(ctx context.Context, id int64) (*HomeworkGrade, error) {
	query := "select id, grade_id, user_id, question, student_answer, subject, grade, question_type, expected_answer, grading_criteria, score, max_score, feedback, strengths, weaknesses, suggestions, correct_answer, explanation, confidence, created_at from homework_grades where id = ?"
	var resp HomeworkGrade
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

func (m *customAIModel) UpdateHomeworkGrade(ctx context.Context, data *HomeworkGrade) error {
	query := "update homework_grades set grade_id = ?, user_id = ?, question = ?, student_answer = ?, subject = ?, grade = ?, question_type = ?, expected_answer = ?, grading_criteria = ?, score = ?, max_score = ?, feedback = ?, strengths = ?, weaknesses = ?, suggestions = ?, correct_answer = ?, explanation = ?, confidence = ? where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.GradeId, data.UserId, data.Question, data.StudentAnswer, data.Subject, data.Grade, data.QuestionType, data.ExpectedAnswer, data.GradingCriteria, data.Score, data.MaxScore, data.Feedback, data.Strengths, data.Weaknesses, data.Suggestions, data.CorrectAnswer, data.Explanation, data.Confidence, data.Id)
	return err
}

func (m *customAIModel) DeleteHomeworkGrade(ctx context.Context, id int64) error {
	query := "delete from homework_grades where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

// Speech Recognition methods
func (m *customAIModel) InsertSpeechRecognition(ctx context.Context, data *SpeechRecognition) (sql.Result, error) {
	query := "insert into speech_recognitions (recognition_id, user_id, audio_data, language, subject, expected_text, text, confidence, pronunciation_score, fluency_score, feedback, created_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ret, err := m.conn.ExecCtx(ctx, query, data.RecognitionId, data.UserId, data.AudioData, data.Language, data.Subject, data.ExpectedText, data.Text, data.Confidence, data.PronunciationScore, data.FluencyScore, data.Feedback, data.CreatedAt)
	return ret, err
}

func (m *customAIModel) FindSpeechRecognitionById(ctx context.Context, id int64) (*SpeechRecognition, error) {
	query := "select id, recognition_id, user_id, audio_data, language, subject, expected_text, text, confidence, pronunciation_score, fluency_score, feedback, created_at from speech_recognitions where id = ?"
	var resp SpeechRecognition
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

func (m *customAIModel) UpdateSpeechRecognition(ctx context.Context, data *SpeechRecognition) error {
	query := "update speech_recognitions set recognition_id = ?, user_id = ?, audio_data = ?, language = ?, subject = ?, expected_text = ?, text = ?, confidence = ?, pronunciation_score = ?, fluency_score = ?, feedback = ? where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, data.RecognitionId, data.UserId, data.AudioData, data.Language, data.Subject, data.ExpectedText, data.Text, data.Confidence, data.PronunciationScore, data.FluencyScore, data.Feedback, data.Id)
	return err
}

func (m *customAIModel) DeleteSpeechRecognition(ctx context.Context, id int64) error {
	query := "delete from speech_recognitions where id = ?"
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
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

func MapToJSON(m map[string]interface{}) (string, error) {
	if m == nil {
		return "{}", nil
	}
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func JSONToMap(jsonStr string) (map[string]interface{}, error) {
	if jsonStr == "" {
		return map[string]interface{}{}, nil
	}
	var m map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &m)
	return m, err
}

var ErrNotFound = sqlx.ErrNotFound