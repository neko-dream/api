// 管理画面のセッション編集フォームの変換と検証。サーバー（admin_usecase.EditTalkSession）と同じ条件にそろえる。

export interface ManagedTalkSession {
  talkSessionID: string;
  theme: string;
  description: string;
  owner: { userID: string; displayID: string; displayName: string; iconURL: string };
  scheduledEndTime: string;
  city?: string;
  prefecture?: string;
  thumbnailURL: string;
  hidden: boolean;
  updatedAt: string;
  createdAt: string;
  hideTop: boolean;
  hideReport: boolean;
  disableAnalysis: boolean;
  restrictions: string[];
  hasSurvey: boolean;
  participantCount: number;
}

export interface RestrictionOption {
  key: string;
  description: string;
  dependsOn?: string[];
}

export interface TalkSessionEditValues {
  theme: string;
  description: string;
  /** <input type="datetime-local"> の値（ローカル時刻、分まで） */
  scheduledEndTime: string;
  prefecture: string;
  city: string;
  hideTop: boolean;
  hideReport: boolean;
  disableAnalysis: boolean;
  restrictions: string[];
}

export interface TalkSessionEditRequest extends Omit<TalkSessionEditValues, "scheduledEndTime"> {
  scheduledEndTime: string;
}

export type TalkSessionEditErrors = Partial<Record<keyof TalkSessionEditValues, string>>;

const SURVEY_ANSWERED = "survey.answered";
const THEME_MAX_LENGTH = 100;
const DESCRIPTION_MAX_LENGTH = 40000;
// Go の strings.TrimSpace と同じ空白
const edgeWhitespace = /^\p{White_Space}+|\p{White_Space}+$/gu;

const pad = (n: number) => String(n).padStart(2, "0");

export function toLocalDateTimeInput(iso: string): string {
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return "";
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function toEditValues(session: ManagedTalkSession): TalkSessionEditValues {
  return {
    theme: session.theme,
    description: session.description,
    scheduledEndTime: toLocalDateTimeInput(session.scheduledEndTime),
    prefecture: session.prefecture ?? "",
    city: session.city ?? "",
    hideTop: session.hideTop,
    hideReport: session.hideReport,
    disableAnalysis: session.disableAnalysis,
    restrictions: [...session.restrictions],
  };
}

export function toEditRequest(
  values: TalkSessionEditValues,
  session: ManagedTalkSession,
): TalkSessionEditRequest {
  // 入力欄は分までしか持たないので、変えていなければ元の値を送る（サーバーは終了日時の変更有無を厳密に比べる）
  const scheduledEndTime =
    values.scheduledEndTime === toLocalDateTimeInput(session.scheduledEndTime)
      ? session.scheduledEndTime
      : new Date(values.scheduledEndTime).toISOString();
  return {
    ...values,
    scheduledEndTime,
    prefecture: values.prefecture.replace(edgeWhitespace, ""),
    city: values.city.replace(edgeWhitespace, ""),
    restrictions: [...values.restrictions],
  };
}

export function validateEditValues(
  values: TalkSessionEditValues,
  session: ManagedTalkSession,
  options: RestrictionOption[],
): TalkSessionEditErrors {
  const errors: TalkSessionEditErrors = {};
  if (!values.theme.replace(edgeWhitespace, "")) {
    errors.theme = "テーマを入力してください";
  } else if (Array.from(values.theme).length > THEME_MAX_LENGTH) {
    errors.theme = `テーマは${THEME_MAX_LENGTH}文字以内で入力してください`;
  }
  if (Array.from(values.description).length > DESCRIPTION_MAX_LENGTH) {
    errors.description = `本文は${DESCRIPTION_MAX_LENGTH}文字以内で入力してください`;
  }
  if (!values.scheduledEndTime || Number.isNaN(new Date(values.scheduledEndTime).getTime())) {
    errors.scheduledEndTime = "終了日時を入力してください";
  }

  const selected = new Set(values.restrictions);
  for (const key of values.restrictions) {
    const option = options.find((o) => o.key === key);
    const missing = (option?.dependsOn ?? []).filter((dep) => !selected.has(dep));
    if (missing.length > 0) {
      const names = missing.map((dep) => options.find((o) => o.key === dep)?.description ?? dep);
      errors.restrictions = `「${option?.description ?? key}」には「${names.join("」「")}」も必要です`;
      break;
    }
  }
  if (!errors.restrictions && selected.has(SURVEY_ANSWERED) && !session.hasSurvey) {
    errors.restrictions = "アンケートがないセッションには「アンケート回答済み」を付けられません";
  }
  return errors;
}

export function addedRestrictions(before: string[], after: string[]): string[] {
  return after.filter((key) => !before.includes(key));
}

/** 参加者がいるセッションに参加制限を追加すると、既存の参加者が参加できなくなる場合がある */
export function needsRestrictionWarning(
  session: ManagedTalkSession,
  values: TalkSessionEditValues,
): boolean {
  return (
    session.participantCount > 0 &&
    addedRestrictions(session.restrictions, values.restrictions).length > 0
  );
}
