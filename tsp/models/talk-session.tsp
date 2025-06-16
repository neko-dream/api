import "@typespec/http";
import "./user.tsp";
import "./common.tsp";

using Http;

namespace kotohiro {
  model talkSession {
    /**
     * トークセッションID
     */
    id: string;

    /**
     * テーマ
     */
    theme: string;

    /**
     * 説明
     */
    description?: string | null;

    /**
     * 作成ユーザー
     */
    owner: user;

    /**
     * 作成日時
     */
    createdAt: string;

    /**
     * 終了予定日時
     */
    scheduledEndTime: string;

    /**
     * 位置情報
     */
    location?: location;

    /**
     * 市区町村
     */
    city?: string | null;

    /**
     * 都道府県
     */
    prefecture?: string | null;

    /**
     * サムネ画像
     */
    thumbnailURL?: string | null;

    /**
     * セッションの参加制限
     */
    restrictions: restriction[];

    hideReport: boolean;
  }

  model Report {
    talkSessionID: string;

    /**
     * レポート本文
     */
    content: string;
  }

  model conclusion {
    /**
     * 作成ユーザー
     */
    user: user;

    /**
     * 結論本文
     */
    content: string;
  }

  model actionItem {
    ActionItemID: string;
    Sequence: integer;
    Content: string;
    Status: string;
    CreatedAt: string;
    UpdatedAt: string;
  }

  model restriction {
    key: string;
    description: string;

    /**
     * 依存しているrestriction
     */
    dependsOn?: string[];
  }

  model policyConsentStatus {
    /**
     * 最新ポリシーのバージョン
     */
    policyVersion: string;

    /**
     * 同意した日時
     */
    consentedAt?: string | null;

    /**
     * 同意したか
     */
    consentGiven: boolean;
  }
}