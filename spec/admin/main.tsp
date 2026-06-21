import "@typespec/http";
import "@typespec/openapi";
import "@typespec/openapi3";

// 管理API は app spec のサブセット。管理画面(Refine)向けに独立した OpenAPI を出力する。
import "../config/service.tsp";
import "../models/manage.tsp";
import "../routes/manage.tsp";

using Http;
using OpenAPI;

@service(#{ title: "Kotohiro Admin API" })
@info(#{ version: "1.0.0" })
@useAuth(CookieAuth)
namespace kotohiro {

}
