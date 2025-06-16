import "@typespec/http";

using Http;

namespace kotohiro {
  model validationErrorItem {
    /**
     * バリデーションエラーのフィールド
     */
    field: string;

    /**
     * バリデーションエラーメッセージ
     */
    message: string;
  }

  model error {
    code: string;
    message: string;
  }

  model offsetPagination {
    totalCount: integer;
    offset: integer;
    limit: integer;
  }

  model success {
    message: string;
  }

  model location {
    /**
     * 緯度
     */
    latitude?: numeric;

    /**
     * 経度
     */
    longitude?: numeric;
  }
}