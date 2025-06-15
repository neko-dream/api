import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

using Http;
using OpenAPI;

@service(#{ title: "Kotohiro API", version: "1.0.0" })
@info(#{ version: "1.0.0" })
@server("http://localhost:3000", "Kotohiro API Server")
@useAuth(CookieAuth)
namespace kotohiro {
  @doc("Cookie-based authentication using JWT tokens stored in secure HTTP-only cookies")
  model CookieAuth is ApiKeyAuth<ApiKeyLocation.cookie, "SessionId">;
  model talkSession {
    /**
     * トークセッションID
     */
    id: string;

    /**
     * テーマ
     */
    theme: string;

    /**
     * 説明
     */
    description?: string | null;

    /**
     * 作成ユーザー
     */
    owner: user;

    /**
     * 作成日時
     */
    createdAt: string;

    /**
     * 終了予定日時
     */
    scheduledEndTime: string;

    /**
     * 位置情報
     */
    location?: location;

    /**
     * 市区町村
     */
    city?: string | null;

    /**
     * 都道府県
     */
    prefecture?: string | null;

    /**
     * サムネ画像
     */
    thumbnailURL?: string | null;

    /**
     * セッションの参加制限
     */
    restrictions: restriction[];

    hideReport: boolean;
  }

  model validationErrorItem {
    /**
     * バリデーションエラーのフィールド
     */
    field: string;

    /**
     * バリデーションエラーメッセージ
     */
    message: string;
  }

  model error {
    code: string;
    message: string;
  }

  model user {
    @minLength(3)
    @maxLength(25)
    @pattern("^[A-Za-z0-9]$")
    displayID: string;

    displayName: string;
    iconURL?: string | null;
  }

  model opinion {
    /**
     * 意見ID
     */
    id: string;

    title?: string;

    /**
     * 意見のテキスト
     */
    @minLength(1) @maxLength(140) content: string;

    /**
     * 親の意見ID。ルートならば無し
     */
    parentID?: string;

    /**
     * 意見投稿主の意見。ルート意見の場合はここには何も入らない
     */
    voteType?: "disagree" | "agree" | "pass" | null;

    /**
     * 画像が返る場合もある
     */
    pictureURL?: string | null;

    /**
     * 参考文献URL
     */
    referenceURL?: string;

    postedAt: string;
    isDeleted: boolean;
  }

  model location {
    /**
     * 緯度
     */
    latitude?: numeric;

    /**
     * 経度
     */
    longitude?: numeric;
  }

  model tokenClaim {
    /**
     * Audience
     */
    aud: string;

    /**
     * 有効期限
     */
    exp: string;

    /**
     * 発行日時
     */
    iat: string;

    /**
     * 発行者
     */
    iss: string;

    /**
     * ユーザID
     */
    `sub`: string;

    /**
     * JWT ID
     */
    jti: string;

    /**
     * ユーザーID
     */
    displayID?: string;

    /**
     * ユーザー名
     */
    displayName?: string;

    /**
     * アイコンURL
     */
    iconURL?: string;

    /**
     * ユーザ登録済みか
     */
    isRegistered: boolean;

    isEmailVerified: boolean;

    /**
     * アカウントの種類。組織がなければ空
     */
    orgType?: integer | null;

    /**
     * パスワードの更新が必要かどうか
     */
    requiredPasswordChange: boolean;
  }

  model offsetPagination {
    totalCount: integer;
    offset: integer;
    limit: integer;
  }

  model userDemographics {
    /**
     * 生年月日
     */
    @summary("20001010") dateOfBirth?: integer | null;

    /**
     * 性別
     */
    @summary("性別") gender?: string | null;

    /**
     * 市町村
     */
    @summary("市町村") city?: string | null;

    /**
     * 都道府県
     */
    @summary("都道府県") prefecture?: string | null;
  }

  model userGroupPosition {
    posX: numeric;
    posY: numeric;
    displayID: string;
    displayName: string;
    iconURL?: string | null;
    groupName: string;
    groupID: integer;

    /**
     * 境界ポイントのインデックス
     */
    perimeterIndex?: integer;
  }

  model Report {
    talkSessionID: string;

    /**
     * レポート本文
     */
    content: string;
  }

  model conclusion {
    /**
     * 作成ユーザー
     */
    user: user;

    /**
     * 結論本文
     */
    content: string;
  }

  model actionItem {
    ActionItemID: string;
    Sequence: integer;
    Content: string;
    Status: string;
    CreatedAt: string;
    UpdatedAt: string;
  }

  model restriction {
    key: string;
    description: string;

    /**
     * 依存しているrestriction
     */
    dependsOn?: string[];
  }

  model policyConsentStatus {
    /**
     * 最新ポリシーのバージョン
     */
    policyVersion: string;

    /**
     * 同意した日時
     */
    consentedAt?: string | null;

    /**
     * 同意したか
     */
    consentGiven: boolean;
  }

  model reportReason {
    /**
     * 1
     */
    reasonID: integer;

    /**
     * 不適切な内容
     */
    reason: string;
  }

  model opinionGroupRatio {
    agreeCount: integer;
    disagreeCount: integer;
    passCount: integer;
    groupID: integer;
    groupName: string;
  }

  model reportDetail {
    opinion: opinion;

    /**
     * 作成ユーザー
     */
    user: user;

    status: "unsolved" | "deleted" | "hold";
    reasons: {
      reason: string;
      content?: string | null;
    }[];

    /**
     * この意見が通報を受けた回数
     */
    reportCount: integer;
  }

  model organization {
    /**
     * 組織ID
     */
    ID: string;

    /**
     * 組織名
     */
    Name: string;

    /**
     * 組織のタイプ
     */
    Type: integer;

    RoleName: string;

    /**
     * ロール
     */
    Role: integer;
  }

  model success {
    message: string;
  }

  // ==================== Models from main.tsp ====================

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

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/{opinionID}")
  @get
  @summary("意見詳細")
  op getOpinionDetail2(@path opinionID: string): Body<{
    opinion: opinion;
    user: user;
    myVoteType?: "agree" | "disagree" | "pass" | null;
  }> | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/{opinionID}/replies")
  @get
  @summary("意見に対するリプライ意見一覧")
  op opinionComments2(
    /**
     * 親意見のID
     */
    @path opinionID: string,
  ): Body<{
    opinions: {
      opinion: opinion;
      user: user;
      myVoteType?: "agree" | "disagree" | "pass" | null;
    }[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("vote")
  @extension("x-ogen-operation-group", "Vote")
  @route("/opinions/{opinionID}/votes")
  @post
  @summary("意思表明API")
  op vote2(
    /**
     * 意見のID
     */
    @path opinionID: string,

    @multipartBody body: {
      voteStatus: HttpPart<"agree" | "disagree" | "pass" | null>;
    },
  ): Body<opinion[]> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/{opinionID}/report")
  @post
  @summary("意見通報API")
  op reportOpinion(
    @path opinionID: string,
    @multipartBody body: {
      reason?: HttpPart<integer>;
      content?: HttpPart<string | null>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * セッション作成者しか取得できない
   */
  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/{opinionID}/reports")
  @get
  @summary("意見に対する通報取得")
  op getOpinionReports(@path opinionID: string): reportDetail | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/{opinionID}/reports/solve")
  @post
  @summary("通報を解決")
  op solveOpinionReport(
    @path opinionID: string,
    @multipartBody body: {
      action: HttpPart<"deleted" | "hold">;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/{opinionID}/analysis")
  @get
  @summary("意見に投票したグループごとの割合")
  op getOpinionAnalysis(@path opinionID: string): Body<opinionGroupRatio[]> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * parentOpinionIDがなければルートの意見として投稿される
   * parentOpinionIDがない場合はtalkSessionIDが必須
   *
   * セッション管理者はisSeedをtrueにするとシード意見として投稿できる
   */
  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions")
  @post
  @summary("セッションに対して意見投稿 or 意見に対するリプライ")
  op postOpinionPost2(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      talkSessionID?: HttpPart<string | null>;
      parentOpinionID?: HttpPart<string | null>;
      title?: HttpPart<string | null>;
      opinionContent: HttpPart<string>;
      referenceURL?: HttpPart<string | null>;
      picture?: HttpPart<bytes>;
      isSeed?: HttpPart<boolean | null>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/report_reasons")
  @get
  @summary("意見への通報理由一覧")
  op getOpinionReportReasons(): Body<reportReason[]> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * 特定のセッションで満たしていない条件があれば返す
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/restrictions")
  @get
  @summary("セッションで満たしていない制限")
  op getTalkSessionRestrictionSatisfied(
    @path talkSessionID: string,
  ): Body<restriction[]> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}")
  @get
  @summary("トークセッションの詳細")
  op getTalkSessionDetail(@path talkSessionID: string): talkSession | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}")
  @put
  @summary("セッション編集")
  op editTalkSession(
    @path talkSessionID: string,
    @body body: {
      theme: string;
      scheduledEndTime: utcDateTime;
      latitude?: numeric | null;
      longitude?: numeric;
      prefecture?: string | null;
      city?: string | null;
      description?: string | null;
      thumbnailURL?: string;
    },
  ): talkSession | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/talksessions/{talkSessionID}/opinions")
  @get
  @summary("セッションに対する意見一覧")
  op getOpinionsForTalkSession(
    @path talkSessionID: string,
    @query(#{ explode: true }) sort?:
      | "latest"
      | "mostReplies"
      | "oldest"
      | null,
    @query(#{ explode: true }) limit?: integer | null,
    @query(#{ explode: true }) offset?: integer | null,
    @query(#{ explode: true }) seed?: boolean | null,
  ): Body<{
    opinions: {
      opinion: opinion;
      user: user;
      replyCount: integer;
      myVoteType?: "pass" | "disagree" | "agree" | null;
    }[];
    pagination: {
      totalCount: integer;
    };
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * セッションの中からまだ投票していない意見をランダムに取得する
   * remainingCountは取得した意見を含めてスワイプできる意見の総数を返す
   */
  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/talksessions/{talkSessionID}/swipe_opinions")
  @get
  @summary("スワイプ用のエンドポイント")
  op swipe_opinions(
    @path talkSessionID: string,
    @query(#{ explode: true }) limit?: integer | null = 10,
  ): Body<{
    opinions: {
      replyCount: integer;
      opinion: opinion;
      user: user;
    }[];
    remainingCount: integer;
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions")
  @get
  @summary("セッション一覧")
  op getTalkSessionList(
    /**
     * 1ページあたりの要素数
     */
    @query(#{ explode: true }) limit?: integer | null,

    /**
     * どの要素から始めるか
     */
    @query(#{ explode: true }) offset?: integer | null,

    @query(#{ explode: true }) theme?: string | null,
    @query(#{ explode: true }) status?: "open" | "finished" | null = "open",

    @extension(
      "x-apidog",
      #{
        enumDescriptions: #{
          latest: "",
          oldest: "",
          mostReplies: "",
          nearest: "",
        },
      }
    )
    @query(#{ explode: true })
    sortKey?: "latest" | "oldest" | "mostReplies" | "nearest",

    @query(#{ explode: true }) latitude?: numeric | null,
    @query(#{ explode: true }) longitude?: numeric | null,
  ): Body<{
    talkSessions: {
      talkSession: talkSession;
      opinionCount: integer;
    }[];
    pagination: offsetPagination;
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  /**
   * ## サムネイル画像について
   * - `Description中に出てくる画像で一番最初のものを使用`。
   * - 画像自体は`POST /images`でサーバにポストしたものを使用してください。
   *
   * ## 投稿制限のキーについて
   * restrictionsに値を入れると一定のデモグラ情報を登録していない限り、セッションへの投稿が制限されるようにできる。
   * restrictionsには [GET /talksessions/restrictions](https://app.apidog.com/link/project/674502/apis/api-14271260) より取れるkeyをカンマ区切りで入力してください。
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions")
  @post
  @summary("セッション作成")
  op createTalkSession(
    @multipartBody body: {
      theme: HttpPart<string>;
      scheduledEndTime: HttpPart<utcDateTime>;
      latitude?: HttpPart<numeric | null>;
      longitude?: HttpPart<numeric | null>;
      city?: HttpPart<string | null>;
      prefecture?: HttpPart<string | null>;
      description?: HttpPart<string | null>;
      thumbnailURL?: HttpPart<string | null>;
      restrictions?: HttpPart<string[]>;
    },
  ): Body<{
    id: string;
    theme: string;
    description?: string | null;
    owner: user;
    createdAt: string;
    scheduledEndTime: string;
    location?: location;
    city?: string | null;
    prefecture?: string | null;
    thumbnailURL?: string | null;
    restrictions: restriction[];
    hideReport: boolean;
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/analysis")
  @get
  @summary("分析結果一覧")
  op talkSessionAnalysis(@path talkSessionID: string): Body<{
    myPosition?: userGroupPosition;
    positions: userGroupPosition[];
    groupOpinions: {
      groupName: string;
      groupID: integer;
      opinions: {
        opinion: opinion;
        user: user;
        agreeCount: integer;
        disagreeCount: integer;
        passCount: integer;
      }[];
    }[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/report")
  @get
  @summary("セッションレポートを返す")
  op getTalkSessionReport(@path talkSessionID: string): Body<{
    report?: string | null;
  }> | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/reports")
  @get
  @summary("通報一覧")
  op getReportsForTalkSession(
    @path talkSessionID: string,
    @query(#{ explode: true }) status?: "unsolved" | "deleted" | "hold",
  ): Body<{
    reports: reportDetail[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/reports/count")
  @get
  @summary("通報件数")
  op getTalkSessionReportCount(
    @path talkSessionID: string,
    @query(#{ explode: true }) status: "unsolved" | "deleted" | "hold",
  ): Body<{
    count: integer;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("timeline")
  @extension("x-ogen-operation-group", "Timeline")
  @route("/talksessions/{talkSessionID}/timelines/{actionItemID}")
  @put
  @summary("タイムライン編集")
  op editTimeLine(
    @path talkSessionID: string,
    @path actionItemID: string,
    @multipartBody body: {
      content?: HttpPart<string | null>;
      status?: HttpPart<string | null>;
    },
  ): actionItem | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/conclusion")
  @get
  @summary("結論取得")
  op getConclusion(@path talkSessionID: string): conclusion | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * 結論（conclusion）はセッションが終了した後にセッっションの作成者が投稿できる文章。
   * セッションの流れやグループの分かれ方などに対するセッション作成者の感想やそれらの意見を受け、これからの方向性などを記入する。
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/conclusion")
  @post
  @summary("結論投稿")
  op postConclusion(
    @path talkSessionID: string,
    @multipartBody body: {
      content: HttpPart<string>;
    },
  ): conclusion | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * タイムラインはセッション終了後にセッション作成者が設定できるその後の予定を知らせるもの
   */
  @tag("timeline")
  @extension("x-ogen-operation-group", "Timeline")
  @route("/talksessions/{talkSessionID}/timelines")
  @get
  @summary("タイムライン取得")
  op getTimeLine(@path talkSessionID: string): Body<{
    items: actionItem[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("timeline")
  @extension("x-ogen-operation-group", "Timeline")
  @route("/talksessions/{talkSessionID}/timeline")
  @post
  @summary("タイムラインアイテム追加")
  op postTimeLineItem(
    @path talkSessionID: string,
    @multipartBody body: {
      content: HttpPart<string>;
      status: HttpPart<string>;
      parentActionItemID?: HttpPart<string | null>;
    },
  ): Body<{
    ActionItemID: string;
    Sequence: integer;
    Content: string;
    Status: string;
    CreatedAt: string;
    UpdatedAt: string;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/consent")
  @get
  @summary("セッションに同意しているか")
  op hasConsent(@path talkSessionID: string): Body<{
    hasConsent: boolean;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/consent")
  @post
  @summary("セッションへの同意")
  op consentTalkSession(@path talkSessionID: string): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * セッションの投稿制限に使用できるキーの一覧を返す
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/restrictions")
  @get
  @summary("セッションで指定可能な制限一覧")
  op getTalkSessionRestrictionKeys(): Body<restriction[]> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/change")
  @put
  @summary("パスワード変更")
  op changePassword(
    /**
     * 古いパスワード
     */
    @query(#{ explode: true }) old_password: string,

    /**
     * 新たなパスワード
     */
    @query(#{ explode: true }) new_password: string,
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("user")
  @extension("x-ogen-operation-group", "User")
  @route("/user")
  @get
  @summary("ユーザー情報の取得")
  op get_user_info(): Body<{
    user: user;
    demographics: userDemographics;
    email?: string | null;
  }> | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("user")
  @extension("x-ogen-operation-group", "user")
  @route("/user")
  @post
  @summary("ユーザー作成")
  op registerUser(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      displayName: HttpPart<string>;
      displayID: HttpPart<string>;
      icon?: HttpPart<bytes>;
      dateOfBirth?: HttpPart<string>;
      gender?: HttpPart<"男性" | "女性" | "その他" | "回答しない" | null>;
      prefecture?: HttpPart<string>;
      city?: HttpPart<string | null>;
      email?: HttpPart<string | null>;
    },
  ): Body<{
    @minLength(3)
    @maxLength(25)
    @pattern("^[A-Za-z0-9]$")
    displayID: string;

    displayName: string;
    iconURL?: string | null;
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("user")
  @extension("x-ogen-operation-group", "user")
  @route("/user")
  @put
  @summary("ユーザー情報の変更")
  op editUserProfile(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      displayName?: HttpPart<string | null>;
      icon?: HttpPart<bytes>;
      deleteIcon?: HttpPart<boolean | null>;
      dateOfBirth?: HttpPart<string>;
      gender?: HttpPart<"男性" | "女性" | "その他" | "回答しない" | null>;
      city?: HttpPart<string | null>;
      prefecture?: HttpPart<string | null>;
      email?: HttpPart<string | null>;
    },
  ): Body<{
    @minLength(3)
    @maxLength(25)
    @pattern("^[A-Za-z0-9]$")
    displayID: string;

    displayName: string;
    iconURL?: string | null;
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("user")
  @extension("x-ogen-operation-group", "User")
  @route("/opinions/histories")
  @get
  @summary("今までに投稿した異見")
  op opinionsHistory(
    /**
     * ソートきー
     */
    @query(#{ explode: true }) sort?:
      | "latest"
      | "mostReplies"
      | "oldest"
      | null,

    @query(#{ explode: true }) limit?: integer,
    @query(#{ explode: true }) offset?: integer,
  ): Body<{
    opinions: {
      user: user;
      opinion: opinion;
      replyCount: integer;
    }[];
    pagination: {
      totalCount: integer;
    };
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/opened")
  @get
  @summary("自分が開いたセッション一覧")
  op getOpenedTalkSession(
    @query(#{ explode: true }) limit?: integer = 10,
    @query(#{ explode: true }) offset?: integer,

    /**
     * テーマ
     */
    @query(#{ explode: true }) theme?: string,

    @extension("x-apidog", #{ enumDescriptions: #{ finished: "", open: "" } })
    @query(#{ explode: true })
    status?: "finished" | "open" | null,
  ): Body<{
    talkSessions: {
      talkSession: talkSession;
      opinionCount: integer;
    }[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("user")
  @extension("x-ogen-operation-group", "User")
  @route("/talksessions/histories")
  @get
  @summary("リアクション済みのセッション一覧")
  op sessionsHistory(
    @query(#{ explode: true }) limit?: integer,
    @query(#{ explode: true }) offset?: integer,

    /**
     * テーマ
     */
    @query(#{ explode: true }) theme?: string | null,

    @query(#{ explode: true }) status?: "open" | "finished" | null,
  ): Body<{
    pagination: offsetPagination;
    talkSessions: {
      talkSession: talkSession;
      opinionCount: integer;
    }[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/token/info")
  @get
  @summary("JWTの内容を返してくれる")
  op oauth_token_info(): Body<{
    aud: string;
    exp: string;
    iat: string;
    iss: string;
    `sub`: string;
    jti: string;
    displayID?: string;
    displayName?: string;
    iconURL?: string;
    isRegistered: boolean;
    isEmailVerified: boolean;
    orgType?: integer | null;
    requiredPasswordChange: boolean;
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: {
      code: string;
      message: string;
    };
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/{provider}/callback")
  @get
  @summary("Auth Callback")
  @useAuth([])
  op oauth_callback(
    @path provider: string,
    @query(#{ explode: true }) code: string,

    /**
     * OAuth State from Query
     */
    @query(#{ explode: true }) state: string,
  ): {
    @statusCode statusCode: 302;
    @header("Location") Location: string;

    /**
     * SessionID
     */
    @header("Set-Cookie") SetCookie: string[];

    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/dev/login")
  @get
  @summary("開発用登録/ログイン")
  @useAuth([])
  op devAuthorize(
    @query(#{ explode: true }) redirect_url: string,

    /**
     * devのみで使用するsubjectです。ここで指定した値はログインした後も確認できないため覚えておいてください。同じ値を指定すると同じアカウントにログインできます。
     */
    @query(#{ explode: true }) id: string,
  ): {
    @statusCode statusCode: 302;
    @header("Location") Location?: string;
    @header("Set-Cookie") SetCookie?: string[];
    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/{provider}/login")
  @get
  @summary("ログイン")
  @useAuth([])
  op authorize(
    @path provider: "google" | "line",

    /**
     * ログイン後にリダイレクトするURL
     */
    @query(#{ explode: true }) redirect_url: string,

    /**
     * 登録していなかった場合に飛ばすURL
     */
    @query(#{ explode: true }) registration_url?: string,
  ): {
    @statusCode statusCode: 302;

    /**
     * IDPのログインページ
     */
    @header("Location") Location: string;

    /**
     * OAuth2.0 State
     */
    @header("Set-Cookie") SetCookie?: string[];

    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * そのアカウントには再度ログインできなくなります。ログインしたければ言ってね！
   */
  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/dev/detach")
  @delete
  @summary("**開発用** 現在のアカウントを切り離す")
  @useAuth([])
  op authAccountDetach(): {
    @header("Set-Cookie") SetCookie?: string[];
    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/revoke")
  @post
  @summary("トークンを失効（ログアウト）")
  op oauth_token_revoke(): {
    @statusCode statusCode: 204;

    /**
     * Cookie削除用のSet-Cookie
     */
    @header("Set-Cookie") SetCookie: string[];

    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/login")
  @post
  @summary("パスワードによるログイン")
  @useAuth([])
  op passwordLogin(
    @body body: {
      id_or_email: string;
      password: string;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/register")
  @post
  @summary("パスワードによる登録（devのみ）")
  @useAuth([])
  op passwordRegister(
    @multipartBody body: {
      id: HttpPart<string>;
      password: HttpPart<string>;
      email: HttpPart<string>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("user")
  @extension("x-ogen-operation-group", "test")
  @route("/test")
  @get
  @summary("OpenAPIテスト用")
  op test(): Body<{
    optInt?: integer | null;
    optNilInt?: integer;
    optNilBool?: boolean | null;
    optBool?: boolean;
    optUrl?: url;

    @format("url")
    optNilUrl?: string | null;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("test")
  @extension("x-ogen-operation-group", "Test")
  @route("/test/dummy")
  @post
  @summary("mudai")
  op dummiInit(): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * 画像を投稿してURLを返すAPI
   */
  @tag("image")
  @extension("x-ogen-operation-group", "Image")
  @route("/images")
  @post
  @summary("画像投稿")
  op postImage(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      image: HttpPart<bytes>;
    },
  ): Body<{
    url: string;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("policy")
  @extension("x-ogen-operation-group", "Policy")
  @route("/policy/consent")
  @get
  @summary("最新のポリシーに同意したかを取得")
  op getPolicyConsentStatus(): policyConsentStatus | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("policy")
  @extension("x-ogen-operation-group", "Policy")
  @route("/policy/consent")
  @post
  @summary("最新のポリシーに同意する")
  op policyConsent(
    @multipartBody body: {
      policyVersion: HttpPart<string>;
    },
  ): policyConsentStatus | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations")
  @get
  @summary("所属組織一覧")
  op getOrganizations(): Body<{
    organizations: organization[];
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * 組織を作成できる。
   * これを作れるユーザーはDBを直接叩いて作るしかない。
   *
   * OrgType
   * - 1: 通常
   * - 2: 自治体
   * - 3: 議員
   */
  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations")
  @post
  @summary("組織作成（運営ユーザーのみ）")
  op createOrganizations(
    @multipartBody body: {
      name: HttpPart<string>;
      orgType?: HttpPart<integer | null>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * Role
   * - 1: Member
   * - 2: Admin
   * - 3: Owner
   */
  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/{organizationID}/invite")
  @post
  @summary("組織ユーザー招待（運営ユーザーのみ）")
  op inviteOrganization(
    @path organizationID: string,
    @multipartBody body: {
      role: HttpPart<integer>;
      email: HttpPart<string>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("organization")
  @extension("x-ogen-operation-group", "Organization")
  @route("/organizations/{organizationID}/invite_user")
  @post
  @summary("組織にユーザーを追加")
  op inviteOrganizationForUser(
    @path organizationID: string,
    @multipartBody body: {
      role: HttpPart<integer>;
      displayID: HttpPart<string | null>;
    },
  ): Body<{
    success: boolean;
  }> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("health")
  @extension("x-ogen-operation-group", "Health")
  @route("/health")
  @get
  @summary("ヘルスチェック")
  @useAuth([])
  op health(): Body<{}> | {
    @statusCode statusCode: 400;
    @multipartBody body: {};
  } | {
    @statusCode statusCode: 500;
    @multipartBody body: {};
  };

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
