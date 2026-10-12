import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";

// 認証は HttpOnly Cookie(AdminSessionId/AdminRefreshToken) で行うため、全リクエストで credentials を含める。
export function createAdminHTTPClient(apiUrl: string) {
  const httpClient = axios.create({
    baseURL: apiUrl,
    withCredentials: true,
  });

  // access(15分の短命JWT)が切れると 401 が返る。1 度だけ /v1/admin/refresh で再発行してリトライする。
  // refresh 呼び出しは並行リクエストでも 1 本に集約する。
  type RetriableConfig = InternalAxiosRequestConfig & { _retried?: boolean };

  const NO_RETRY_PATHS = new Set(["/v1/admin/refresh", "/v1/admin/login", "/v1/admin/logout"]);

  let refreshInFlight: Promise<void> | null = null;

  const refreshAccessToken = (): Promise<void> => {
    if (!refreshInFlight) {
      refreshInFlight = httpClient
        .post("/v1/admin/refresh")
        .then(() => undefined)
        .finally(() => {
          refreshInFlight = null;
        });
    }
    return refreshInFlight;
  };

  httpClient.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      const status = error.response?.status;
      const original = error.config as RetriableConfig | undefined;

      if (status !== 401 || !original || original._retried) {
        return Promise.reject(error);
      }
      // 認証系エンドポイント自身の 401 は再試行しない(ログイン画面に戻す)。
      if (NO_RETRY_PATHS.has(original.url ?? "")) {
        return Promise.reject(error);
      }

      try {
        await refreshAccessToken();
      } catch (refreshError) {
        return Promise.reject(refreshError);
      }

      original._retried = true;
      return httpClient.request(original);
    },
  );

  // NOTE: /v1/manage/* のCRUD構造は未確定のため、ベースでは simple-rest を素のまま接続するだけ。
  // 各リソースの一覧/詳細の実装は今後詰める。
  return httpClient;
}
