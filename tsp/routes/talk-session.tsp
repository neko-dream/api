import "@typespec/http";
import "@typespec/openapi";
import "../models/talk-session.tsp";
import "../models/user.tsp";
import "../models/common.tsp";
import "../models/opinion.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
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
      organizationAliasID?: HttpPart<string | null>;
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
}