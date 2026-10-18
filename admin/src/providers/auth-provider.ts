import type { AuthProvider } from "@refinedev/core";
import axios from "axios";

import { authErrorResult } from "./auth-error";
import { httpClient } from "./data-provider";

// 認証状態の正は HttpOnly Cookie(AdminSessionId)。ここに保存するのは UI 表示用キャッシュのみ。
const USER_KEY = "admin_user";

interface AdminTokenInfo {
  sub: string;
  email: string;
  displayName: string;
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
    } catch (error) {
      const rejected = axios.isAxiosError(error) && error.response?.status === 401;
      return {
        success: false,
        error: {
          name: "LoginError",
          message: rejected
            ? "このGoogleアカウントは管理者として登録されていません"
            : "ログインできませんでした。接続を確認して、もう一度試してください",
        },
      };
    }
  },

  logout: async () => {
    try {
      await httpClient.post("/v1/admin/logout");
    } catch {
      // Do not claim logout while HttpOnly cookies may still be valid.
      return {
        success: false,
        error: {
          name: "LogoutError",
          message: "ログアウトできませんでした。もう一度試してください",
        },
      };
    }
    localStorage.removeItem(USER_KEY);
    return { success: true, redirectTo: "/login" };
  },

  check: async () => {
    try {
      const { data } = await httpClient.get<AdminTokenInfo>("/v1/admin/token/info");
      localStorage.setItem(USER_KEY, JSON.stringify(data));
      return { authenticated: true };
    } catch (error) {
      const status = axios.isAxiosError(error) ? error.response?.status : undefined;
      if (status === 401 || status === 403) {
        localStorage.removeItem(USER_KEY);
        return { authenticated: false, redirectTo: "/login" };
      }
      return {
        authenticated: false,
        error: {
          name: "ConnectionError",
          message: "認証状態を確認できませんでした。接続を確認してください",
        },
      };
    }
  },

  onError: async (error) => {
    const result = authErrorResult(error);
    if (result.redirectTo) localStorage.removeItem(USER_KEY);
    return result;
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
