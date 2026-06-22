import "@typespec/http";
import "@typespec/openapi";
import "../models/talk-session-survey.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  /**
   * セッションのアンケート定義を取得（参加前に質問を表示するため）
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/survey")
  @get
  @summary("アンケート定義取得")
  @useAuth([])
  op getTalkSessionSurvey(
    @path talkSessionID: string,
  ): TalkSessionSurvey | {
    @statusCode statusCode: 404;
    @body body: {};
  };

  /**
   * セッションにアンケートを作成（オーナーのみ、セッションあたり 1 件まで）
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/survey")
  @post
  @summary("アンケート作成")
  op createTalkSessionSurvey(
    @path talkSessionID: string,
    @body body: {
      questions: TalkSessionSurveyQuestionInput[];
    },
  ): TalkSessionSurvey | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 403;
    @body body: {};
  } | {
    @statusCode statusCode: 404;
    @body body: {};
  } | {
    @statusCode statusCode: 409;
    @body body: {
      code: string;
      message: string;
    };
  };

  /**
   * アンケートを編集（オーナーのみ）。
   * 既に回答がある質問・選択肢は soft delete 経由でのみ非表示化可能。
   * questionID / choiceID 指定で既存項目を編集、省略で新規追加。
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/survey")
  @put
  @summary("アンケート編集")
  op updateTalkSessionSurvey(
    @path talkSessionID: string,
    @body body: {
      questions: TalkSessionSurveyQuestionInput[];
      /**
       * 削除する質問ID（回答があれば soft delete、なければ hard delete）
       */
      deletedQuestionIDs?: string[];
      /**
       * 削除する選択肢ID
       */
      deletedChoiceIDs?: string[];
    },
  ): TalkSessionSurvey | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 403;
    @body body: {};
  } | {
    @statusCode statusCode: 404;
    @body body: {};
  };

  /**
   * アンケートに回答を送信（1 ユーザー 1 送信、修正不可）。
   * 成功すると survey.answered restriction が自動的に満たされた状態になる。
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/survey/submissions")
  @post
  @summary("アンケート回答送信")
  op submitTalkSessionSurvey(
    @path talkSessionID: string,
    @body body: {
      answers: TalkSessionSurveyAnswerInput[];
    },
  ): Body<{
    submissionID: string;
  }> | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 404;
    @body body: {};
  } | {
    @statusCode statusCode: 409;
    @body body: {
      code: string;
      message: string;
    };
  };

  /**
   * セッションオーナー向け：回答一覧を取得
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/survey/submissions")
  @get
  @summary("アンケート回答一覧（オーナー）")
  op listTalkSessionSurveySubmissions(
    @path talkSessionID: string,
    @query limit?: int32,
    @query offset?: int32,
  ): Body<{
    submissions: TalkSessionSurveySubmission[];
    totalCount: int32;
  }> | {
    @statusCode statusCode: 403;
    @body body: {};
  } | {
    @statusCode statusCode: 404;
    @body body: {};
  };

  /**
   * セッションオーナー向け：質問ごとに集計した結果を取得
   */
  @tag("talk_session")
  @extension("x-ogen-operation-group", "TalkSession")
  @route("/talksessions/{talkSessionID}/survey/summary")
  @get
  @summary("アンケート集計結果（オーナー）")
  op getTalkSessionSurveySummary(
    @path talkSessionID: string,
  ): TalkSessionSurveySummary | {
    @statusCode statusCode: 403;
    @body body: {};
  } | {
    @statusCode statusCode: 404;
    @body body: {};
  };
}
