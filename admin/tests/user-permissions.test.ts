import assert from "node:assert/strict";
import test from "node:test";
import axios from "axios";
import simpleRest from "@refinedev/simple-rest";
import { canGrantSessionCreation, sessionCreationPermissionLabel } from "../src/pages/users/types.ts";

test("permission UI distinguishes stored grants from currently inactive users", () => {
  for (const [providers, withdrawalDate, eligible] of [
    [["GOOGLE"], undefined, true],
    [["DEV"], undefined, true],
    [["GUEST", "GOOGLE"], undefined, true],
    [["GUEST"], undefined, false],
    [[], undefined, false],
    [["GOOGLE"], "2026-10-12T00:00:00Z", false],
  ] as const) {
    const user = { providers: [...providers], withdrawalDate, permissions: { canCreateTalkSession: true } };
    assert.equal(canGrantSessionCreation(user), eligible);
    assert.equal(sessionCreationPermissionLabel(user), eligible ? "許可あり" : "許可あり（現在は無効）");
    assert.equal(sessionCreationPermissionLabel({ ...user, permissions: { canCreateTalkSession: false } }), "許可なし");
  }
});

test("the existing custom data provider sends the permission mutation without a new provider", async () => {
  const client = axios.create({ withCredentials: true });
  const requests: any[] = [];
  client.defaults.adapter = async (config) => {
    requests.push(config);
    return { status: 200, statusText: "OK", headers: {}, config, data: { canCreateTalkSession: false } };
  };
  const provider = simpleRest("https://api.example.test/v1/manage", client);
  assert.ok(provider.custom);
  for (const allowed of [true, false]) {
    await provider.custom({ url: "https://api.example.test/v1/manage/users/user-a/permissions", method: "put", payload: { canCreateTalkSession: allowed } });
    const request = requests.at(-1);
    assert.equal(request.method, "put");
    assert.equal(request.withCredentials, true);
    assert.deepEqual(JSON.parse(request.data), { canCreateTalkSession: allowed });
  }
});
