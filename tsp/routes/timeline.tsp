import "@typespec/http";
import "@typespec/openapi";
import "../models/talk-session.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
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
}