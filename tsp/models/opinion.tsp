import "@typespec/http";
import "./user.tsp";

using Http;

namespace kotohiro {
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
    voteType?: "disagree" | "agree" | "pass" | null;

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

    status: "unsolved" | "deleted" | "hold";
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

