// セッション詳細画面（意見・AI分析タブ）の表示用の変換。

export interface ManagedOpinion {
  opinionID: string;
  content: string;
  pictureURL?: string;
  referenceURL?: string;
  parentOpinionID?: string;
  author: { userID: string; displayID: string; displayName: string; iconURL: string };
  createdAt: string;
  agreeCount: number;
  disagreeCount: number;
  passCount: number;
  reportCount: number;
  unsolvedReportCount: number;
  deleted: boolean;
}

export interface OpinionStatusTag {
  label: string;
  color?: string;
}

export function opinionStatusTags(opinion: ManagedOpinion): OpinionStatusTag[] {
  const tags: OpinionStatusTag[] = [];
  if (opinion.parentOpinionID) tags.push({ label: "返信", color: "blue" });
  if (opinion.deleted) tags.push({ label: "削除済み", color: "red" });
  if (opinion.reportCount > 0) {
    tags.push({
      label:
        opinion.unsolvedReportCount > 0
          ? `通報 ${opinion.reportCount}件（未対応 ${opinion.unsolvedReportCount}件）`
          : `通報 ${opinion.reportCount}件`,
      color: opinion.unsolvedReportCount > 0 ? "orange" : undefined,
    });
  }
  return tags;
}

/** 賛成率は保留を除いた賛成・反対のうちの割合（%）。賛否の投票がなければ null */
export function voteSummary(opinion: ManagedOpinion): { total: number; agreeRate: number | null } {
  const decided = opinion.agreeCount + opinion.disagreeCount;
  return {
    total: decided + opinion.passCount,
    agreeRate: decided === 0 ? null : Math.round((opinion.agreeCount / decided) * 100),
  };
}

// GET /talksessions/{id}/analysis のうち、管理画面で使う部分
export interface TalkSessionAnalysis {
  positions: { groupID: number; groupName: string }[];
  groupOpinions: {
    groupID: number;
    groupName: string;
    opinions: {
      opinion: { id: string; content: string };
      user: { displayID: string; displayName: string };
      agreeCount: number;
      disagreeCount: number;
      passCount: number;
    }[];
  }[];
}

export interface GroupSummary {
  groupID: number;
  groupName: string;
  memberCount: number;
  opinions: TalkSessionAnalysis["groupOpinions"][number]["opinions"];
}

export function groupSummaries(analysis: TalkSessionAnalysis): GroupSummary[] {
  const members = new Map<number, number>();
  for (const p of analysis.positions) members.set(p.groupID, (members.get(p.groupID) ?? 0) + 1);
  return analysis.groupOpinions
    .map((g) => ({
      groupID: g.groupID,
      groupName: g.groupName,
      memberCount: members.get(g.groupID) ?? 0,
      opinions: g.opinions,
    }))
    .sort((a, b) => a.groupID - b.groupID);
}

/** POST /v1/manage/talksessions/{id}/analysis/report は application/x-www-form-urlencoded を受け取る */
export function reportVisibilityBody(hidden: boolean): string {
  return new URLSearchParams({ hidden: String(hidden) }).toString();
}

export type RegenerateType = "report" | "group" | "image";

/** POST /v1/manage/talksessions/{id}/analysis/regenerate も application/x-www-form-urlencoded を受け取る */
export function regenerateBody(type: RegenerateType): string {
  return new URLSearchParams({ type }).toString();
}
