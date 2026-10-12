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

    isAdmin: boolean;
  }

  @doc("管理者")
  model AdminUserItem {
    id: string;
    email: string;
    displayName: string;
  }

  @doc("管理者アカウント一覧")
  model AdminUserListResponse {
    adminUsers: AdminUserItem[];
  }

  @doc("管理者作成リクエスト")
  model CreateAdminUserRequest {
    /**
     * Google アカウントのメールアドレス
     */
    email: string;

    displayName?: string;
  }
}
