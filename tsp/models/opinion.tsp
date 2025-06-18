import "@typespec/http";
import "./user.tsp";

using Http;

namespace kotohiro {
  /**
   * 投票タイプ
   */
  enum VoteType {
    Agree: "agree",
    Disagree: "disagree",
    Pass: "pass",
  }

  /**
   * 通報ステータス
   */
  enum ReportStatus {
    Unsolved: "unsolved",
    Deleted: "deleted",
    Hold: "hold",
  }

  /**
   * 通報解決アクション
   */
  enum ReportAction {
    Deleted: "deleted",
    Hold: "hold",
  }
  model Opinion {
    /**
     * 意見ID
     */
    id: string;

    title?: string;

    /**
     * 意見のテキスト
     */
    @minLength(1) @maxLength(140) content: string;

    /**
     * 親の意見ID。ルートならば無し
     */
    parentID?: string;

    /**
     * 意見投稿主の意見。ルート意見の場合はここには何も入らない
     */
    voteType?: VoteType | null;

    /**
     * 画像が返る場合もある
     */
    pictureURL?: string | null;

    /**
     * 参考文献URL
     */
    referenceURL?: string;

    postedAt: string;
    isDeleted: boolean;
  }

  model ReportReason {
    /**
     * 1
     */
    reasonID: integer;

    /**
     * 不適切な内容
     */
    reason: string;
  }

  model OpinionGroupRatio {
    agreeCount: integer;
    disagreeCount: integer;
    passCount: integer;
    groupID: integer;
    groupName: string;
  }

  model ReportDetail {
    opinion: Opinion;

    /**
     * 作成ユーザー
     */
    user: User;

    status: ReportStatus;
    reasons: {
      reason: string;
      content?: string | null;
    }[];

    /**
     * この意見が通報を受けた回数
     */
    reportCount: integer;
  }
}

