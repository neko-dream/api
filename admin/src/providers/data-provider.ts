import simpleRest from "@refinedev/simple-rest";
import axios from "axios";

import { env } from "../config/env";

// 認証は HttpOnly Cookie(AdminSessionId) で行うため、全リクエストで credentials を含める。
export const httpClient = axios.create({
  baseURL: env.apiUrl,
  withCredentials: true,
});

// NOTE: /v1/manage/* のCRUD構造は未確定のため、ベースでは simple-rest を素のまま接続するだけ。
// 各リソースの一覧/詳細の実装は今後詰める。
export const dataProvider = simpleRest(`${env.apiUrl}/v1/manage`, httpClient);
