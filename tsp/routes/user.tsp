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
  op get_user_info(): Body<{
    user: user;
    demographics: userDemographics;
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
  op registerUser(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      displayName: HttpPart<string>;
      displayID: HttpPart<string>;
      icon?: HttpPart<bytes>;
      dateOfBirth?: HttpPart<string>;
      gender?: HttpPart<"男性" | "女性" | "その他" | "回答しない" | null>;
      prefecture?: HttpPart<string>;
      city?: HttpPart<string | null>;
      email?: HttpPart<string | null>;
    },
  ): Body<{
    @minLength(3)
    @maxLength(25)
    @pattern("^[A-Za-z0-9]$")
    displayID: string;

    displayName: string;
    iconURL?: string | null;
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
  @extension("x-ogen-operation-group", "user")
  @route("/user")
  @put
  @summary("ユーザー情報の変更")
  op editUserProfile(
    @header contentType: "multipart/form-data",
    @multipartBody body: {
      displayName?: HttpPart<string | null>;
      icon?: HttpPart<bytes>;
      deleteIcon?: HttpPart<boolean | null>;
      dateOfBirth?: HttpPart<string>;
      gender?: HttpPart<"男性" | "女性" | "その他" | "回答しない" | null>;
      city?: HttpPart<string | null>;
      prefecture?: HttpPart<string | null>;
      email?: HttpPart<string | null>;
    },
  ): Body<{
    @minLength(3)
    @maxLength(25)
    @pattern("^[A-Za-z0-9]$")
    displayID: string;

    displayName: string;
    iconURL?: string | null;
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
}