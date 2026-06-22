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

  @doc("管理者アカウント")
  model AdminUserItem {
    id: string;
    email: string;
    displayName: string;

    /**
     * 管理者ロール (数値が小さいほど高権限)
     */
    role: integer;
  }

  @doc("管理者アカウント一覧")
  model AdminUserListResponse {
    adminUsers: AdminUserItem[];
  }

  @doc("管理者アカウント作成リクエスト")
  model CreateAdminUserRequest {
    /**
     * Google アカウントのメールアドレス
     */
    email: string;

    displayName?: string;

    /**
     * 管理者ロール (未指定は一般管理者)
     */
    role?: integer;
  }
}
