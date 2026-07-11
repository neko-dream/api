import "@typespec/http";

using Http;

namespace kotohiro {
  /**
   * アンケートの質問タイプ（選択肢系のみ）
   */
  enum TalkSessionSurveyQuestionType {
    singleChoice: "single_choice",
    dropdown: "dropdown",
  }

  /**
   * アンケート選択肢
   */
  model TalkSessionSurveyChoice {
    /**
     * 選択肢ID
     */
    choiceID: string;

    /**
     * 選択肢テキスト
     */
    text: string;

    /**
     * 表示順
     */
    displayOrder: int32;
  }

  /**
   * アンケート質問
   */
  model TalkSessionSurveyQuestion {
    /**
     * 質問ID
     */
    questionID: string;

    /**
     * 質問文
     */
    text: string;

    /**
     * 質問タイプ
     */
    type: TalkSessionSurveyQuestionType;

    /**
     * 必須回答かどうか
     */
    isRequired: boolean;

    /**
     * 「その他」入力を許可するか（single_choice / dropdown のみ有効）
     */
    allowOther: boolean;

    /**
     * 「その他」自由記述の最大文字数（allowOther の入力に適用）
     */
    maxLength?: int32 | null;

    /**
     * 表示順
     */
    displayOrder: int32;

    /**
     * 選択肢（single_choice / dropdown）
     */
    choices: TalkSessionSurveyChoice[];
  }

  /**
   * アンケート定義
   */
  model TalkSessionSurvey {
    /**
     * アンケートID
     */
    surveyID: string;

    /**
     * 紐付くセッションID
     */
    talkSessionID: string;

    /**
     * 質問一覧（display_order 昇順）
     */
    questions: TalkSessionSurveyQuestion[];
  }

  /**
   * 回答送信用：単一回答
   * 質問タイプは single_choice / dropdown のみ：
   * - selectedChoiceIDs を 1 件、または allowOther の場合 other
   */
  model TalkSessionSurveyAnswerInput {
    /**
     * 質問ID
     */
    questionID: string;

    /**
     * 選択された choice_id 配列（single_choice / dropdown は 1 件）
     */
    selectedChoiceIDs?: string[];

    /**
     * 「その他」自由記述（allowOther=true の場合に使用）
     */
    other?: string | null;
  }

  /**
   * アンケート回答（取得用）
   */
  model TalkSessionSurveyAnswer {
    /**
     * 質問ID
     */
    questionID: string;

    /**
     * 選択された choice_id 配列
     */
    selectedChoiceIDs: string[];

    /**
     * 「その他」自由記述
     */
    other?: string | null;
  }

  /**
   * 回答送信レコード
   */
  model TalkSessionSurveySubmission {
    /**
     * 送信ID
     */
    submissionID: string;

    /**
     * ユーザーID
     */
    userID: string;

    /**
     * 送信日時
     */
    submittedAt: utcDateTime;

    /**
     * 回答（質問ごと）
     */
    answers: TalkSessionSurveyAnswer[];
  }

  /**
   * アンケート作成用：質問入力
   */
  model TalkSessionSurveyQuestionInput {
    /**
     * 既存質問を編集する場合は questionID を指定。新規追加時は省略。
     */
    questionID?: string | null;

    text: string;
    type: TalkSessionSurveyQuestionType;
    isRequired: boolean;
    allowOther: boolean;
    maxLength?: int32 | null;
    displayOrder: int32;
    choices: TalkSessionSurveyChoiceInput[];
  }

  /**
   * アンケート作成用：選択肢入力
   */
  model TalkSessionSurveyChoiceInput {
    /**
     * 既存選択肢を編集する場合は choiceID を指定。新規追加時は省略。
     */
    choiceID?: string | null;

    text: string;
    displayOrder: int32;
  }

  /**
   * 集計：選択肢ごとの件数
   */
  model TalkSessionSurveyChoiceCount {
    choiceID: string;
    text: string;
    count: int32;
  }

  /**
   * 集計：質問ごとの結果（single_choice / dropdown）
   * choiceCounts + otherCount を持つ。
   */
  model TalkSessionSurveyQuestionSummary {
    questionID: string;
    text: string;
    type: TalkSessionSurveyQuestionType;
    /**
     * この質問に有効な回答を行ったユーザー数
     */
    answerCount: int32;
    /**
     * 選択肢ごとの件数
     */
    choiceCounts: TalkSessionSurveyChoiceCount[];
    /**
     * allowOther=true のみ：その他テキストを記入したユーザー数
     */
    otherCount: int32;
  }

  /**
   * 集計：survey 全体
   */
  model TalkSessionSurveySummary {
    surveyID: string;
    talkSessionID: string;
    /**
     * 送信ユーザー総数
     */
    submissionCount: int32;
    /**
     * 質問ごとの集計（display_order 昇順、削除済みは除外）
     */
    questionSummaries: TalkSessionSurveyQuestionSummary[];
  }
}
