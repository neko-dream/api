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
  permissions: { canCreateTalkSession: boolean };
}

export function canGrantSessionCreation(user: Pick<ManagedUser, "withdrawalDate" | "providers">): boolean {
  return !user.withdrawalDate && user.providers.some((provider) => provider !== "GUEST");
}

export function sessionCreationPermissionLabel(user: Pick<ManagedUser, "withdrawalDate" | "providers" | "permissions">): string {
  if (!user.permissions.canCreateTalkSession) return "許可なし";
  return canGrantSessionCreation(user) ? "許可あり" : "許可あり（現在は無効）";
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
