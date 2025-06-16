import "@typespec/http";

using Http;

namespace kotohiro {
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
    @doc("ユーザーID")
    userID: string;

    @doc("表示ID")
    displayID: string;

    @doc("表示名")
    displayName: string;

    @doc("アイコンURL")
    iconURL: string;

    @doc("最終ログイン日時")
    lastLoginAt: utcDateTime;

    @doc("作成日時")
    createdAt: utcDateTime;

    @doc("更新日時")
    updatedAt: utcDateTime;
  }

  model TalkSessionListResponse {
    TalkSessionStats: TalkSessionStats[];
    TotalCount: int32;
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

  model UserStatsResponse {
    @doc("ユニークアクション数")
    uniqueActionUserCount: int32;

    @doc("登録ユーザー数")
    userCount: int32;

    @doc("セッション数")
    talkSessionCount: int32;

    @doc("日付")
    date: utcDateTime;
  }
}