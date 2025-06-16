import "@typespec/http";

using Http;

namespace kotohiro {
  model TokenClaim {
    /**
     * Audience
     */
    aud: string;

    /**
     * 有効期限
     */
    exp: string;

    /**
     * 発行日時
     */
    iat: string;

    /**
     * 発行者
     */
    iss: string;

    /**
     * ユーザID
     */
    `sub`: string;

    /**
     * JWT ID
     */
    jti: string;

    /**
     * ユーザーID
     */
    displayID?: string;

    /**
     * ユーザー名
     */
    displayName?: string;

    /**
     * アイコンURL
     */
    iconURL?: string;

    /**
     * ユーザ登録済みか
     */
    isRegistered: boolean;

    isEmailVerified: boolean;

    /**
     * アカウントの種類。組織がなければ空
     */
    orgType?: integer | null;

    /**
     * パスワードの更新が必要かどうか
     */
    requiredPasswordChange: boolean;

        organizationRole?: string | null;
    organizationCode?: string | null;
    organizationID?: string | null;
  }
}
