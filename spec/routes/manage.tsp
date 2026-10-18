import "@typespec/http";
import "@typespec/openapi";
import "../config/service.tsp";
import "../models/manage.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @route("/v1/manage")
  @useAuth(AdminCookieAuth)
  namespace KotohiroManageService {
    @route("/talksessions")
    @tag("manage")
    interface TalkSessions {
      @route("/list")
      @operationId("getTalkSessionListManage")
      @extension("x-ogen-operation-group", "manage")
      @get
      getTalkSessionListManage(
        @query status?: "active" | "inactive",
        @query offset?: int32,
        @query limit?: int32,
      ): kotohiro.TalkSessionListResponse;

      @route("/{talkSessionID}")
      @operationId("getTalkSessionManage")
      @extension("x-ogen-operation-group", "manage")
      @get
      getTalkSession(
        @path talkSessionID: string,
      ): kotohiro.TalkSessionForManage;

      /** 管理者がセッションの内容・表示設定・参加制限を編集する。変更は監査ログに残す。 */
      @route("/{talkSessionID}")
      @operationId("editTalkSessionManage")
      @extension("x-ogen-operation-group", "manage")
      @put
      editTalkSession(
        @path talkSessionID: string,
        @body body: kotohiro.TalkSessionEditForManage,
      ): kotohiro.TalkSessionForManage;
    }

    @route("/talksessions/{talkSessionID}/analysis")
    @tag("manage")
    interface Analysis {
      @route("/report")
      @operationId("getAnalysisReportManage")
      @extension("x-ogen-operation-group", "manage")
      @get
      getAnalysisReport(
        @path talkSessionID: string,
      ): kotohiro.AnalysisReportResponse;

      @route("/report")
      @extension("x-ogen-operation-group", "manage")
      @operationId("toggleReportVisibilityManage")
      @post
      toggleReportVisibility(
        @header contentType: "application/x-www-form-urlencoded",
        @path talkSessionID: string,
        @body body: kotohiro.ToggleReportVisibilityRequest,
      ): kotohiro.ToggleReportVisibilityResponse;

      @route("/regenerate")
      @operationId("manageRegenerateManage")
      @extension("x-ogen-operation-group", "manage")
      @post
      regenerateAnalysis(
        @header contentType: "application/x-www-form-urlencoded",
        @path talkSessionID: string,
        @body body: kotohiro.RegenerateRequest,
      ): kotohiro.RegenerateResponse;
    }

    @route("/users")
    @tag("manage")
    interface UserStats {
      @extension("x-ogen-operation-group", "manage")
      @operationId("createUserManage")
      @post
      @summary("メールログイン用ユーザーを作成し、初回パスワード設定メールを送信する")
      createUser(@body body: InviteManagedUserRequest): {
        @statusCode statusCode: 201;
        @body body: ManagedUserInvitation;
      };

      @route("/{userID}/invitation")
      @extension("x-ogen-operation-group", "manage")
      @operationId("resendUserInvitationManage")
      @post
      resendUserInvitation(@path userID: string): ManagedUserInvitation;

      @route("/{userID}/permissions")
      @extension("x-ogen-operation-group", "manage")
      @operationId("setUserPermissionsManage")
      @put
      setUserPermissions(@path userID: string, @body permissions: UserPermissions): UserPermissions;

      @route("/list")
      @extension("x-ogen-operation-group", "manage")
      @operationId("getUserListManage")
      @get
      getUserList(
        @minValue(0)
        @query offset?: int32,
        @minValue(1)
        @maxValue(100)
        @query limit?: int32,
        @maxLength(200)
        @query search?: string,
      ): kotohiro.ManagedUserList;

      @route("/{userID}")
      @extension("x-ogen-operation-group", "manage")
      @operationId("getUserDetailManage")
      @get
      getUserDetail(@path userID: string): Body<kotohiro.ManagedUser> | {
        @statusCode statusCode: 404;
        @body body: {};
      };

      @route("/stats/list")
      @extension("x-ogen-operation-group", "manage")
      @operationId("getUserStatsListManage")
      @get
      getUserStats(
        @query range: string,
        @query offset?: int32,
        @query limit?: int32,
      ): kotohiro.UserStatsResponse[];

      @route("/stats/total")
      @extension("x-ogen-operation-group", "manage")
      @operationId("getUserStatsTotalManage")
      @get
      getUserStatsTotal(): kotohiro.UserStatsResponse;
    }
  }
}
