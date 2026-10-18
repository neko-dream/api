// 管理APIの認証エラーへの対応を決める（authProvider.onError から使う）。
// refine は logout を返すと authProvider.logout（/v1/admin/logout）を呼び、サーバー側のセッションを消す。
// リフレッシュの一時的な失敗でも同じブラウザの他のタブまでログアウトされるので、ログイン画面に移動するだけにする。
export function authErrorResult(error: unknown): { redirectTo?: string } {
  const status = (error as { response?: { status?: number } } | undefined)?.response?.status;
  if (status === 401 || status === 403) {
    return { redirectTo: "/login" };
  }
  return {};
}
