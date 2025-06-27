# Google Form風アンケート機能設計書

## 1. 概要

Kotohiro（ことひろ）に、Google Formのような柔軟性のあるアンケート機能を追加する。既存のTalkSessionとは独立した機能として実装し、将来的に連携可能な設計とする。

## 2. 機能要件

### 2.1 主要機能
- **柔軟な質問タイプ**:
  - テキスト（短文・長文）
  - 選択肢（単一選択・複数選択）
  - 数値
  - 日付・時刻
  - スケール（1-5、1-10など）

- **バリデーション**:
  - 必須/任意設定
  - 文字数制限
  - 数値範囲
  - 正規表現パターン
  - カスタムバリデーション

- **回答管理**:
  - リアルタイム集計
  - 個別回答の閲覧
  - CSV/JSONエクスポート
  - 統計分析

- **アクセス制御**:
  - 公開/限定公開
  - 組織内限定
  - パスワード保護
  - 回答回数制限（1人1回など）

### 2.2 技術要件
- 既存のKotohiroアーキテクチャに準拠
- Clean Architecture / CQRSパターンの採用
- 高いパフォーマンスとスケーラビリティ
- セキュリティ（個人情報保護）

### 2.3 実装の詳細

#### 2.3.1 スキーマの柔軟性
アンケートの質問と回答は動的に変化するため、以下の方針で実装：

- **質問の設定**: JSONB型で`options`フィールドに格納
  - 質問タイプごとに異なる設定を柔軟に保持
  - 選択肢、プレースホルダー、デフォルト値など

- **バリデーション**: JSONB型で`validation`フィールドに格納
  - 文字数制限、正規表現、数値範囲など
  - カスタムバリデーションルールも定義可能

- **回答データ**: JSONB型で`value`フィールドに格納
  - 質問タイプに応じた様々な形式のデータを保存
  - テキスト、数値、配列、オブジェクトなど

#### 2.3.2 エラーハンドリング戦略

**ドメインレベルのエラー**:
```go
// internal/domain/messages/error.go に追加
var (
    ErrSurveyNotFound = errors.New("survey not found")
    ErrSurveyAlreadyPublished = errors.New("survey is already published")
    ErrSurveyClosed = errors.New("survey is closed")
    ErrQuestionRequired = errors.New("this question is required")
    ErrInvalidAnswerFormat = errors.New("invalid answer format")
    ErrResponseAlreadySubmitted = errors.New("response already submitted")
)
```

**バリデーションエラー**:
- 各質問タイプに応じたバリデーション実装
- エラーメッセージは質問ごとにカスタマイズ可能
- 複数のバリデーションエラーをまとめて返却

**APIレベルのエラー**:
- 4xx系: クライアントエラー（バリデーション、権限など）
- 5xx系: サーバーエラー（DB接続エラーなど）
- 詳細なエラー情報をJSONで返却

## 3. ドメインモデル設計

### 3.1 Survey（アンケート）

```go
package survey

import (
    "context"
    "errors"
    "time"
    "github.com/neko-dream/server/internal/domain/model/shared"
    "github.com/neko-dream/server/internal/domain/model/user"
    "github.com/neko-dream/server/internal/domain/model/organization"
)

type SurveyStatus string

const (
    SurveyStatusDraft     SurveyStatus = "draft"
    SurveyStatusPublished SurveyStatus = "published"
    SurveyStatusClosed    SurveyStatus = "closed"
)

type (
    SurveyRepository interface {
        Create(ctx context.Context, survey *Survey) error
        Update(ctx context.Context, survey *Survey) error
        FindByID(ctx context.Context, surveyID shared.UUID[Survey]) (*Survey, error)
        FindByOrganizationID(ctx context.Context, organizationID shared.UUID[organization.Organization]) ([]*Survey, error)
    }

    Survey struct {
        surveyID            shared.UUID[Survey]
        title               string
        description         *string
        creatorID           shared.UUID[user.User]
        organizationID      *shared.UUID[organization.Organization]
        organizationAliasID *shared.UUID[organization.OrganizationAlias]
        status              SurveyStatus
        settings            SurveySettings
        startAt             *time.Time
        endAt               *time.Time
        createdAt           time.Time
        updatedAt           time.Time
    }
)

type SurveySettings struct {
    RequireLogin        bool   `json:"require_login"`
    AllowMultipleSubmit bool   `json:"allow_multiple_submit"`
    ShowProgressBar     bool   `json:"show_progress_bar"`
    RandomizeQuestions  bool   `json:"randomize_questions"`
    ConfirmationMessage string `json:"confirmation_message"`
}

func NewSurvey(
    surveyID shared.UUID[Survey],
    title string,
    description *string,
    creatorID shared.UUID[user.User],
    organizationID *shared.UUID[organization.Organization],
    organizationAliasID *shared.UUID[organization.OrganizationAlias],
    settings SurveySettings,
    startAt *time.Time,
    endAt *time.Time,
    createdAt time.Time,
) *Survey {
    return &Survey{
        surveyID:            surveyID,
        title:               title,
        description:         description,
        creatorID:           creatorID,
        organizationID:      organizationID,
        organizationAliasID: organizationAliasID,
        status:              SurveyStatusDraft,
        settings:            settings,
        startAt:             startAt,
        endAt:               endAt,
        createdAt:           createdAt,
        updatedAt:           createdAt,
    }
}

// ゲッターメソッド
func (s *Survey) SurveyID() shared.UUID[Survey] {
    return s.surveyID
}

func (s *Survey) Title() string {
    return s.title
}

func (s *Survey) Description() *string {
    return s.description
}

func (s *Survey) CreatorID() shared.UUID[user.User] {
    return s.creatorID
}

func (s *Survey) OrganizationID() *shared.UUID[organization.Organization] {
    return s.organizationID
}

func (s *Survey) OrganizationAliasID() *shared.UUID[organization.OrganizationAlias] {
    return s.organizationAliasID
}

func (s *Survey) Status() SurveyStatus {
    return s.status
}

// セッターメソッド
func (s *Survey) ChangeTitle(title string) {
    s.title = title
    s.updatedAt = time.Now()
}

func (s *Survey) ChangeDescription(description *string) {
    s.description = description
    s.updatedAt = time.Now()
}

func (s *Survey) Publish() error {
    if s.status != SurveyStatusDraft {
        return errors.New("survey is not in draft status")
    }
    s.status = SurveyStatusPublished
    s.updatedAt = time.Now()
    return nil
}

func (s *Survey) Close() error {
    if s.status != SurveyStatusPublished {
        return errors.New("survey is not published")
    }
    s.status = SurveyStatusClosed
    s.updatedAt = time.Now()
    return nil
}
```

### 3.2 Question（質問）

```go
package survey

type QuestionType string

const (
    QuestionTypeText          QuestionType = "text"
    QuestionTypeTextArea      QuestionType = "textarea"
    QuestionTypeRadio         QuestionType = "radio"
    QuestionTypeCheckbox      QuestionType = "checkbox"
    QuestionTypeDropdown      QuestionType = "dropdown"
    QuestionTypeNumber        QuestionType = "number"
    QuestionTypeDate          QuestionType = "date"
    QuestionTypeTime          QuestionType = "time"
    QuestionTypeDateTime      QuestionType = "datetime"
    QuestionTypeFile          QuestionType = "file"
    QuestionTypeScale         QuestionType = "scale"
    QuestionTypeGrid          QuestionType = "grid"
)

type (
    QuestionRepository interface {
        Create(ctx context.Context, question *Question) error
        Update(ctx context.Context, question *Question) error
        Delete(ctx context.Context, questionID shared.UUID[Question]) error
        FindByID(ctx context.Context, questionID shared.UUID[Question]) (*Question, error)
        FindBySurveyID(ctx context.Context, surveyID shared.UUID[Survey]) ([]*Question, error)
    }

    Question struct {
        questionID   shared.UUID[Question]
        surveyID     shared.UUID[Survey]
        questionType QuestionType
        title        string
        description  *string
        required     bool
        order        int
        options      []QuestionOption
        validation   *ValidationRule
    }
)

type QuestionOption struct {
    ID    string `json:"id"`
    Label string `json:"label"`
    Value string `json:"value"`
    Order int    `json:"order"`
}

type ValidationRule struct {
    Type       string         `json:"type"`
    Parameters map[string]any `json:"parameters"`
}
```

### 3.3 SurveyResponse（回答セッション）

```go
package survey

type ResponseStatus string

const (
    ResponseStatusInProgress ResponseStatus = "in_progress"
    ResponseStatusSubmitted  ResponseStatus = "submitted"
    ResponseStatusAbandoned  ResponseStatus = "abandoned"
)

type (
    SurveyResponseRepository interface {
        Create(ctx context.Context, response *SurveyResponse) error
        Update(ctx context.Context, response *SurveyResponse) error
        FindByID(ctx context.Context, responseID shared.UUID[SurveyResponse]) (*SurveyResponse, error)
        FindBySurveyID(ctx context.Context, surveyID shared.UUID[Survey]) ([]*SurveyResponse, error)
        FindByResponderID(ctx context.Context, responderID shared.UUID[user.User]) ([]*SurveyResponse, error)
    }

    SurveyResponse struct {
        surveyResponseID shared.UUID[SurveyResponse]
        surveyID         shared.UUID[Survey]
        responderID      *shared.UUID[user.User]
        sessionToken     string
        ipAddress        string
        userAgent        string
        startedAt        time.Time
        submittedAt      *time.Time
        status           ResponseStatus
    }
)
```

### 3.4 Answer（個別回答）

```go
package survey

type (
    AnswerRepository interface {
        Create(ctx context.Context, answer *Answer) error
        Update(ctx context.Context, answer *Answer) error
        Delete(ctx context.Context, answerID shared.UUID[Answer]) error
        FindByID(ctx context.Context, answerID shared.UUID[Answer]) (*Answer, error)
        FindByResponseID(ctx context.Context, responseID shared.UUID[SurveyResponse]) ([]*Answer, error)
        FindByQuestionID(ctx context.Context, questionID shared.UUID[Question]) ([]*Answer, error)
    }

    Answer struct {
        answerID         shared.UUID[Answer]
        surveyResponseID shared.UUID[SurveyResponse]
        questionID       shared.UUID[Question]
        value            AnswerValue
        answeredAt       time.Time
    }
)

type AnswerValue struct {
    Type  string `json:"type"`
    Value any    `json:"value"`
}
```

## 4. データベーススキーマ

### 4.1 基本テーブル

```sql
-- surveys table
CREATE TABLE IF NOT EXISTS surveys (
    survey_id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    creator_id UUID NOT NULL,
    organization_id UUID,
    organization_alias_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    settings JSONB NOT NULL DEFAULT '{}',
    start_at TIMESTAMP,
    end_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);

-- questions table
CREATE TABLE IF NOT EXISTS questions (
    question_id UUID PRIMARY KEY,
    survey_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    required BOOLEAN NOT NULL DEFAULT false,
    "order" INTEGER NOT NULL,
    options JSONB,
    validation JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT (now()),
    updated_at TIMESTAMP NOT NULL DEFAULT (now())
);

-- survey_responses table
CREATE TABLE IF NOT EXISTS survey_responses (
    survey_response_id UUID PRIMARY KEY,
    survey_id UUID NOT NULL,
    responder_id UUID,
    session_token VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    started_at TIMESTAMP NOT NULL DEFAULT (now()),
    submitted_at TIMESTAMP,
    status VARCHAR(20) NOT NULL DEFAULT 'in_progress'
);

-- answers table
CREATE TABLE IF NOT EXISTS answers (
    answer_id UUID PRIMARY KEY,
    survey_response_id UUID NOT NULL,
    question_id UUID NOT NULL,
    value JSONB NOT NULL,
    answered_at TIMESTAMP NOT NULL DEFAULT (now())
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_surveys_creator_id ON surveys(creator_id);
CREATE INDEX IF NOT EXISTS idx_surveys_organization_id ON surveys(organization_id);
CREATE INDEX IF NOT EXISTS idx_surveys_status ON surveys(status);
CREATE INDEX IF NOT EXISTS idx_questions_survey_id ON questions(survey_id);
CREATE INDEX IF NOT EXISTS idx_survey_responses_survey_id ON survey_responses(survey_id);
CREATE INDEX IF NOT EXISTS idx_survey_responses_responder_id ON survey_responses(responder_id);
CREATE INDEX IF NOT EXISTS idx_answers_survey_response_id ON answers(survey_response_id);
CREATE INDEX IF NOT EXISTS idx_answers_question_id ON answers(question_id);

-- ユニーク制約
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_user_survey_response ON survey_responses(survey_id, responder_id) WHERE responder_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_response_question_answer ON answers(survey_response_id, question_id);

-- 外部キー制約
ALTER TABLE questions ADD CONSTRAINT fk_questions_survey_id
    FOREIGN KEY (survey_id) REFERENCES surveys(survey_id) ON DELETE CASCADE;
ALTER TABLE survey_responses ADD CONSTRAINT fk_survey_responses_survey_id
    FOREIGN KEY (survey_id) REFERENCES surveys(survey_id) ON DELETE CASCADE;
ALTER TABLE answers ADD CONSTRAINT fk_answers_survey_response_id
    FOREIGN KEY (survey_response_id) REFERENCES survey_responses(survey_response_id) ON DELETE CASCADE;
ALTER TABLE answers ADD CONSTRAINT fk_answers_question_id
    FOREIGN KEY (question_id) REFERENCES questions(question_id) ON DELETE CASCADE;
```

### 4.2 JSONBフィールドの構造例

#### 4.2.1 options JSONB

```json
// テキスト型（text/textarea）
{
  "placeholder": "回答を入力してください",
  "defaultValue": ""
}

// 単一選択（radio/dropdown）
{
  "choices": [
    {"id": "opt1", "label": "とても満足", "value": "5"},
    {"id": "opt2", "label": "満足", "value": "4"},
    {"id": "opt3", "label": "普通", "value": "3"},
    {"id": "opt4", "label": "不満", "value": "2"},
    {"id": "opt5", "label": "とても不満", "value": "1"}
  ],
  "allowOther": true,
  "otherLabel": "その他（具体的に）",
  "randomizeOrder": false
}

// 複数選択（checkbox）
{
  "choices": [
    {"id": "feat1", "label": "使いやすさ", "value": "usability"},
    {"id": "feat2", "label": "デザイン", "value": "design"},
    {"id": "feat3", "label": "価格", "value": "price"},
    {"id": "feat4", "label": "サポート", "value": "support"}
  ],
  "minSelections": 1,
  "maxSelections": 3,
  "allowOther": true
}

// スケール型
{
  "min": 1,
  "max": 10,
  "minLabel": "全く同意しない",
  "maxLabel": "強く同意する",
  "showLabels": true
}

// 日付型
{
  "dateFormat": "YYYY-MM-DD",
  "minDate": "2024-01-01",
  "maxDate": "2024-12-31",
  "disabledDates": ["2024-01-01", "2024-12-25"]
}
```

#### 4.2.2 validation JSONB

```json
// テキスト型のバリデーション
{
  "minLength": 10,
  "maxLength": 500,
  "pattern": "^[A-Za-z0-9\\s]+$",
  "patternMessage": "英数字とスペースのみ使用可能です"
}

// 数値型のバリデーション
{
  "min": 0,
  "max": 100,
  "step": 0.5,
  "decimalPlaces": 2
}

// メールアドレスのバリデーション
{
  "type": "email",
  "allowedDomains": ["company.com", "example.com"],
  "blockedDomains": ["tempmail.com"]
}

// 電話番号のバリデーション
{
  "type": "phone",
  "format": "JP", // 国コード
  "pattern": "^0\\d{9,10}$"
}

// カスタムバリデーション
{
  "custom": {
    "function": "validateZipCode",
    "message": "有効な郵便番号を入力してください"
  }
}
```

#### 4.2.3 value JSONB（answers table）

```json
// テキスト回答
{
  "text": "とても使いやすく、満足しています。"
}

// 数値回答
{
  "number": 85.5
}

// 単一選択回答
{
  "selected": "opt1"
}

// 単一選択（その他を選択）
{
  "selected": "other",
  "otherValue": "カスタムオプションの内容"
}

// 複数選択回答
{
  "selected": ["feat1", "feat3", "feat4"]
}

// 複数選択（その他を含む）
{
  "selected": ["feat1", "other"],
  "otherValue": "セキュリティ機能"
}

// スケール回答
{
  "scale": 8
}

// 日付回答
{
  "date": "2024-06-15"
}

// 日時回答
{
  "datetime": "2024-06-15T14:30:00Z"
}

// 複合型回答（住所など）
{
  "address": {
    "postalCode": "150-0001",
    "prefecture": "東京都",
    "city": "渋谷区",
    "street": "神宮前1-1-1",
    "building": "〇〇ビル5F"
  }
}
```

## 5. APIエンドポイント設計

### 5.1 TypeSpec定義

```typescript
// tsp/models/survey.tsp
import "@typespec/http";
import "@typespec/rest";
import "./common.tsp";

namespace kotohiro;

using TypeSpec.Http;
using TypeSpec.Rest;

@summary("アンケートステータス")
enum SurveyStatus {
  draft: "draft",
  published: "published",
  closed: "closed"
}

@summary("質問タイプ")
enum QuestionType {
  text: "text",
  textarea: "textarea",
  radio: "radio",
  checkbox: "checkbox",
  dropdown: "dropdown",
  number: "number",
  date: "date",
  time: "time",
  datetime: "datetime",
  scale: "scale"
}

@summary("回答ステータス")
enum ResponseStatus {
  inProgress: "in_progress",
  submitted: "submitted",
  abandoned: "abandoned"
}

@summary("アンケート")
model Survey {
  surveyID: string;
  title: string;
  description?: string;
  creatorID: string;
  organizationID?: string;
  organizationAliasID?: string;
  status: SurveyStatus;
  settings: SurveySettings;
  startAt?: utcDateTime;
  endAt?: utcDateTime;
  createdAt: utcDateTime;
  updatedAt: utcDateTime;
}

@summary("アンケート設定")
model SurveySettings {
  requireLogin: boolean = false;
  allowMultipleSubmit: boolean = false;
  showProgressBar: boolean = true;
  randomizeQuestions: boolean = false;
  confirmationMessage?: string;
}

@summary("質問")
model Question {
  questionID: string;
  surveyID: string;
  type: QuestionType;
  title: string;
  description?: string;
  required: boolean = false;
  order: int32;
  options?: QuestionOption[];
  validation?: ValidationRule;
}

@summary("質問の選択肢")
model QuestionOption {
  id: string;
  label: string;
  value: string;
  order: int32;
}

@summary("バリデーションルール")
model ValidationRule {
  type: string;
  parameters: Record<unknown>;
}

@summary("回答セッション")
model SurveyResponse {
  surveyResponseID: string;
  surveyID: string;
  responderID?: string;
  sessionToken?: string;
  startedAt: utcDateTime;
  submittedAt?: utcDateTime;
  status: ResponseStatus;
}

@summary("個別回答")
model Answer {
  answerID: string;
  surveyResponseID: string;
  questionID: string;
  value: unknown;
  answeredAt: utcDateTime;
}

// リクエスト/レスポンスモデル
@summary("アンケート作成リクエスト")
model CreateSurveyRequest {
  title: string;
  description?: string;
  organizationID?: string;
  organizationAliasID?: string;
  settings?: SurveySettings;
  startAt?: utcDateTime;
  endAt?: utcDateTime;
}

@summary("アンケート更新リクエスト")
model UpdateSurveyRequest {
  title?: string;
  description?: string;
  settings?: SurveySettings;
  startAt?: utcDateTime;
  endAt?: utcDateTime;
}

@summary("アンケートフィルター")
model SurveyFilter {
  organizationID?: string;
  creatorID?: string;
  status?: SurveyStatus;
  offset?: int32 = 0;
  limit?: int32 = 20;
}

@summary("質問作成リクエスト")
model CreateQuestionRequest {
  type: QuestionType;
  title: string;
  description?: string;
  required?: boolean = false;
  order?: int32;
  options?: QuestionOption[];
  validation?: ValidationRule;
}

@summary("質問更新リクエスト")
model UpdateQuestionRequest {
  type?: QuestionType;
  title?: string;
  description?: string;
  required?: boolean;
  order?: int32;
  options?: QuestionOption[];
  validation?: ValidationRule;
}

@summary("質問順序変更リクエスト")
model ReorderRequest {
  questionOrders: QuestionOrder[];
}

@summary("質問順序")
model QuestionOrder {
  questionID: string;
  order: int32;
}

@summary("回答送信リクエスト")
model SubmitAnswerRequest {
  questionID: string;
  value: unknown;
}

@summary("回答フィルター")
model ResponseFilter {
  responderID?: string;
  status?: ResponseStatus;
  startedFrom?: utcDateTime;
  startedTo?: utcDateTime;
  submittedFrom?: utcDateTime;
  submittedTo?: utcDateTime;
  offset?: int32 = 0;
  limit?: int32 = 20;
}

@summary("ページネーションレスポンス")
model PagedResponse<T> {
  items: T[];
  totalCount: int32;
  offset: int32;
  limit: int32;
  hasMore: boolean;
}

@summary("集計結果")
model SurveyResults {
  surveyID: string;
  totalResponses: int32;
  completedResponses: int32;
  inProgressResponses: int32;
  abandonedResponses: int32;
  questionResults: QuestionResult[];
}

@summary("質問ごとの集計結果")
model QuestionResult {
  questionID: string;
  questionTitle: string;
  questionType: QuestionType;
  totalAnswers: int32;
  results: AnswerStatistics;
}

@summary("回答統計")
model AnswerStatistics {
  // 選択肢型の集計
  optionCounts?: OptionCount[];
  // テキスト型の回答リスト
  textAnswers?: string[];
  // 数値型の統計
  numericStats?: NumericStatistics;
}

@summary("選択肢カウント")
model OptionCount {
  optionID: string;
  optionLabel: string;
  count: int32;
  percentage: float32;
}

@summary("数値統計")
model NumericStatistics {
  min: float64;
  max: float64;
  average: float64;
  median: float64;
  standardDeviation: float64;
}

@summary("回答詳細レスポンス")
model SurveyResponseDetail {
  response: SurveyResponse;
  answers: AnswerDetail[];
}

@summary("回答詳細")
model AnswerDetail {
  questionID: string;
  questionTitle: string;
  answer: Answer;
}
```

### 5.2 エンドポイント一覧

```typescript
// tsp/routes/survey.tsp
import "../models/survey.tsp";

using kotohiro;

@route("/surveys")
@tag("survey")
interface SurveyService {
  @post
  @extension("x-ogen-operation-group", "survey")
  @summary("アンケート作成")
  createSurvey(@body request: CreateSurveyRequest): Survey;

  @get
  @extension("x-ogen-operation-group", "survey")
  @summary("アンケート一覧取得")
  listSurveys(@query filter?: SurveyFilter): PagedResponse<Survey>;

  @get
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}")
  @summary("アンケート詳細取得")
  getSurvey(@path surveyId: string): Survey;

  @put
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}")
  @summary("アンケート更新")
  updateSurvey(@path surveyId: string, @body request: UpdateSurveyRequest): Survey;

  @delete
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}")
  @summary("アンケート削除")
  deleteSurvey(@path surveyId: string): void;

  @post
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/publish")
  @summary("アンケート公開")
  publishSurvey(@path surveyId: string): Survey;

  @post
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/close")
  @summary("アンケート終了")
  closeSurvey(@path surveyId: string): Survey;

  @post
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/questions")
  @summary("質問追加")
  addQuestion(@path surveyId: string, @body question: CreateQuestionRequest): Question;

  @put
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/questions/{questionId}")
  @summary("質問更新")
  updateQuestion(@path surveyId: string, @path questionId: string, @body question: UpdateQuestionRequest): Question;

  @delete
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/questions/{questionId}")
  @summary("質問削除")
  deleteQuestion(@path surveyId: string, @path questionId: string): void;

  @post
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/questions/reorder")
  @summary("質問順序変更")
  reorderQuestions(@path surveyId: string, @body order: ReorderRequest): void;

  @post
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/responses")
  @summary("回答開始")
  startResponse(@path surveyId: string): SurveyResponse;

  @post
  @extension("x-ogen-operation-group", "survey")
  @route("/responses/{responseId}/answers")
  @summary("回答送信")
  submitAnswer(@path responseId: string, @body answer: SubmitAnswerRequest): Answer;

  @post
  @extension("x-ogen-operation-group", "survey")
  @route("/responses/{responseId}/submit")
  @summary("回答完了")
  submitResponse(@path responseId: string): SurveyResponse;

  @get
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/results")
  @summary("集計結果取得")
  getSurveyResults(@path surveyId: string): SurveyResults;

  @get
  @extension("x-ogen-operation-group", "survey")
  @route("/{surveyId}/responses")
  @summary("個別回答一覧取得")
  listResponses(@path surveyId: string, @query filter?: ResponseFilter): PagedResponse<SurveyResponseDetail>;
}
```

## 6. 実装アーキテクチャ

### 6.1 条件付き質問表示（Question Logic）

質問の表示条件を管理する仕組み：

```go
// 質問の条件定義
type QuestionCondition struct {
    DependsOn   string      `json:"depends_on"` // 依存する質問ID
    Operator    string      `json:"operator"`   // equals, contains, greater_than等
    Value       interface{} `json:"value"`      // 比較値
}

// 使用例：「満足度が低い場合のみ理由を質問」
{
  "depends_on": "satisfaction_question_id",
  "operator": "less_than",
  "value": 3
}
```

### 6.2 回答セッション管理

**実装方針**: WebSocketを使用せず、HTTPベースのセッション管理を採用します。

#### セッショントークン管理
- 未ログインユーザーの回答を追跡
- UUIDv4形式で生成（crypto/rand使用）
- Cookieまたはローカルストレージに保存
- デフォルトTTL: 24時間（設定可能）

#### 自動保存機能
- **クライアント側**:
  - 定期的な自動保存（30秒ごと）
  - ユーザー操作時のデバウンス保存（2秒後）
  - ネットワークエラー時のローカルストレージバックアップ
- **サーバー側**:
  - Redis/Memcachedでセッション状態を管理
  - TTL付きでセッションデータを保存
  - デバイス変更の検出と記録

#### セッション復帰
- セッショントークンによる状態復元
- 回答済み質問と下書きの復元
- 現在のページ位置の保持
- 進捗率の表示

#### タイムアウト管理
- 非アクティブ検出（30分）
- 警告表示とセッション延長オプション
- 自動的な有効期限延長（アクティブ時）
- グレースフルな期限切れ処理

### 6.3 レイヤー構成

1. **Domain Layer** (`/internal/domain/model/survey/`)
   - エンティティ: Survey, Question, Answer, SurveyResponse
   - 値オブジェクト: QuestionType, ValidationRule, AnswerValue
   - ドメインサービス: SurveyValidator, AnswerValidator
   - リポジトリインターフェース

2. **Application Layer** (`/internal/application/`)
   - Use Cases:
     - CreateSurveyUseCase
     - UpdateSurveyUseCase
     - PublishSurveyUseCase
     - SubmitResponseUseCase
   - Query Services:
     - GetSurveyDetailQuery
     - GetSurveyResultsQuery
     - ListSurveysQuery
   - 例:
     ```go
     type CreateSurveyUseCase interface {
         Execute(context.Context, CreateSurveyUseCaseInput) (CreateSurveyUseCaseOutput, error)
     }
     ```

3. **Infrastructure Layer** (`/internal/infrastructure/`)
   - Repository実装（PostgreSQL + SQLC）
     - インターフェイスはdomain層に定義
     - 実装例:
       ```go
       type surveyRepository struct {
           *db.DBManager
       }

       func NewSurveyRepository(dbManager *db.DBManager) survey.SurveyRepository {
           return &surveyRepository{
               DBManager: dbManager,
           }
       }
       ```
     - SQLクエリは `/internal/infrastructure/persistence/sqlc/queries/survey/` に配置
     - クエリを書いたら`./scripts/gen.sh`でコード生成
     - 生成されたクエリは `db.GetQueries(ctx).CreateSurvey(ctx, params)` のように使用

4. **Presentation Layer** (`/internal/presentation/`)
   - HTTPハンドラー（OpenAPI/TypeSpecから生成）
   - 認証・認可ミドルウェア
   - セッション管理

### 実装の規約

1. **依存性注入**
   - 新しいインターフェイスを伴うコンストラクタを作成したら `internal/infrastructure/di/` に登録
   - ユースケース: `internal/infrastructure/di/application.go`
   - リポジトリ: `internal/infrastructure/di/infrastructure.go`

2. **UUIDの扱い**
   - 必ず `shared.UUID[T]` を使用
   - 新規作成: `shared.NewUUID[survey.Survey]()`
   - パース: `shared.ParseUUID[survey.Survey]("uuid-string")`

3. **エラーハンドリング**
   - ドメインエラーは `internal/domain/messages/error.go` に定義

4. **命名規則**
   - UseCase: `ExecuteXxxUseCase`
   - Query: `GetXxxQuery`, `BrowseXxxQuery`
   - Repository: `XxxRepository`

### 6.4 主要な処理フロー

#### アンケート回答フロー
1. ユーザーがアンケートにアクセス
2. SurveyResponseを作成（セッション開始）
3. 質問に対する回答を順次保存
4. 全回答完了後、SurveyResponseを提出済みに更新

#### 集計処理フロー
1. リアルタイムで回答を集計
2. 質問タイプに応じた統計処理
3. グラフ表示用のデータ変換
4. キャッシュによる高速化

## 7. 実装計画

### Phase 1: 基本機能実装（2週間）
- [ ] ドメインモデル実装
  - [ ] `/internal/domain/model/survey/` にエンティティ定義
  - [ ] リポジトリインターフェース定義
- [ ] データベースマイグレーション作成
  - [ ] `./scripts/migrate-create.sh survey_tables` でマイグレーション作成
- [ ] リポジトリ実装
  - [ ] `/internal/infrastructure/persistence/repository/survey_repository.go`
  - [ ] SQLクエリを `/internal/infrastructure/persistence/sqlc/queries/survey/` に作成
- [ ] 基本的なCRUD機能のユースケース
  - [ ] `/internal/application/usecase/survey_usecase/`
  - [ ] `/internal/application/query/survey_query/`
- [ ] TypeSpec定義とAPIエンドポイント実装
  - [ ] `/tsp/models/survey.tsp` にモデル定義
  - [ ] `/tsp/routes/survey.tsp` にエンドポイント定義
  - [ ] `./scripts/gen.sh` でコード生成
- [ ] DI設定
  - [ ] `/internal/infrastructure/di/application.go` にユースケース登録
  - [ ] `/internal/infrastructure/di/infrastructure.go` にリポジトリ登録
- [ ] テキスト・選択肢型質問の対応

### Phase 2: 高度な機能（2週間）
- [ ] バリデーション機能
- [ ] リアルタイム集計機能
- [ ] エクスポート機能

### Phase 3: UI実装（1週間）
- [ ] 管理画面でのアンケート作成・編集
- [ ] 回答画面の実装
- [ ] 集計結果表示画面

### Phase 4: 最適化とテスト（1週間）
- [ ] パフォーマンス最適化
- [ ] 統合テスト
- [ ] セキュリティ監査

## 8. 技術的考慮事項

### 8.1 セキュリティ

- XSS対策（HTMLサニタイジング）
- CSRF対策（トークン検証）
- SQLインジェクション対策（パラメータバインディング）
- 個人情報の暗号化（必要に応じて）
- レート制限（DoS攻撃対策）
- 入力値の厳密な検証

### 8.2 パフォーマンス

- 大量回答への対応（ページネーション）
  - デフォルト: 20件/ページ
  - 最大: 100件/ページ
- 集計結果のキャッシング
  - Redis使用（将来実装）
  - TTL: 5分（設定可能）
- データベースインデックスの最適化
  - 主要な検索フィールドにインデックス作成済み
- 非同期処理の活用
  - 大量データのエクスポートは非同期ジョブで処理

### 8.3 拡張性

- カスタムバリデーターの実装
  - バリデーターインターフェースを定義
  - プラグイン形式で追加可能
- 新しい質問タイプの追加が容易な設計
  - QuestionType enumに追加
  - 対応するバリデーターとレンダラーを実装
- Webhook機能（将来実装）
  - 回答完了時の通知
  - 外部システムとの連携

## 9. 既存システムとの統合

- 認証システム: 既存のJWT認証を活用
- 組織管理: Organization機能との連携
- 権限管理: 既存のロールベースアクセス制御
- 通知機能: 回答完了通知等（将来実装）

## 10. 今後の拡張案

- TalkSessionとの連携（アンケート結果をもとにした議論）
- AI分析機能（自由記述の要約・分類）
- テンプレート機能
- 多言語対応
- リアルタイムコラボレーション編集
