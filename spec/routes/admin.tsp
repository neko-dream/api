import "@typespec/http";
import "@typespec/openapi";
import "../config/service.tsp";
import "../models/admin.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @route("/v1/admin")
  namespace KotohiroAdminService {
    @route("/login")
    @operationId("adminLogin")
    @extension("x-ogen-operation-group", "Admin")
    @tag("admin")
    @post
    @summary("管理者ログイン(Google認証)")
    @useAuth([])
    op adminLogin(
      @body body: {
        /**
         * Google Identity Services で取得した ID トークン(credential)
         */
        idToken: string;
      },
    ): {
      @statusCode statusCode: 200;

      /**
       * AdminSessionId Cookie
       */
      @header("Set-Cookie") SetCookie: string[];

      @body body: AdminTokenInfo;
    } | {
      @statusCode statusCode: 401;
      @body body: {};
    } | {
      @statusCode statusCode: 500;
      @body body: {};
    };

    @route("/token/info")
    @operationId("getAdminTokenInfo")
    @extension("x-ogen-operation-group", "Admin")
    @tag("admin")
    @get
    @summary("管理者トークンの内容を返す")
    @useAuth(AdminCookieAuth)
    op getAdminTokenInfo(): Body<AdminTokenInfo> | {
      @statusCode statusCode: 403;
      @body body: {};
    };

    @route("/refresh")
    @operationId("adminRefresh")
    @extension("x-ogen-operation-group", "Admin")
    @tag("admin")
    @post
    @summary("リフレッシュトークンでアクセストークンを再発行")
    @useAuth([])
    op adminRefresh(): {
      @statusCode statusCode: 200;

      /**
       * 新しい AdminSessionId / AdminRefreshToken Cookie
       */
      @header("Set-Cookie") SetCookie: string[];

      @body body: AdminTokenInfo;
    } | {
      @statusCode statusCode: 401;
      @body body: {};
    };

    @route("/logout")
    @operationId("adminLogout")
    @extension("x-ogen-operation-group", "Admin")
    @tag("admin")
    @post
    @summary("管理者ログアウト")
    @useAuth([])
    op adminLogout(): {
      @statusCode statusCode: 204;

      /**
       * Cookie削除用のSet-Cookie
       */
      @header("Set-Cookie") SetCookie: string[];
    };
  }
}
