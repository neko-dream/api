import assert from "node:assert/strict";
import test from "node:test";
import { AxiosError, type InternalAxiosRequestConfig } from "axios";
import { createAdminHTTPClient } from "../src/providers/http-client.ts";

function response(config: InternalAxiosRequestConfig, status = 200) {
  return { config, status, statusText: String(status), headers: {}, data: {} };
}
function reject(config: InternalAxiosRequestConfig, status: number) {
  return Promise.reject(
    new AxiosError("failed", undefined, config, undefined, response(config, status)),
  );
}

test("credentials are included and parallel requests share one refresh", async () => {
  const client = createAdminHTTPClient("");
  let refreshed = false;
  let refreshes = 0;
  client.defaults.adapter = async (config) => {
    assert.equal(config.withCredentials, true);
    if (config.url === "/v1/admin/refresh") {
      refreshes++;
      await new Promise((resolve) => setTimeout(resolve, 10));
      refreshed = true;
      return response(config);
    }
    return refreshed ? response(config) : reject(config, 401);
  };
  await Promise.all([client.get("/v1/admin/token/info"), client.get("/v1/admin/admin-users")]);
  assert.equal(refreshes, 1);
});

test("protected requests are retried at most once", async () => {
  const client = createAdminHTTPClient("");
  let refreshes = 0;
  let attempts = 0;
  client.defaults.adapter = async (config) => {
    if (config.url === "/v1/admin/refresh") {
      refreshes++;
      return response(config);
    }
    attempts++;
    return reject(config, 401);
  };
  await assert.rejects(client.get("/v1/admin/admin-users"));
  assert.equal(refreshes, 1);
  assert.equal(attempts, 2);
});

test("failed refresh does not loop", async () => {
  const client = createAdminHTTPClient("");
  const paths: string[] = [];
  client.defaults.adapter = async (config) => {
    paths.push(config.url!);
    return reject(config, 401);
  };
  await assert.rejects(client.get("/v1/admin/admin-users"));
  assert.deepEqual(paths, ["/v1/admin/admin-users", "/v1/admin/refresh"]);
});

test("refresh server failure remains a server failure, not authentication rejection", async () => {
  const client = createAdminHTTPClient("");
  client.defaults.adapter = async (config) => reject(config, config.url === "/v1/admin/refresh" ? 500 : 401);
  await assert.rejects(client.get("/v1/admin/admin-users"), (error: any) => error.response?.status === 500);
});

test("authentication endpoints and non-401 errors never trigger refresh", async () => {
  for (const [path, status] of [
    ["/v1/admin/login", 401],
    ["/v1/admin/logout", 401],
    ["/v1/admin/refresh", 401],
    ["/v1/admin/admin-users", 403],
    ["/v1/admin/admin-users", 500],
  ] as const) {
    const client = createAdminHTTPClient("");
    let attempts = 0;
    client.defaults.adapter = async (config) => {
      attempts++;
      return reject(config, status);
    };
    await assert.rejects(client.get(path));
    assert.equal(attempts, 1);
  }
});
