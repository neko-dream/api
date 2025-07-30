import "@typespec/http";
import "@typespec/openapi";
import "../models/talk-session.tsp";

using Http;
using OpenAPI;

namespace kotohiro {
  @tag("policy")
  @extension("x-ogen-operation-group", "Policy")
  @route("/policy/consent")
  @get
  @summary("最新のポリシーに同意したかを取得")
  @useAuth([])
  op getPolicyConsentStatus(): PolicyConsentStatus | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };

  @tag("policy")
  @extension("x-ogen-operation-group", "Policy")
  @route("/policy/consent")
  @post
  @summary("最新のポリシーに同意する")
  op policyConsent(
    @multipartBody body: {
      policyVersion: HttpPart<string>;
    },
  ): PolicyConsentStatus | {
    @statusCode statusCode: 400;
    @body body: {};
  } | {
    @statusCode statusCode: 500;
    @body body: {};
  };
}
