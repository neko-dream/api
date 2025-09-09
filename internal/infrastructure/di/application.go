package di

import (
	analysis_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/analysis"
	opinion_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/opinion"
	report_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/report"
	talksession_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/talksession"
	user_query "github.com/neko-dream/server/internal/infrastructure/persistence/query/user"

	"github.com/neko-dream/server/internal/application/event_processor"
	"github.com/neko-dream/server/internal/application/event_processor/handlers"
	opinion_q "github.com/neko-dream/server/internal/application/query/opinion"
	"github.com/neko-dream/server/internal/application/query/organization_query"
	"github.com/neko-dream/server/internal/application/query/policy_query"
	report_q "github.com/neko-dream/server/internal/application/query/report_query"
	"github.com/neko-dream/server/internal/application/query/talksession"
	"github.com/neko-dream/server/internal/application/query/timeline_query"
	"github.com/neko-dream/server/internal/application/usecase/analysis_usecase"
	"github.com/neko-dream/server/internal/application/usecase/auth_usecase"
	"github.com/neko-dream/server/internal/application/usecase/image_usecase"
	"github.com/neko-dream/server/internal/application/usecase/opinion_usecase"
	"github.com/neko-dream/server/internal/application/usecase/organization_usecase"
	"github.com/neko-dream/server/internal/application/usecase/policy_usecase"
	"github.com/neko-dream/server/internal/application/usecase/report_usecase"
	"github.com/neko-dream/server/internal/application/usecase/talksession_usecase"
	"github.com/neko-dream/server/internal/application/usecase/timeline_usecase"
	"github.com/neko-dream/server/internal/application/usecase/user_usecase"
	"github.com/neko-dream/server/internal/application/usecase/vote_usecase"
)

// このファイルはアプリケーション層（ユースケース/クエリ）のコンストラクタを管理します。
// 新しいユースケースやクエリを追加した場合は必ずここに追記してください。

func useCaseDeps() []ProvideArg {
	return []ProvideArg{
		{talksession_usecase.NewAddConclusionCommandHandler, nil},
		{talksession_usecase.NewStartTalkSessionUseCase, nil},
		{talksession_usecase.NewTakeConsentUseCase, nil},
		{talksession_usecase.NewEditTalkSessionUseCase, nil},
		{talksession_query.NewBrowseTalkSessionQueryHandler, nil},
		{talksession_query.NewBrowseOpenedByUserQueryHandler, nil},
		{talksession_query.NewBrowseJoinedTalkSessionQueryHandler, nil},
		{talksession_query.NewGetTalkSessionDetailByIDQueryHandler, nil},
		{talksession_query.NewGetConclusionByIDQueryHandler, nil},
		{talksession_query.NewGetRestrictionsQuery, nil},
		{talksession_query.NewHasConsentQuery, nil},
		{talksession.NewIsTalkSessionSatisfiedInteractor, nil},
		{opinion_usecase.NewSubmitOpinionHandler, nil},
		{opinion_usecase.NewReportOpinion, nil},
		{opinion_query.NewGetOpinionsByTalkSessionIDQueryHandler, nil},
		{opinion_query.NewGetOpinionDetailByIDQueryHandler, nil},
		{opinion_query.NewGetOpinionRepliesQueryHandler, nil},
		{opinion_query.NewSwipeOpinionsQueryHandler, nil},
		{opinion_query.NewGetMyOpinionsQueryHandler, nil},
		{opinion_query.NewGetOpinionGroupRatioInteractor, nil},
		{opinion_q.NewGetReportReasons, nil},
		{user_usecase.NewEditHandler, nil},
		{user_usecase.NewRegisterHandler, nil},
		{user_usecase.NewWithdraw, nil},
		{user_query.NewDetailHandler, nil},
		{user_query.NewGetByDisplayIDHandler, nil},
		{vote_usecase.NewVoteHandler, nil},
		{auth_usecase.NewAuthLogin, nil},
		{auth_usecase.NewRevoke, nil},
		{auth_usecase.NewAuthCallback, nil},
		{auth_usecase.NewLoginForDev, nil},
		{auth_usecase.NewDetachAccount, nil},
		{auth_usecase.NewPasswordRegister, nil},
		{auth_usecase.NewPasswordLogin, nil},
		{auth_usecase.NewChangePassword, nil},
		{auth_usecase.NewReactivate, nil},
		{timeline_usecase.NewAddTimeLine, nil},
		{timeline_usecase.NewEditTimeLine, nil},
		{timeline_query.NewGetTimeLine, nil},
		{analysis_query.NewGetAnalysisResultHandler, nil},
		{analysis_query.NewGetReportQueryHandler, nil},
		{report_query.NewGetByTalkSessionQueryInteractor, nil},
		{report_query.NewGetOpinionReportQueryInteractor, nil},
		{report_usecase.NewSolveReportCommandInteractor, nil},
		{report_q.NewGetCountQueryInteractor, nil},
		{image_usecase.NewUploadImageHandler, nil},
		{policy_usecase.NewAcceptPolicy, nil},
		{policy_query.NewCheckConsent, nil},
		{organization_usecase.NewCreateOrganizationInteractor, nil},
		{organization_usecase.NewInviteOrganizationInteractor, nil},
		{organization_usecase.NewInviteOrganizationForUserInteractor, nil},
		{organization_usecase.NewCreateOrganizationAliasUseCase, nil},
		{organization_usecase.NewDeactivateOrganizationAliasUseCase, nil},
		{organization_usecase.NewListOrganizationAliasesUseCase, nil},
		{organization_usecase.NewSwitchOrganizationUseCase, nil},
		{organization_usecase.NewUpdateOrganizationInteractor, nil},
		{organization_query.NewListOrganizationUsersQuery, nil},
		{analysis_usecase.NewApplyFeedbackInteractor, nil},
		{event_processor.NewEventHandlerRegistry, nil},
		{handlers.NewTalkSessionPushNotificationHandler, nil},
		{SetupEventProcessor, nil},
	}
}
