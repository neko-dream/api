import simpleRest from "@refinedev/simple-rest";
import { env } from "../config/env";
import { createAdminHTTPClient } from "./http-client";

export const httpClient = createAdminHTTPClient(env.apiUrl);
export const dataProvider = simpleRest(`${env.apiUrl}/v1/manage`, httpClient);
