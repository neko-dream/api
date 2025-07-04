import "@typespec/http";
import "@typespec/openapi";
import "../config/service.tsp";
import "../models/auth.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/change")
  @put
  @summary("パスワード変更")
  op changePassword(
    /**
     * 古いパスワード
     */
    @query(#{ explode: true }) oldPassword: string,

    /**
     * 新たなパスワード
     */
    @query(#{ explode: true }) newPassword: string,
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/token/info")
  @get
  @summary("JWTの内容を返してくれる")
  op getTokenInfo(): Body<TokenClaim> | {
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

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/{provider}/callback")
  @get
  @summary("Auth Callback")
  @useAuth([])
  op handleAuthCallback(
    @path provider: string,
    @query(#{ explode: true }) code: string,

    /**
     * OAuth State from Query
     */
    @query(#{ explode: true }) state: string,
  ): {
    @statusCode statusCode: 302;
    @header("Location") Location: string;

    /**
     * SessionID
     */
    @header("Set-Cookie") SetCookie: string[];

    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/dev/login")
  @get
  @summary("開発用登録/ログイン")
  @useAuth([])
  op devAuthorize(
    @query(#{ explode: true }) redirect_url: string,

    /**
     * devのみで使用するsubjectです。ここで指定した値はログインした後も確認できないため覚えておいてください。同じ値を指定すると同じアカウントにログインできます。
     */
    @query(#{ explode: true }) id: string,

    /**
     * 組織コード（組織ログインの場合）
     */
    @query(#{ explode: true }) organizationCode?: string,
  ): {
    @statusCode statusCode: 302;
    @header("Location") Location?: string;
    @header("Set-Cookie") SetCookie?: string[];
    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/{provider}/login")
  @get
  @summary("ログイン")
  @useAuth([])
  op authorize(
    @path provider: "google" | "line",

    /**
     * ログイン後にリダイレクトするURL
     */
    @query(#{ explode: true }) redirect_url: string,

    /**
     * 組織コード（組織ログインの場合）
     */
    @query(#{ explode: true }) organization_code?: string,

    /**
     * 登録していなかった場合に飛ばすURL
     */
    @query(#{ explode: true }) registration_url?: string,
  ): {
    @statusCode statusCode: 302;

    /**
     * IDPのログインページ
     */
    @header("Location") Location: string;

    /**
     * OAuth2.0 State
     */
    @header("Set-Cookie") SetCookie?: string[];

    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /**
   * そのアカウントには再度ログインできなくなります。ログインしたければ言ってね！
   */
  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/dev/detach")
  @delete
  @summary("**開発用** 現在のアカウントを切り離す")
  op authAccountDetach(): {
    @header("Set-Cookie") SetCookie?: string[];
    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/revoke")
  @post
  @summary("トークンを失効（ログアウト）")
  op revokeToken(): {
    @statusCode statusCode: 204;

    /**
     * Cookie削除用のSet-Cookie
     */
    @header("Set-Cookie") SetCookie: string[];

    @body body: {};
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/login")
  @post
  @summary("パスワードによるログイン")
  @useAuth([])
  op passwordLogin(
    @multipartBody body: {
      idOrEmail: HttpPart<string>;
      password: HttpPart<string>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/register")
  @post
  @summary("パスワードによる登録（devのみ）")
  @useAuth([])
  op passwordRegister(
    @multipartBody body: {
      password: HttpPart<string>;
      email: HttpPart<string>;
    },
  ): Body<{}> | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };
}
