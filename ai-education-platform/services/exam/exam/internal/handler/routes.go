package handler

import (
	"net/http"

	"ai-education-platform/services/exam/exam/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/questions",
				Handler: GetQuestionListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/question/:id",
				Handler: GetQuestionDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/exams",
				Handler: GetExamListHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/exam/:id",
				Handler: GetExamDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/exam/:id/start",
				Handler: StartExamHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/exam/answer",
				Handler: SubmitAnswerHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/exam/:id/finish",
				Handler: FinishExamHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/exam/:id/result",
				Handler: GetExamResultHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/my/exam-records",
				Handler: GetMyExamRecordsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/practice/start",
				Handler: StartPracticeHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/practice/answer",
				Handler: SubmitPracticeAnswerHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/practice/:id/finish",
				Handler: FinishPracticeHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/wrong-book",
				Handler: GetWrongBookHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/wrong-question/:id/master",
				Handler: MarkQuestionMasteredHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/practice/recommend",
				Handler: SmartRecommendPracticeHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/question",
				Handler: CreateQuestionHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/question/:id",
				Handler: UpdateQuestionHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/question/:id",
				Handler: DeleteQuestionHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/exam",
				Handler: CreateExamHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/exam/:id",
				Handler: UpdateExamHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/exam/:id",
				Handler: DeleteExamHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/questions/generate",
				Handler: SmartGenerateQuestionsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/questions/stats",
				Handler: GetQuestionStatsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/questions/import",
				Handler: ImportQuestionsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/questions/export",
				Handler: ExportQuestionsHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/teacher"),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/questions/stats",
				Handler: GetInstitutionQuestionStatsHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/exams/stats",
				Handler: GetInstitutionExamStatsHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/question/:id/review",
				Handler: ReviewQuestionHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/exam/:id/review",
				Handler: ReviewExamHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/institution"),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)
}