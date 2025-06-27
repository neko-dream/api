import "@typespec/http";
import "@typespec/openapi";
import "../models/opinion.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
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
      voteStatus: HttpPart<string>;
    },
  ): Body<Opinion[]> | {
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
}
