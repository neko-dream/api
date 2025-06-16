import "@typespec/http";
import "@typespec/openapi";
import "../models/opinion.tsp";
import "../models/user.tsp";
import "../models/common.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
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

  @tag("opinion")
  @extension("x-ogen-operation-group", "Opinion")
  @route("/opinions/{opinionID}/report")
  @post
  @summary("意見通報API")
  op reportOpinion(
    @path opinionID: string,
    @multipartBody body: {
      reason?: HttpPart<integer | null>;
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
}