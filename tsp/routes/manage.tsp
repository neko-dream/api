import "@typespec/http";
import "@typespec/openapi";
import "../models/manage.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @route("/v1/manage")
  namespace KotohiroManageService {
    @route("/talksessions")
    @tag("manage")
    interface TalkSessions {
      @route("/list")
      @operationId("getTalkSessionListManage")
      @extension("x-ogen-operation-group", "manage")
      @get
      getTalkSessionList(
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
        @multipartBody body: kotohiro.ToggleReportVisibilityRequest,
      ): kotohiro.ToggleReportVisibilityResponse;

      @route("/regenerate")
      @operationId("manageRegenerateManage")
      @extension("x-ogen-operation-group", "manage")
      @post
      regenerateAnalysis(
        @header contentType: "application/x-www-form-urlencoded",
        @path talkSessionID: string,
        @multipartBody body: kotohiro.RegenerateRequest,
      ): kotohiro.RegenerateResponse;
    }

    @route("/users")
    @tag("manage")
    interface UserStats {
      @route("/list")
      @extension("x-ogen-operation-group", "manage")
      @operationId("getUserListManage")
      @get
      getUserList(
        @query offset?: int32,
        @query limit?: int32,
        @query search?: string,
        @query orderBy?:
          | "createdAt"
          | "updatedAt"
          | "displayName"
          | "lastLoginAt",
        @query order?: "asc" | "desc",
      ): kotohiro.UserForManage[];

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