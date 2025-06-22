import "@typespec/http";
import "@typespec/openapi";

using Http;
using OpenAPI;

namespace kotohiro {
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
}