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

  /**
   * ゲストユーザーとしてログインします。
   * ゲストユーザーは登録制限のないセッションでのみ意見投稿・投票が可能です。
   * Google/LINEログインすると正規ユーザーに昇格します。
   */
  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/guest/login")
  @post
  @summary("ゲストログイン")
  @useAuth([])
  op guestLogin(): {
    @statusCode statusCode: 200;
    @header("Set-Cookie") SetCookie: string[];
    @body body: {};
  } | {
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
  @summary("パスワードを確認し、ログインごとのメール本人確認を開始する")
  @useAuth([])
  op passwordLogin(
    @multipartBody body: {
      idOrEmail: HttpPart<string>;
      password: HttpPart<string>;
    },
  ): {
    @statusCode statusCode: 202;
    @header("Set-Cookie") SetCookie: string[];
    @body body: { verificationRequired: true; expiresAt: utcDateTime; };
  } | {
    @statusCode statusCode: 401;
    @body body: { code: string; message: string; };
  } | {
    @statusCode statusCode: 429;
    @body body: { code: string; message: string; };
  } | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  /** GETではリンクを消費しない。ログインを開始したブラウザで確認ボタンを押す。 */
  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/verify")
  @get
  @useAuth([])
  @summary("メール本人確認ページ（リンクはURLフラグメントに格納）")
  op passwordEmailConfirmation(): {
    @statusCode statusCode: 200;
    @header contentType: "text/html; charset=utf-8";
    @body body: string;
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/verify")
  @post
  @useAuth([])
  @summary("メール本人確認を完了し、初めてログインCookieを発行する")
  op verifyPasswordEmail(@body body: { token: string; }): {
    @statusCode statusCode: 200;
    @header("Set-Cookie") SetCookie: string[];
    @body body: { redirectURL: string; };
  } | {
    @statusCode statusCode: 400;
    @body body: { code: string; message: string; };
  } | {
    @statusCode statusCode: 500;
    @body body: { code: string; message: string; };
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/verify-pin")
  @post
  @useAuth([])
  @summary("メールの6桁PINで本人確認を完了する")
  op verifyPasswordEmailPIN(@body body: { pin: string; }): {
    @statusCode statusCode: 200;
    @header("Set-Cookie") SetCookie: string[];
    @body body: { redirectURL: string; };
  } | {
    @statusCode statusCode: 400;
    @body body: { code: string; message: string; };
  } | {
    @statusCode statusCode: 429;
    @body body: { code: string; message: string; };
  } | {
    @statusCode statusCode: 500;
    @body body: { code: string; message: string; };
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

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/setup")
  @post
  @useAuth([])
  @summary("メールから初回パスワードを設定する（ログインCookieは発行しない）")
  op setupPassword(@body body: {
    @minLength(44)
    @maxLength(44)
    token: string;
    password: string;
  }): {
    @statusCode statusCode: 204;
  };

  /** 入力条件はサーバーの共通パスワード検証から取得する。秘密値は返さない。 */
  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/password/policy")
  @get
  @useAuth([])
  @summary("新規設定・変更用のパスワード入力条件")
  op getPasswordPolicy(): {
    @statusCode statusCode: 200;
    @body body: {
      /** 最小文字数（Unicodeコードポイント数） */
      minLength: int32;

      /** 最大バイト数（UTF-8） */
      maxBytes: int32;

      /** 改行・タブなどの制御文字を禁止する */
      forbidControlCharacters: boolean;

      /** 先頭・末尾の空白（全角空白を含む）を禁止する */
      forbidSurroundingWhitespace: boolean;
    };
  };

  @tag("auth")
  @extension("x-ogen-operation-group", "Auth")
  @route("/auth/reactivate")
  @post
  @summary("退会ユーザーの復活")
  op reactivateUser(): {
    @statusCode statusCode: 200;
    @body body: {
      message: string;
      user: User;
    };
  } | {
    @statusCode statusCode: 400;
    @body body: {
      code: string;
      message: string;
    };
  } | {
    @statusCode statusCode: 403;
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
