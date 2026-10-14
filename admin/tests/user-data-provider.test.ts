import assert from "node:assert/strict";
import test from "node:test";
import axios from "axios";
import simpleRest from "@refinedev/simple-rest";
import { withUserManagement } from "../src/providers/user-data-provider.ts";

test("user reads map pagination, search, and stable record IDs", async () => {
  const client = axios.create();
  const requests: any[] = [];
  client.defaults.adapter = async (config) => {
    requests.push(config);
    return {
      status: 200,
      statusText: "OK",
      headers: {},
      config,
      data: config.url?.endsWith("/list")
        ? { users: [{ userID: "user-a" }], totalCount: 21 }
        : { userID: "user-a" },
    };
  };
  const provider = withUserManagement(client, simpleRest("", client));
  const list = await provider.getList({
    resource: "users",
    pagination: { currentPage: 2, pageSize: 20 },
    filters: [{ field: "search", operator: "contains", value: "city" }],
  });
  assert.deepEqual(requests[0].params, {
    offset: 20,
    limit: 20,
    search: "city",
  });
  assert.equal(list.total, 21);
  assert.equal(list.data[0].id, "user-a");
  const detail = await provider.getOne({ resource: "users", id: "user-a" });
  assert.equal(requests[1].url, "/v1/manage/users/user-a");
  assert.equal(detail.data.id, "user-a");
});

test("missing users and API failures propagate without fake records", async () => {
  const client = axios.create();
  client.defaults.adapter = async () => {
    throw new Error("request failed");
  };
  const provider = withUserManagement(client, simpleRest("", client));
  await assert.rejects(provider.getOne({ resource: "users", id: "missing" }));
  await assert.rejects(provider.getList({ resource: "users" }));
});

test("other resources retain their original data provider", async () => {
  const client = axios.create();
  const urls: string[] = [];
  client.defaults.adapter = async (config) => {
    urls.push(config.url ?? "");
    return { status: 200, statusText: "OK", headers: {}, config, data: [] };
  };
  const provider = withUserManagement(client, simpleRest("/v1/manage", client));
  await provider.getList({ resource: "talksessions" });
  assert.equal(urls[0].startsWith("/v1/manage/talksessions"), true);
});

test("user creation sends only profile fields and preserves delivery status and stable ID", async () => {
  for (const invitationSent of [true, false]) {
    const client = axios.create({ withCredentials: true });
    client.defaults.adapter = async (config) => {
      assert.equal(config.url, "/v1/manage/users");
      assert.equal(config.method, "post");
      assert.equal(config.withCredentials, true);
      assert.deepEqual(JSON.parse(config.data), {
        displayName: "表示名",
        email: "user@example.com",
      });
      return {
        status: 201,
        statusText: "Created",
        headers: {},
        config,
        data: { userID: "user-a", invitationSent },
      };
    };
    const provider = withUserManagement(client, simpleRest("", client));
    const result = await provider.create({
      resource: "users",
      variables: { displayName: "表示名", email: "user@example.com" },
    });
    assert.deepEqual(result.data, {
      userID: "user-a",
      id: "user-a",
      invitationSent,
    });
  }
});

test("creation errors propagate and other resources still use their original write adapter", async () => {
  const client = axios.create();
  const provider = withUserManagement(client, simpleRest("/v1/manage", client));
  client.defaults.adapter = async () => {
    throw new Error("failed creation");
  };
  await assert.rejects(provider.create({ resource: "users", variables: {} }));
  client.defaults.adapter = async (config) => {
    assert.equal(config.url, "/v1/manage/talksessions");
    return {
      status: 201,
      statusText: "Created",
      headers: {},
      config,
      data: { id: "session-a" },
    };
  };
  assert.equal(
    (await provider.create({ resource: "talksessions", variables: {} })).data
      .id,
    "session-a",
  );
});
