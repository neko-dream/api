import "@typespec/http";

using Http;

namespace kotohiro {
  /**
   * アンケートの質問タイプ
   */
  enum TalkSessionSurveyQuestionType {
    freeText: "free_text",
    singleChoice: "single_choice",
    multiChoice: "multi_choice",
    rating: "rating",
    date: "date",
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
     * 「その他」入力を許可するか（single_choice / multi_choice / dropdown のみ有効）
     */
    allowOther: boolean;

    /**
     * 自由記述の最大文字数（free_text または allowOther の入力に適用）
     */
    maxLength?: int32 | null;

    /**
     * rating の下限値（rating のみ有効）
     */
    minValue?: int32 | null;

    /**
     * rating の上限値（rating のみ有効）
     */
    maxValue?: int32 | null;

    /**
     * rating の下限ラベル（例：「全く満足していない」）
     */
    minLabel?: string | null;

    /**
     * rating の上限ラベル（例：「非常に満足している」）
     */
    maxLabel?: string | null;

    /**
     * 表示順
     */
    displayOrder: int32;

    /**
     * 選択肢（single_choice / multi_choice / dropdown のみ）
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
   * 質問タイプによって意味が変わる：
   * - free_text: text を使用
   * - single_choice / dropdown: selectedChoiceIDs を 1 件、または allowOther の場合 other
   * - multi_choice: selectedChoiceIDs を 1 件以上、加えて allowOther の場合 other
   * - rating: rating を使用
   * - date: date を使用 (YYYY-MM-DD)
   */
  model TalkSessionSurveyAnswerInput {
    /**
     * 質問ID
     */
    questionID: string;

    /**
     * 自由記述（free_text）
     */
    text?: string | null;

    /**
     * 選択された choice_id 配列（single/multi/dropdown）
     */
    selectedChoiceIDs?: string[];

    /**
     * 「その他」自由記述（allowOther=true の choice 系で使用）
     */
    other?: string | null;

    /**
     * 評価値（rating）
     */
    rating?: int32 | null;

    /**
     * 日付（date, YYYY-MM-DD）
     */
    date?: string | null;
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
     * 自由記述（free_text）
     */
    text?: string | null;

    /**
     * 選択された choice_id 配列
     */
    selectedChoiceIDs: string[];

    /**
     * 「その他」自由記述
     */
    other?: string | null;

    /**
     * 評価値
     */
    rating?: int32 | null;

    /**
     * 日付
     */
    date?: string | null;
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
    minValue?: int32 | null;
    maxValue?: int32 | null;
    minLabel?: string | null;
    maxLabel?: string | null;
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
}
