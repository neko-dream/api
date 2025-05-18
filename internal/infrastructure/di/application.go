package di

import (
	analysis_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/analysis"
	opinion_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/opinion"
	report_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/report"
	talksession_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/talksession"
	user_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/user"

	"github.com/neko-dream/server/internal/application/command/auth_command"
	"github.com/neko-dream/server/internal/application/command/image_command"
	"github.com/neko-dream/server/internal/application/command/opinion_command"
	"github.com/neko-dream/server/internal/application/command/organization_command"
	"github.com/neko-dream/server/internal/application/command/policy_command"
	"github.com/neko-dream/server/internal/application/command/report_command"
	"github.com/neko-dream/server/internal/application/command/talksession_command"
	"github.com/neko-dream/server/internal/application/command/timeline_command"
	"github.com/neko-dream/server/internal/application/command/user_command"
	"github.com/neko-dream/server/internal/application/command/vote_command"
	opinion_q "github.com/neko-dream/server/internal/application/query/opinion"
	"github.com/neko-dream/server/internal/application/query/policy_query"
	report_q "github.com/neko-dream/server/internal/application/query/report_query"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/internal/application/query/timeline_query"
)

// このファイルはアプリケーション層（ユースケース/クエリ）のコンストラクタを管理します。
// 新しいユースケースやクエリを追加した場合は必ずここに追記してください。

func useCaseDeps() []ProvideArg {
	return []ProvideArg{
		{talksession_command.NewAddConclusionCommandHandler, nil},
		{talksession_command.NewStartTalkSessionCommand, nil},
		{talksession_command.NewTakeConsentUseCase, nil},
		{talksession_command.NewEditCommand, nil},
		{talksession_query.NewBrowseTalkSessionQueryHandler, nil},
		{talksession_query.NewBrowseOpenedByUserQueryHandler, nil},
		{talksession_query.NewBrowseJoinedTalkSessionQueryHandler, nil},
		{talksession_query.NewGetTalkSessionDetailByIDQueryHandler, nil},
		{talksession_query.NewGetConclusionByIDQueryHandler, nil},
		{talksession_query.NewGetRestrictionsQuery, nil},
		{talksession_query.NewHasConsentQuery, nil},
		{talksession.NewIsTalkSessionSatisfiedInteractor, nil},
		{opinion_command.NewSubmitOpinionHandler, nil},
		{opinion_command.NewReportOpinion, nil},
		{opinion_query.NewGetOpinionsByTalkSessionIDQueryHandler, nil},
		{opinion_query.NewGetOpinionDetailByIDQueryHandler, nil},
		{opinion_query.NewGetOpinionRepliesQueryHandler, nil},
		{opinion_query.NewSwipeOpinionsQueryHandler, nil},
		{opinion_query.NewGetMyOpinionsQueryHandler, nil},
		{opinion_query.NewGetOpinionGroupRatioInteractor, nil},
		{opinion_q.NewGetReportReasons, nil},
		{user_command.NewEditHandler, nil},
		{user_command.NewRegisterHandler, nil},
		{user_query.NewDetailHandler, nil},
		{vote_command.NewVoteHandler, nil},
		{auth_command.NewAuthLogin, nil},
		{auth_command.NewRevoke, nil},
		{auth_command.NewAuthCallback, nil},
		{auth_command.NewLoginForDev, nil},
		{auth_command.NewDetachAccount, nil},
		{auth_command.NewPasswordRegister, nil},
		{auth_command.NewPasswordLogin, nil},
		{auth_command.NewChangePassword, nil},
		{timeline_command.NewAddTimeLine, nil},
		{timeline_command.NewEditTimeLine, nil},
		{timeline_query.NewGetTimeLine, nil},
		{analysis_query.NewGetAnalysisResultHandler, nil},
		{analysis_query.NewGetReportQueryHandler, nil},
		{report_query.NewGetByTalkSessionQueryInteractor, nil},
		{report_query.NewGetOpinionReportQueryInteractor, nil},
		{report_command.NewSolveReportCommandInteractor, nil},
		{report_q.NewGetCountQueryInteractor, nil},
		{image_command.NewUploadImageHandler, nil},
		{policy_command.NewAcceptPolicy, nil},
		{policy_query.NewCheckConsent, nil},
		{organization_command.NewCreateOrganizationInteractor, nil},
		{organization_command.NewInviteOrganizationInteractor, nil},
		{organization_command.NewInviteOrganizationForUserInteractor, nil},
	}
}
