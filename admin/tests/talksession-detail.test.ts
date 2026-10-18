import assert from "node:assert/strict";
import test from "node:test";
import {
  groupSummaries,
  opinionStatusTags,
  regenerateBody,
  reportVisibilityBody,
  voteSummary,
  type ManagedOpinion,
  type TalkSessionAnalysis,
} from "../src/pages/talksessions/detail-view.ts";

const opinion: ManagedOpinion = {
  opinionID: "o-1",
  content: "意見",
  author: { userID: "u-1", displayID: "user", displayName: "ユーザー", iconURL: "" },
  createdAt: "2026-10-01T00:00:00Z",
  agreeCount: 3,
  disagreeCount: 1,
  passCount: 0,
  reportCount: 0,
  unsolvedReportCount: 0,
  hidden: false,
};

test("status tags show reply, deletion and reports only when relevant", () => {
  assert.deepEqual(opinionStatusTags(opinion), []);
  const labels = opinionStatusTags({
    ...opinion,
    parentOpinionID: "o-0",
    hidden: true,
    reportCount: 3,
    unsolvedReportCount: 2,
  }).map((t) => t.label);
  assert.deepEqual(labels, ["返信", "非表示（削除済み）", "通報 3件（未対応 2件）"]);
  assert.deepEqual(
    opinionStatusTags({ ...opinion, reportCount: 1 }).map((t) => t.label),
    ["通報 1件"],
  );
});

test("vote summary totals votes and rounds the agree ratio", () => {
  assert.deepEqual(voteSummary(opinion), { total: 4, agreeRate: 75 });
  assert.deepEqual(voteSummary({ ...opinion, agreeCount: 0, disagreeCount: 0 }), {
    total: 0,
    agreeRate: null,
  });
  assert.equal(voteSummary({ ...opinion, agreeCount: 1, disagreeCount: 2 }).agreeRate, 33);
});

test("groups are ordered by id with member counts from positions", () => {
  const analysis: TalkSessionAnalysis = {
    positions: [
      { groupID: 1, groupName: "B" },
      { groupID: 0, groupName: "A" },
      { groupID: 1, groupName: "B" },
    ],
    groupOpinions: [
      { groupID: 1, groupName: "B", opinions: [] },
      {
        groupID: 0,
        groupName: "A",
        opinions: [
          {
            opinion: { id: "o-1", content: "代表的な意見" },
            user: { displayID: "user", displayName: "ユーザー" },
            agreeCount: 2,
            disagreeCount: 0,
            passCount: 1,
          },
        ],
      },
    ],
  };
  const groups = groupSummaries(analysis);
  assert.deepEqual(
    groups.map((g) => [g.groupID, g.groupName, g.memberCount, g.opinions.length]),
    [
      [0, "A", 1, 1],
      [1, "B", 2, 0],
    ],
  );
  assert.deepEqual(groupSummaries({ positions: [], groupOpinions: [] }), []);
});

test("report visibility is sent as the form body the API expects", () => {
  assert.equal(reportVisibilityBody(true), "hidden=true");
  assert.equal(reportVisibilityBody(false), "hidden=false");
});

test("regenerate request is sent as a form body", () => {
  assert.equal(regenerateBody("report"), "type=report");
  assert.equal(regenerateBody("group"), "type=group");
});
