package di

import (
	opinion_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/opinion"
	talksession_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/talksession"

	analysis_usecase "github.com/neko-dream/server/internal/usecase/analysis"
	auth_usecase "github.com/neko-dream/server/internal/usecase/auth"
	"github.com/neko-dream/server/internal/usecase/command"
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
		{analysis_usecase.NewGetAnalysisResultUseCase, nil},
		{analysis_usecase.NewGetReportQueryHandler, nil},
		{timeline_usecase.NewAddTimeLineUseCase, nil},
		{timeline_usecase.NewGetTimeLineUseCase, nil},
		{timeline_usecase.NewEditTimeLineUseCase, nil},
		{vote_usecase.NewPostVoteUseCase, nil},
		{command.NewAddConclusionCommandHandler, nil},
		{command.NewStartTalkSessionCommand, nil},
		{command.NewSubmitOpinionCommandHandler, nil},
		{talksession_query.NewBrowseTalkSessionQueryHandler, nil},
		{talksession_query.NewBrowseOpenedByUserQueryHandler, nil},
		{talksession_query.NewBrowseJoinedTalkSessionQueryHandler, nil},
		{talksession_query.NewGetTalkSessionDetailByIDQueryHandler, nil},
		{talksession_query.NewGetConclusionByIDQueryHandler, nil},
		{opinion_query.NewGetOpinionsByTalkSessionIDQueryHandler, nil},
		{opinion_query.NewGetOpinionDetailByIDQueryHandler, nil},
		{opinion_query.NewGetOpinionRepliesQueryHandler, nil},
		{opinion_query.NewSwipeOpinionsQueryHandler, nil},
		{opinion_query.NewGetMyOpinionsQueryHandler, nil},
	}
}
