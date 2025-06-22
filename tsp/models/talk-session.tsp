import "@typespec/http";
import "./user.tsp";
import "./common.tsp";

using Http;

namespace kotohiro {
  model TalkSession {
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
    owner: User;

    /**
     * 作成組織名
     */
    organizationAlias?: OrganizationAlias | null;

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
    location?: Location;

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
    restrictions: Restriction[];

    /**
     * レポートを隠すかどうか
     */
    hideReport: boolean;
  }

  model Report {
    talkSessionID: string;

    /**
     * レポート本文
     */
    content: string;
  }

  model Conclusion {
    /**
     * 作成ユーザー
     */
    user: User;

    /**
     * 結論本文
     */
    content: string;
  }

  model ActionItem {
    actionItemID: string;
    sequence: integer;
    content: string;
    status: string;
    createdAt: string;
    updatedAt: string;
  }

  model Restriction {
    key: string;
    description: string;

    /**
     * 依存しているrestriction
     */
    dependsOn?: string[];
  }

  model PolicyConsentStatus {
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

