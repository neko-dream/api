// API のベースURL。
// dev では vite proxy 経由で同一オリジン(/v1, /auth)に流すため空文字。
// 本番等では VITE_API_URL を設定する。
export const env = {
  apiUrl: import.meta.env.VITE_API_URL ?? "",
};
