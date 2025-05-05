import "@typespec/http";
import "@typespec/openapi";

using TypeSpec.Http;
using TypeSpec.OpenAPI;

@service(#{
  title: "Kotohiro Manage API Server",
  version: "1.0.0",
})
@server("http://localhost:3000", "Kotohiro Manage API Server")
@useAuth(ApiKeyAuth<ApiKeyLocation.cookie, "SessionId">)
@route("/v1/manage")
namespace KotohiroManageService {
  model TalkSessionForManage {
    TalkSessionID: string;
    Theme: string;
    Description: string;
    Owner: UserForManage;
    ScheduledEndTime: utcDateTime;
    City?: string;
    Prefecture?: string;
    ThumbnailURL: string;
    Hidden: boolean;
    UpdatedAt: string;
    CreatedAt: string;
  }

  model TalkSessionStats {
    TalkSessionID: string;
    Theme: string;
    Description: string;
    Owner: UserForManage;
    ScheduledEndTime: utcDateTime;
    City?: string;
    Prefecture?: string;
    ThumbnailURL: string;
    Hidden: boolean;
    UpdatedAt: string;
    CreatedAt: string;
    OpinionCount: int32;
    OpinionUserCount: int32;
    VoteCount: int32;
    VoteUserCount: int32;
  }

  model UserForManage {
    UserID: string;
    DisplayID: string;
    DisplayName: string;
    IconURL: string;
  }

  @route("/talksessions")
  @tag("manage")
  interface TalkSessions {

    @route("/list")
    @operationId("getTalkSessionListManage")
    @extension("x-ogen-operation-group", "manage")
    @get getTalkSessionList(
      @query status?: "active" | "inactive",
      @query page?: int32,
      @query limit?: int32,
    ): TalkSessionStats[];


    @route("/{talkSessionID}")
    @operationId("getTalkSessionManage")
    @extension("x-ogen-operation-group", "manage")
    @get getTalkSession(
      @path talkSessionID: string
    ): TalkSessionForManage;
  }

  model AnalysisReportResponse {
    @doc("レポート本文")
    report?: string;
  }

  model RegenerateResponse {
    @doc("ステータス")
    status: string;

    @doc("メッセージ")
    message: string;
  }

  model ToggleReportVisibilityResponse {
    @doc("ステータス")
    status: string;

    @doc("メッセージ")
    message: string;
  }

  model ToggleReportVisibilityRequest {
    @doc("非表示にするかどうか")
    hidden: boolean;
  }

  model RegenerateRequest {
    @doc("再生成するタイプ")
    type: "report" | "group" | "image";
  }


  @route("/talksessions/{talkSessionID}/analysis")
  @tag("manage")
  interface Analysis {

    @route("/report")
    @operationId("getAnalysisReportManage")
    @extension("x-ogen-operation-group", "manage")
    @get getAnalysisReport(
      @path talkSessionID: string
    ): AnalysisReportResponse;


    @route("/report")
    @extension("x-ogen-operation-group", "manage")
    @operationId("toggleReportVisibilityManage")
    @post toggleReportVisibility(
      @header contentType: "application/x-www-form-urlencoded",
      @path talkSessionID: string,
      @body body: ToggleReportVisibilityRequest,
    ): ToggleReportVisibilityResponse;


    @route("/regenerate")
    @operationId("manageRegenerateManage")
    @extension("x-ogen-operation-group", "manage")
    @post regenerateAnalysis(
      @header contentType: "application/x-www-form-urlencoded",
      @path talkSessionID: string,
      @body body: RegenerateRequest,
    ): RegenerateResponse;
  }
}

