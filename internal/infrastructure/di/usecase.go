package di

import (
	opinion_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/opinion"
	talksession_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/talksession"
	user_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/user"

	analysis_usecase "github.com/neko-dream/server/internal/usecase/analysis"
	"github.com/neko-dream/server/internal/usecase/command/auth_command"
	"github.com/neko-dream/server/internal/usecase/command/opinion_command"
	"github.com/neko-dream/server/internal/usecase/command/talksession_command"
	"github.com/neko-dream/server/internal/usecase/command/timeline_command"
	"github.com/neko-dream/server/internal/usecase/command/user_command"
	"github.com/neko-dream/server/internal/usecase/command/vote_command"
	"github.com/neko-dream/server/internal/usecase/query/timeline_query"
)

// useCaseDeps returns a slice of dependency injection arguments for various use case handlers and query handlers across different domains of the application.
// The function constructs a comprehensive list of constructors for commands and queries related to analysis, talk sessions, opinions, users, voting, authentication, and timelines.
// Each constructor is paired with a nil argument, indicating no additional configuration is required during initialization.
// The returned slice can be used to register and configure use case dependencies in the application's dependency injection container.
func useCaseDeps() []ProvideArg {
	return []ProvideArg{
		{analysis_usecase.NewGetAnalysisResultUseCase, nil},
		{analysis_usecase.NewGetReportQueryHandler, nil},
		{talksession_command.NewAddConclusionCommandHandler, nil},
		{talksession_command.NewStartTalkSessionCommand, nil},
		{opinion_command.NewSubmitOpinionHandler, nil},
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
		{user_command.NewEditHandler, nil},
		{user_command.NewRegisterHandler, nil},
		{user_query.NewDetailHandler, nil},
		{vote_command.NewVoteHandler, nil},
		{auth_command.NewAuthLogin, nil},
		{auth_command.NewRevoke, nil},
		{auth_command.NewAuthCallback, nil},
		{timeline_command.NewAddTimeLine, nil},
		{timeline_command.NewEditTimeLine, nil},
		{timeline_query.NewGetTimeLine, nil},
	}
}
