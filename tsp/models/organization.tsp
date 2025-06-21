import "@typespec/http";

using Http;

namespace kotohiro {
  model Organization {
    /**
     * 組織ID
     */
    ID: string;

    /**
     * 組織名
     */
    name: string;

    /**
     * 組織コード
     */
    code: string;

    /**
     * 組織のタイプ
     */
    type: integer;

    roleName: string;

    /**
     * ロール
     */
    role: integer;
  }

  /**
   * 組織エイリアス
   */
  model OrganizationAlias {
    aliasID: string;
    aliasName: string;
    createdAt: string;
  }
}