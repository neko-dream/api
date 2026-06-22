import type { AuthProvider } from "@refinedev/core";

import { httpClient } from "./data-provider";

// 認証状態の正は HttpOnly Cookie(AdminSessionId)。ここに保存するのは UI 表示用キャッシュのみ。
const USER_KEY = "admin_user";

interface AdminTokenInfo {
  sub: string;
  email: string;
  displayName: string;
  role: number;
  isAdmin: boolean;
}

export const authProvider: AuthProvider = {
  login: async (params) => {
    // Google Identity Services から受け取った credential(ID トークン)。
    const { credential } = params as { credential?: string };
    if (!credential) {
      return {
        success: false,
        error: { name: "LoginError", message: "Google 認証に失敗しました" },
      };
    }
    try {
      const res = await httpClient.post<AdminTokenInfo>("/v1/admin/login", {
        idToken: credential,
      });
      localStorage.setItem(USER_KEY, JSON.stringify(res.data));
      return { success: true, redirectTo: "/" };
    } catch {
      return {
        success: false,
        error: {
          name: "LoginError",
          message: "このGoogleアカウントは管理者として登録されていません",
        },
      };
    }
  },

  logout: async () => {
    try {
      await httpClient.post("/v1/admin/logout");
    } catch {
      // Cookie はサーバー側で失効。失敗してもローカルは掃除してログイン画面へ。
    }
    localStorage.removeItem(USER_KEY);
    return { success: true, redirectTo: "/login" };
  },

  check: async () => {
    try {
      await httpClient.get("/v1/admin/token/info");
      return { authenticated: true };
    } catch {
      localStorage.removeItem(USER_KEY);
      return { authenticated: false, redirectTo: "/login" };
    }
  },

  onError: async (error) => {
    if (error?.response?.status === 403 || error?.response?.status === 401) {
      return { logout: true, redirectTo: "/login" };
    }
    return {};
  },

  getIdentity: async () => {
    const raw = localStorage.getItem(USER_KEY);
    if (!raw) return null;
    try {
      const user = JSON.parse(raw) as AdminTokenInfo;
      return { id: user.sub, name: user.displayName, email: user.email };
    } catch {
      return null;
    }
  },
};
