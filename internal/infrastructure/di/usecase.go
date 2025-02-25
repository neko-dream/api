package di

import (
	analysis_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/analysis"
	opinion_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/opinion"
	talksession_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/talksession"
	user_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/user"

	"github.com/neko-dream/server/internal/usecase/command/auth_command"
	"github.com/neko-dream/server/internal/usecase/command/image_command"
	"github.com/neko-dream/server/internal/usecase/command/opinion_command"
	"github.com/neko-dream/server/internal/usecase/command/talksession_command"
	"github.com/neko-dream/server/internal/usecase/command/timeline_command"
	"github.com/neko-dream/server/internal/usecase/command/user_command"
	"github.com/neko-dream/server/internal/usecase/command/vote_command"
	"github.com/neko-dream/server/internal/usecase/query/timeline_query"
)

func useCaseDeps() []ProvideArg {
	return []ProvideArg{
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
		{auth_command.NewLoginForDev, nil},
		{timeline_command.NewAddTimeLine, nil},
		{timeline_command.NewEditTimeLine, nil},
		{timeline_query.NewGetTimeLine, nil},
		{analysis_query.NewGetAnalysisResultHandler, nil},
		{analysis_query.NewGetReportQueryHandler, nil},
		{image_command.NewUploadImageHandler, nil},
		{talksession_query.NewGetRestrictionsQuery, nil},
	}
}
