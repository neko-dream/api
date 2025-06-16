import "@typespec/http";
import "@typespec/openapi";
import "../config/service.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("health")
  @extension("x-ogen-operation-group", "Health")
  @route("/health")
  @get
  @summary("ヘルスチェック")
  @useAuth([])
  op health(): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };
}