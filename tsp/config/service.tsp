import "@typespec/http";
import "@typespec/openapi";

using Http;
using OpenAPI;

namespace kotohiro {
  @doc("Cookie-based authentication using JWT tokens stored in secure HTTP-only cookies")
  model CookieAuth is ApiKeyAuth<ApiKeyLocation.cookie, "SessionId">;
}