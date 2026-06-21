import "@typespec/http";

using Http;

namespace kotohiro {
  @doc("管理者トークンの内容")
  model AdminTokenInfo {
    /**
     * 管理者ID
     */
    sub: string;

    email: string;

    displayName: string;

    /**
     * 管理者ロール (数値が小さいほど高権限)
     */
    role: integer;

    isAdmin: boolean;
  }
}
