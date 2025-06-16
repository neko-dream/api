import "@typespec/http";

using Http;

namespace kotohiro {
  model organization {
    /**
     * 組織ID
     */
    ID: string;

    /**
     * 組織名
     */
    Name: string;

    /**
     * 組織コード
     */
    Code: string;

    /**
     * 組織のタイプ
     */
    Type: integer;

    RoleName: string;

    /**
     * ロール
     */
    Role: integer;
  }
}