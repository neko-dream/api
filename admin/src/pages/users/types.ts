export interface ManagedUser {
  id: string;
  userID: string;
  displayID: string;
  displayName: string;
  email?: string;
  emailVerified: boolean;
  providers: string[];
  createdAt: string;
  updatedAt: string;
  withdrawalDate?: string;
}

export function formatUserDate(value?: string): string {
  return value ? new Date(value).toLocaleString("ja-JP") : "—";
}

export function providerLabel(provider: string): string {
  return (
    (
      {
        GOOGLE: "Google",
        LINE: "LINE",
        PASSWORD: "メール・パスワード",
        GUEST: "ゲスト",
        DEV: "開発用",
      } as Record<string, string>
    )[provider] ?? provider
  );
}
