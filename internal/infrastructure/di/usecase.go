package di

import (
	queryimpl "github.com/neko-dream/server/internal/infrastructure/persistence/query"
	"github.com/neko-dream/server/internal/usecase/command"

	analysis_usecase "github.com/neko-dream/server/internal/usecase/analysis"
	auth_usecase "github.com/neko-dream/server/internal/usecase/auth"
	opinion_usecase "github.com/neko-dream/server/internal/usecase/opinion"
	talk_session_usecase "github.com/neko-dream/server/internal/usecase/talk_session"
	timeline_usecase "github.com/neko-dream/server/internal/usecase/timeline"
	user_usecase "github.com/neko-dream/server/internal/usecase/user"
	vote_usecase "github.com/neko-dream/server/internal/usecase/vote"
)

func useCaseDeps() []ProvideArg {
	return []ProvideArg{
		{auth_usecase.NewAuthLoginUseCase, nil},
		{auth_usecase.NewAuthCallbackUseCase, nil},
		{auth_usecase.NewRevokeUseCase, nil},
		{user_usecase.NewRegisterUserUseCase, nil},
		{user_usecase.NewEditUserUseCase, nil},
		{user_usecase.NewGetUserInformationQueryHandler, nil},
		{talk_session_usecase.NewStartTalkSessionCommand, nil},
		{talk_session_usecase.NewViewTalkSessionDetailQuery, nil},
		{talk_session_usecase.NewSearchTalkSessionsQuery, nil},
		{talk_session_usecase.NewBrowseUsersTalkSessionHistoriesQueryHandler, nil},
		{talk_session_usecase.NewGetTalkSessionConclusionQuery, nil},
		{opinion_usecase.NewPostOpinionUseCase, nil},
		{opinion_usecase.NewGetOpinionRepliesUseCase, nil},
		{opinion_usecase.NewGetSwipeOpinionsQueryHandler, nil},
		{opinion_usecase.NewGetOpinionDetailUseCase, nil},
		{opinion_usecase.NewGetUserOpinionListQueryHandler, nil},
		{opinion_usecase.NewGetOpinionsByTalkSessionUseCase, nil},
		{analysis_usecase.NewGetAnalysisResultUseCase, nil},
		{analysis_usecase.NewGetReportQueryHandler, nil},
		{timeline_usecase.NewAddTimeLineUseCase, nil},
		{timeline_usecase.NewGetTimeLineUseCase, nil},
		{timeline_usecase.NewEditTimeLineUseCase, nil},
		{vote_usecase.NewPostVoteUseCase, nil},
		{command.NewAddConclusionCommandHandler, nil},
		{queryimpl.NewBrowseTalkSessionQuery, nil},
	}
}
