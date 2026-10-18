import assert from "node:assert/strict";
import test from "node:test";
import { authErrorResult } from "../src/providers/auth-error.ts";

// 一度の認証エラーでサーバー側のセッションまで消すと、同じブラウザの他のタブもログアウトされる。
// ログイン画面に移動するだけにし、/v1/admin/logout は呼ばせない（logout を返さない）。
test("authentication errors redirect to login without revoking the server session", () => {
  for (const status of [401, 403]) {
    const result = authErrorResult({ response: { status } });
    assert.deepEqual(result, { redirectTo: "/login" });
    assert.equal("logout" in result, false);
  }
});

test("other errors are left to the caller", () => {
  for (const error of [{ response: { status: 404 } }, { response: { status: 500 } }, new Error("offline"), undefined]) {
    assert.deepEqual(authErrorResult(error), {});
  }
});
