import assert from "node:assert/strict";
import test from "node:test";
import {
  addedRestrictions,
  needsRestrictionWarning,
  toEditRequest,
  toEditValues,
  toLocalDateTimeInput,
  validateEditValues,
  type ManagedTalkSession,
  type RestrictionOption,
} from "../src/pages/talksessions/edit-form.ts";

const session: ManagedTalkSession = {
  talkSessionID: "ts-1",
  theme: "地域イベントのこれから",
  description: "<p>本文</p>",
  owner: { userID: "u-1", displayID: "owner", displayName: "オーナー", iconURL: "" },
  scheduledEndTime: "2026-10-20T03:04:05.678Z",
  prefecture: "東京都",
  thumbnailURL: "",
  hidden: false,
  updatedAt: "2026-10-01T00:00:00Z",
  createdAt: "2026-10-01T00:00:00Z",
  hideTop: false,
  hideReport: false,
  disableAnalysis: false,
  restrictions: ["demographics.prefecture"],
  hasSurvey: false,
  participantCount: 3,
};

const options: RestrictionOption[] = [
  { key: "demographics.prefecture", description: "都道府県" },
  { key: "demographics.city", description: "市区町村", dependsOn: ["demographics.prefecture"] },
  { key: "survey.answered", description: "アンケート回答済み" },
];

test("unchanged end time is sent as the original value, not the minute-rounded input", () => {
  const values = toEditValues(session);
  assert.equal(values.city, "");
  assert.equal(values.prefecture, "東京都");
  const request = toEditRequest(values, session);
  assert.equal(request.scheduledEndTime, session.scheduledEndTime);
  assert.deepEqual(request.restrictions, ["demographics.prefecture"]);
});

test("changed end time is sent as an ISO timestamp from local input", () => {
  const values = { ...toEditValues(session), scheduledEndTime: "2026-11-01T09:30" };
  const request = toEditRequest(values, session);
  assert.equal(request.scheduledEndTime, new Date("2026-11-01T09:30").toISOString());
  assert.notEqual(toLocalDateTimeInput(session.scheduledEndTime), "2026-11-01T09:30");
});

test("location text is trimmed and empty means unset", () => {
  const values = { ...toEditValues(session), prefecture: "  ", city: " 千代田区 " };
  const request = toEditRequest(values, session);
  assert.equal(request.prefecture, "");
  assert.equal(request.city, "千代田区");
});

test("validation matches the server limits and restriction rules", () => {
  const valid = toEditValues(session);
  assert.deepEqual(validateEditValues(valid, session, options), {});

  assert.ok(validateEditValues({ ...valid, theme: " " }, session, options).theme);
  assert.ok(validateEditValues({ ...valid, theme: "あ".repeat(101) }, session, options).theme);
  assert.deepEqual(validateEditValues({ ...valid, theme: "😀".repeat(100) }, session, options), {});
  assert.ok(validateEditValues({ ...valid, description: "a".repeat(40001) }, session, options).description);
  assert.ok(validateEditValues({ ...valid, scheduledEndTime: "" }, session, options).scheduledEndTime);

  const cityOnly = { ...valid, restrictions: ["demographics.city"] };
  assert.match(validateEditValues(cityOnly, session, options).restrictions ?? "", /都道府県/);

  const survey = { ...valid, restrictions: ["survey.answered"] };
  assert.match(validateEditValues(survey, session, options).restrictions ?? "", /アンケート/);
  assert.deepEqual(validateEditValues(survey, { ...session, hasSurvey: true }, options), {});
});

test("warning is shown only when restrictions are added to a session with participants", () => {
  assert.deepEqual(addedRestrictions(["a"], ["a", "b"]), ["b"]);
  assert.deepEqual(addedRestrictions(["a", "b"], ["a"]), []);

  const values = toEditValues(session);
  assert.equal(needsRestrictionWarning(session, values), false);
  const tightened = { ...values, restrictions: [...values.restrictions, "demographics.city"] };
  assert.equal(needsRestrictionWarning(session, tightened), true);
  assert.equal(needsRestrictionWarning({ ...session, participantCount: 0 }, tightened), false);
  assert.equal(needsRestrictionWarning(session, { ...values, restrictions: [] }), false);
});
