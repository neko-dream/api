import "@typespec/http";

using Http;

namespace kotohiro {
  model ValidationErrorItem {
    /**
     * バリデーションエラーのフィールド
     */
    field: string;

    /**
     * バリデーションエラーメッセージ
     */
    message: string;
  }

  model Error {
    code: string;
    message: string;
  }

  model OffsetPagination {
    totalCount: integer;
    offset: integer;
    limit: integer;
  }

  model Success {
    message: string;
  }

  model Location {
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