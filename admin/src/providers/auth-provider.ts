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
    const { email, password } = params as { email: string; password: string };
    try {
      const form = new FormData();
      form.append("email", email);
      form.append("password", password);
      const res = await httpClient.post<AdminTokenInfo>("/v1/admin/login", form);
      localStorage.setItem(USER_KEY, JSON.stringify(res.data));
      return { success: true, redirectTo: "/" };
    } catch {
      return {
        success: false,
        error: {
          name: "LoginError",
          message: "メールアドレスまたはパスワードが正しくありません",
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
