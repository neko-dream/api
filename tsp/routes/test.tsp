import "@typespec/http";
import "@typespec/openapi";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("user")
  @extension("x-ogen-operation-group", "test")
  @route("/test")
  @get
  @summary("OpenAPIテスト用")
  @useAuth([])
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
  @useAuth([])
  op dummiInit(): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };
}
