import "@typespec/http";
import "./user.tsp";

using Http;

namespace kotohiro {
  model ManagedUser {
    permissions: UserPermissions;
    userID: string;
    displayID: string;
    displayName: string;
    email?: string;
    emailVerified: boolean;
    passwordSetupRequired: boolean;
    providers: string[];
    createdAt: utcDateTime;
    updatedAt: utcDateTime;
    withdrawalDate?: utcDateTime;
  }

  model ManagedUserList {
    users: ManagedUser[];
    totalCount: int64;
  }

  model InviteManagedUserRequest {
    @minLength(1)
    @maxLength(50)
    displayName: string;
    @minLength(3)
    @maxLength(254)
    email: string;
  }

  model ManagedUserInvitation {
    userID: string;
    invitationSent: boolean;
  }

  model TalkSessionForManage {
    talkSessionID: string;
    theme: string;
    description: string;
    owner: UserForManage;
    scheduledEndTime: utcDateTime;
    city?: string;
    prefecture?: string;
    thumbnailURL: string;
    hidden: boolean;
    updatedAt: string;
    createdAt: string;

    /** ホームの一覧に表示しない */
    hideTop: boolean;

    /** レポートを表示しない */
    hideReport: boolean;

    /** 分析を無効にする */
    disableAnalysis: boolean;

    /** 参加制限のキー（GET /talksessions/restrictions の key） */
    restrictions: string[];

    /** アンケートがあるか */
    hasSurvey: boolean;

    /** 意見を投稿したか投票したユーザーの数。参加制限を変えるときの警告に使う */
    participantCount: int32;
  }

  model TalkSessionEditForManage {
    /** 1〜100文字 */
    theme: string;

    /** 本文（HTML）。40000文字まで */
    description: string;

    /** 変更する場合は未来の日時にする。変更しない場合は現在の値をそのまま送る */
    scheduledEndTime: utcDateTime;

    /** 空文字は未設定として扱う */
    prefecture: string;

    /** 空文字は未設定として扱う */
    city: string;

    hideTop: boolean;
    hideReport: boolean;
    disableAnalysis: boolean;

    /** 参加制限のキー。survey.answered はアンケートがあるセッションにだけ付けられる */
    restrictions: string[];
  }

  model TalkSessionOpinionsForManage {
    totalCount: int32;
    opinions: OpinionForManage[];
  }

  model OpinionForManage {
    opinionID: string;
    content: string;
    pictureURL?: string;
    referenceURL?: string;

    /** 返信の場合は返信先の意見ID */
    parentOpinionID?: string;

    author: UserForManage;
    createdAt: utcDateTime;
    agreeCount: int32;
    disagreeCount: int32;
    passCount: int32;

    /** 通報の件数（すべての状態） */
    reportCount: int32;

    /** 未対応の通報の件数 */
    unsolvedReportCount: int32;

    /** 通報の対応で削除された。削除済みの通報が1件でもあると、一般の意見一覧・返信・詳細では「運営により削除されました」と表示される */
    deleted: boolean;
  }

  model TalkSessionStats {
    talkSessionID: string;
    theme: string;
    description: string;
    owner: UserForManage;
    scheduledEndTime: utcDateTime;
    city?: string;
    prefecture?: string;
    thumbnailURL: string;

    /** レポートを表示しない（hideReport と同じ値） */
    hidden: boolean;

    /** ホームの一覧に表示しない */
    hideTop: boolean;

    updatedAt: string;
    createdAt: string;
    opinionCount: int32;
    opinionUserCount: int32;
    voteCount: int32;
    voteUserCount: int32;
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
    talkSessionStats: TalkSessionStats[];
    totalCount: int32;
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
