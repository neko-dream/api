import simpleRest from "@refinedev/simple-rest";
import { env } from "../config/env";
import { createAdminHTTPClient } from "./http-client";
import { withUserManagement } from "./user-data-provider";

export const httpClient = createAdminHTTPClient(env.apiUrl);
export const dataProvider = withUserManagement(
  httpClient,
  simpleRest(`${env.apiUrl}/v1/manage`, httpClient),
);
