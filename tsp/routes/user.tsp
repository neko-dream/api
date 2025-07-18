import "@typespec/http";
import "@typespec/openapi";
import "../models/user.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("user")
  @extension("x-ogen-operation-group", "User")
  @route("/user")
  @get
  @summary("ユーザー情報の取得")
  op getUserInfo(): Body<{
    user: User;
    demographics: UserDemographics;
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
  op establishUser(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      displayName: HttpPart<string>;
      displayID: HttpPart<string>;
      icon?: HttpPart<bytes>;
      dateOfBirth?: HttpPart<numeric>;
      gender?: HttpPart<string>;
      prefecture?: HttpPart<string>;
      city?: HttpPart<string>;
      email?: HttpPart<string>;
    },
  ): User | {
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
  op updateUserProfile(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      displayName?: HttpPart<string>;
      icon?: HttpPart<bytes>;
      deleteIcon?: HttpPart<boolean | null>;
      dateOfBirth?: HttpPart<numeric>;
      gender?: HttpPart<string>;
      city?: HttpPart<string>;
      prefecture?: HttpPart<string>;
      email?: HttpPart<string>;
    },
  ): User | {
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
