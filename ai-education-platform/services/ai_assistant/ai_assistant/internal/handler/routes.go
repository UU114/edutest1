package handler

import (
	"net/http"

	"ai-education-platform/services/ai_assistant/ai_assistant/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/ai/explain",
				Handler: AIExplainHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/ai/chat",
				Handler: AIChatHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/ai/learning-path",
				Handler: GetLearningPathHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/ai/recommendations",
				Handler: GetStudyRecommendationsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/ai/grade-homework",
				Handler: GradeHomeworkHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/ai/speech-recognition",
				Handler: SpeechRecognitionHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/v1/ai/chat/history",
				Handler: GetChatHistoryHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/v1/ai/chat/session/:id",
				Handler: DeleteChatSessionHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/teacher/ai/teaching-assistant",
				Handler: TeachingAssistantHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/teacher/ai/generate-exercises",
				Handler: GenerateExercisesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/teacher/ai/course-optimization",
				Handler: CourseOptimizationHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/teacher/ai/student-analysis",
				Handler: StudentAnalysisHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/parent/ai/parenting-advice",
				Handler: ParentingAdviceHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/parent/ai/child-report-analysis",
				Handler: ChildReportAnalysisHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/parent/ai/family-activity",
				Handler: FamilyActivitySuggestionHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}