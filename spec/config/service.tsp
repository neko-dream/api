import "@typespec/http";

using Http;
using OpenAPI;

namespace kotohiro {
  @doc("Cookie-based authentication using JWT tokens stored in secure HTTP-only cookies")
  model CookieAuth is ApiKeyAuth<ApiKeyLocation.cookie, "SessionId">;

  @doc("Optional cookie-based authentication - will populate session context if authenticated but won't require it")
  model OptionalCookieAuth is ApiKeyAuth<ApiKeyLocation.cookie, "SessionId">;

  @doc("管理画面用のCookie認証。user の SessionId とは独立した AdminSessionId を使用する")
  model AdminCookieAuth is ApiKeyAuth<ApiKeyLocation.cookie, "AdminSessionId">;
}
