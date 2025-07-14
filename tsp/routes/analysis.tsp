import "@typespec/http";
import "@typespec/openapi";
import "../config/service.tsp";
import "../models/auth.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("analysis")
  @extension("x-ogen-operation-group", "Analysis")
  @route("/report/feedback")
  @post
  @summary("セッションのレポートにフィードバックを適用する")
  op applyFeedbackToReport(
    @multipartBody body: {
      reportID: HttpPart<string>;
      feedbackType: HttpPart<string>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 500;
    @body body: {};
  };
}
