# アンケート回答システム詳細設計書

## 1. 概要

本ドキュメントは、Kotohiroのアンケート機能における回答システムの詳細設計を記述する。Google Formのような柔軟性と、エンタープライズレベルの信頼性を両立するシステムを目指す。対応する質問タイプは、テキスト、選択式（ラジオボタン、チェックボックス、ドロップダウン）、数値、スケール、日時の基本的な入力形式に対応する。

## 2. 回答値の型システム

### 2.1 基本設計思想

- **型安全性**: 各質問タイプに対応した専用の回答型を定義
- **拡張性**: インターフェースによる抽象化で新しい質問タイプを追加可能
- **バリデーション**: 各回答型が自身のバリデーションロジックを持つ

### 2.2 回答値インターフェース

```go
// internal/domain/model/survey/answer_value.go
package survey

import (
    "encoding/json"
    "fmt"
    "time"
    "regexp"
    "strings"
    "math"
)

// AnswerValue は質問タイプごとの回答値を表現する
type AnswerValue interface {
    Type() string
    Validate(question Question) error
    ToString() string // 人間が読める形式への変換
}
```

### 2.3 基本的な回答型

#### 2.3.1 テキスト回答

```go
// TextAnswer テキスト回答（短文・長文）
type TextAnswer struct {
    Text string
}

func (a TextAnswer) Type() string { return "text" }

func (a TextAnswer) Validate(q Question) error {
    if q.Required() && a.Text == "" {
        return fmt.Errorf("回答は必須です")
    }

    // バリデーションルールの適用
    if q.Validation() != nil {
        if minLen, ok := q.Validation().Parameters["minLength"].(float64); ok {
            if len(a.Text) < int(minLen) {
                return fmt.Errorf("最小文字数は%d文字です", int(minLen))
            }
        }
        if maxLen, ok := q.Validation().Parameters["maxLength"].(float64); ok {
            if len(a.Text) > int(maxLen) {
                return fmt.Errorf("最大文字数は%d文字です", int(maxLen))
            }
        }
        if pattern, ok := q.Validation().Parameters["pattern"].(string); ok {
            if matched, _ := regexp.MatchString(pattern, a.Text); !matched {
                if msg, ok := q.Validation().Parameters["patternMessage"].(string); ok {
                    return fmt.Errorf(msg)
                }
                return fmt.Errorf("入力形式が正しくありません")
            }
        }
    }
    return nil
}

func (a TextAnswer) ToString() string {
    return a.Text
}

```

#### 2.3.2 選択式回答

```go
// SelectionAnswer 選択式回答（単一・複数）
type SelectionAnswer struct {
    SelectedIDs []string
    OtherValue  *string
}

func (a SelectionAnswer) Type() string { return "selection" }

func (a SelectionAnswer) Validate(q Question) error {
    if q.Required() && len(a.SelectedIDs) == 0 {
        return fmt.Errorf("選択は必須です")
    }

    // 選択肢の検証
    validOptions := make(map[string]bool)
    for _, opt := range q.Options() {
        validOptions[opt.ID] = true
    }

    for _, id := range a.SelectedIDs {
        if id != "other" && !validOptions[id] {
            return fmt.Errorf("無効な選択肢: %s", id)
        }
    }

    // その他選択時の検証
    hasOther := false
    for _, id := range a.SelectedIDs {
        if id == "other" {
            hasOther = true
            break
        }
    }

    if hasOther && (a.OtherValue == nil || *a.OtherValue == "") {
        return fmt.Errorf("その他を選択した場合は詳細を入力してください")
    }

    // 選択数の検証（複数選択の場合）
    if q.Type() == QuestionTypeCheckbox && q.Validation() != nil {
        if minSel, ok := q.Validation().Parameters["minSelections"].(float64); ok {
            if len(a.SelectedIDs) < int(minSel) {
                return fmt.Errorf("最低%d個選択してください", int(minSel))
            }
        }
        if maxSel, ok := q.Validation().Parameters["maxSelections"].(float64); ok {
            if len(a.SelectedIDs) > int(maxSel) {
                return fmt.Errorf("最大%d個まで選択可能です", int(maxSel))
            }
        }
    }

    return nil
}

func (a SelectionAnswer) ToString() string {
    if a.OtherValue != nil && *a.OtherValue != "" {
        return fmt.Sprintf("%s (その他: %s)", strings.Join(a.SelectedIDs, ", "), *a.OtherValue)
    }
    return strings.Join(a.SelectedIDs, ", ")
}

```

#### 2.3.3 数値回答

```go
// NumericAnswer 数値回答
type NumericAnswer struct {
    Number float64
}

func (a NumericAnswer) Type() string { return "number" }

func (a NumericAnswer) Validate(q Question) error {
    if q.Validation() != nil {
        if min, ok := q.Validation().Parameters["min"].(float64); ok {
            if a.Number < min {
                return fmt.Errorf("最小値は%gです", min)
            }
        }
        if max, ok := q.Validation().Parameters["max"].(float64); ok {
            if a.Number > max {
                return fmt.Errorf("最大値は%gです", max)
            }
        }
        if step, ok := q.Validation().Parameters["step"].(float64); ok {
            // ステップ値の検証
            if min, ok := q.Validation().Parameters["min"].(float64); ok {
                diff := a.Number - min
                if remainder := math.Mod(diff, step); remainder != 0 {
                    return fmt.Errorf("値は%gの倍数である必要があります", step)
                }
            }
        }
    }
    return nil
}

func (a NumericAnswer) ToString() string {
    return fmt.Sprintf("%g", a.Number)
}

```

#### 2.3.4 スケール回答

```go
// ScaleAnswer スケール回答
type ScaleAnswer struct {
    Scale int
}

func (a ScaleAnswer) Type() string { return "scale" }

func (a ScaleAnswer) Validate(q Question) error {
    if q.Validation() != nil {
        if min, ok := q.Validation().Parameters["min"].(float64); ok {
            if a.Scale < int(min) {
                return fmt.Errorf("最小値は%dです", int(min))
            }
        }
        if max, ok := q.Validation().Parameters["max"].(float64); ok {
            if a.Scale > int(max) {
                return fmt.Errorf("最大値は%dです", int(max))
            }
        }
    }
    return nil
}

func (a ScaleAnswer) ToString() string {
    return fmt.Sprintf("%d", a.Scale)
}

```

#### 2.3.5 日時回答

```go
// DateTimeAnswer 日時回答
type DateTimeAnswer struct {
    DateTime time.Time
}

func (a DateTimeAnswer) Type() string { return "datetime" }

func (a DateTimeAnswer) Validate(q Question) error {
    if q.Validation() != nil {
        if minDateStr, ok := q.Validation().Parameters["minDate"].(string); ok {
            minDate, _ := time.Parse(time.RFC3339, minDateStr)
            if a.DateTime.Before(minDate) {
                return fmt.Errorf("最小日付は%sです", minDate.Format("2006-01-02"))
            }
        }
        if maxDateStr, ok := q.Validation().Parameters["maxDate"].(string); ok {
            maxDate, _ := time.Parse(time.RFC3339, maxDateStr)
            if a.DateTime.After(maxDate) {
                return fmt.Errorf("最大日付は%sです", maxDate.Format("2006-01-02"))
            }
        }
    }
    return nil
}

func (a DateTimeAnswer) ToString() string {
    return a.DateTime.Format("2006-01-02 15:04:05")
}

```


## 3. 回答エンティティとセッション管理

### 3.1 回答エンティティ

```go
// internal/domain/model/survey/answer.go
package survey

import (
    "context"
    "fmt"
    "time"
    "github.com/neko-dream/server/internal/domain/model/shared"
)

type Answer struct {
    answerID         shared.UUID[Answer]
    surveyResponseID shared.UUID[SurveyResponse]
    questionID       shared.UUID[Question]
    value            AnswerValue
    answeredAt       time.Time
    // Answerに関する追加情報
    duration         *time.Duration // 回答にかかった時間
    previousValues   []AnswerValue  // 修正履歴
    metadata         AnswerMetadata // メタデータ
}

type AnswerMetadata struct {
    DeviceType   string            // desktop, mobile, tablet
    BrowserInfo  string
    PageLoadTime time.Duration
    Interactions map[string]int    // クリック数、フォーカス回数など
}

// ファクトリメソッド
func NewAnswer(
    responseID shared.UUID[SurveyResponse],
    questionID shared.UUID[Question],
    value AnswerValue,
) (*Answer, error) {
    return &Answer{
        answerID:         shared.NewUUID[Answer](),
        surveyResponseID: responseID,
        questionID:       questionID,
        value:            value,
        answeredAt:       time.Now(),
        previousValues:   []AnswerValue{},
    }, nil
}

// 回答の更新（修正履歴を保持）
func (a *Answer) Update(newValue AnswerValue) error {
    a.previousValues = append(a.previousValues, a.value)
    a.value = newValue
    a.answeredAt = time.Now()
    return nil
}

// 回答時間の記録
func (a *Answer) SetDuration(duration time.Duration) {
    a.duration = &duration
}

// 回答の取り消し可能性チェック
func (a *Answer) CanBeReverted() bool {
    return len(a.previousValues) > 0
}

// 前の回答に戻す
func (a *Answer) Revert() error {
    if !a.CanBeReverted() {
        return fmt.Errorf("no previous answer to revert to")
    }

    lastIndex := len(a.previousValues) - 1
    a.value = a.previousValues[lastIndex]
    a.previousValues = a.previousValues[:lastIndex]
    a.answeredAt = time.Now()

    return nil
}

func (a *Answer) AnswerID() shared.UUID[Answer] {
    return a.answerID
}

func (a *Answer) SurveyResponseID() shared.UUID[SurveyResponse] {
    return a.surveyResponseID
}

func (a *Answer) QuestionID() shared.UUID[Question] {
    return a.questionID
}

func (a *Answer) Value() AnswerValue {
    return a.value
}

func (a *Answer) AnsweredAt() time.Time {
    return a.answeredAt
}

func (a *Answer) Duration() *time.Duration {
    return a.duration
}

func (a *Answer) PreviousValues() []AnswerValue {
    return a.previousValues
}

func (a *Answer) Metadata() AnswerMetadata {
    return a.metadata
}
```

### 3.2 回答セッション管理

```go
// internal/domain/model/survey/survey_response.go
package survey

import (
    "context"
    "fmt"
    "sync"
    "time"
    "github.com/neko-dream/server/internal/domain/model/shared"
    "github.com/neko-dream/server/internal/domain/model/user"
)

type SurveyResponse struct {
    surveyResponseID shared.UUID[SurveyResponse]
    surveyID         shared.UUID[Survey]
    responderID      *shared.UUID[user.User]
    sessionToken     string
    ipAddress        string
    userAgent        string
    startedAt        time.Time
    lastActivityAt   time.Time
    submittedAt      *time.Time
    status           ResponseStatus
    answers          map[shared.UUID[Question]]*Answer
    answerOrder      []shared.UUID[Question] // 回答順序の記録
    currentPage      int
    progress         SurveyProgress
    mutex            sync.RWMutex // 並行アクセス制御
}

type SurveyProgress struct {
    TotalQuestions    int
    AnsweredQuestions int
    RequiredQuestions int
    AnsweredRequired  int
    Percentage        float32
    EstimatedTimeLeft time.Duration
}

// ファクトリメソッド
func NewSurveyResponse(
    surveyID shared.UUID[Survey],
    responderID *shared.UUID[user.User],
    sessionToken string,
    ipAddress string,
    userAgent string,
    totalQuestions int,
    requiredQuestions int,
) *SurveyResponse {
    return &SurveyResponse{
        surveyResponseID: shared.NewUUID[SurveyResponse](),
        surveyID:         surveyID,
        responderID:      responderID,
        sessionToken:     sessionToken,
        ipAddress:        ipAddress,
        userAgent:        userAgent,
        startedAt:        time.Now(),
        lastActivityAt:   time.Now(),
        status:           ResponseStatusInProgress,
        answers:          make(map[shared.UUID[Question]]*Answer),
        answerOrder:      []shared.UUID[Question]{},
        currentPage:      0,
        progress: SurveyProgress{
            TotalQuestions:    totalQuestions,
            RequiredQuestions: requiredQuestions,
        },
    }
}

// 回答の保存（スレッドセーフ）
func (sr *SurveyResponse) SaveAnswer(
    questionID shared.UUID[Question],
    value AnswerValue,
    question Question,
) error {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()

    // バリデーション
    if err := value.Validate(question); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // 回答時間の計算
    var duration time.Duration
    if len(sr.answerOrder) > 0 {
        lastAnswerID := sr.answerOrder[len(sr.answerOrder)-1]
        if lastAnswer, exists := sr.answers[lastAnswerID]; exists {
            duration = time.Since(lastAnswer.AnsweredAt())
        }
    } else {
        duration = time.Since(sr.startedAt)
    }

    // 既存の回答があれば更新、なければ新規作成
    if existingAnswer, exists := sr.answers[questionID]; exists {
        if err := existingAnswer.Update(value); err != nil {
            return err
        }
    } else {
        newAnswer, err := NewAnswer(sr.surveyResponseID, questionID, value)
        if err != nil {
            return err
        }
        newAnswer.SetDuration(duration)
        sr.answers[questionID] = newAnswer
        sr.answerOrder = append(sr.answerOrder, questionID)
    }

    sr.lastActivityAt = time.Now()
    sr.updateProgress()

    return nil
}

// 進捗率の更新
func (sr *SurveyResponse) updateProgress() {
    sr.progress.AnsweredQuestions = len(sr.answers)

    if sr.progress.TotalQuestions > 0 {
        sr.progress.Percentage = float32(sr.progress.AnsweredQuestions) / float32(sr.progress.TotalQuestions) * 100

        // 推定残り時間の計算
        if sr.progress.AnsweredQuestions > 0 {
            avgTimePerQuestion := time.Since(sr.startedAt) / time.Duration(sr.progress.AnsweredQuestions)
            remainingQuestions := sr.progress.TotalQuestions - sr.progress.AnsweredQuestions
            sr.progress.EstimatedTimeLeft = avgTimePerQuestion * time.Duration(remainingQuestions)
        }
    }
}

// ページング対応
func (sr *SurveyResponse) GetCurrentPageQuestions(
    allQuestions []Question,
    questionsPerPage int,
) []Question {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()

    start := sr.currentPage * questionsPerPage
    end := start + questionsPerPage

    if start >= len(allQuestions) {
        return []Question{}
    }

    if end > len(allQuestions) {
        end = len(allQuestions)
    }

    return allQuestions[start:end]
}

// 次のページへ
func (sr *SurveyResponse) NextPage() {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()

    sr.currentPage++
    sr.lastActivityAt = time.Now()
}

// 前のページへ
func (sr *SurveyResponse) PreviousPage() {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()

    if sr.currentPage > 0 {
        sr.currentPage--
    }
    sr.lastActivityAt = time.Now()
}

// 回答の提出
func (sr *SurveyResponse) Submit() error {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()

    if sr.status != ResponseStatusInProgress {
        return fmt.Errorf("response is not in progress")
    }

    sr.status = ResponseStatusSubmitted
    now := time.Now()
    sr.submittedAt = &now
    sr.lastActivityAt = now

    return nil
}

func (sr *SurveyResponse) SurveyResponseID() shared.UUID[SurveyResponse] {
    return sr.surveyResponseID
}

func (sr *SurveyResponse) Status() ResponseStatus {
    return sr.status
}

func (sr *SurveyResponse) Progress() SurveyProgress {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()
    return sr.progress
}

func (sr *SurveyResponse) GetAnswer(questionID shared.UUID[Question]) *Answer {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()
    return sr.answers[questionID]
}

func (sr *SurveyResponse) GetNewAnswers() []*Answer {
    sr.mutex.RLock()
    defer sr.mutex.RUnlock()

    answers := make([]*Answer, 0, len(sr.answers))
    for _, answer := range sr.answers {
        answers = append(answers, answer)
    }
    return answers
}
```

### 3.3 セッション復帰とタイムアウト管理

**実装方針**: WebSocketを使用せず、HTTPベースのセッショントークン管理を採用します。

#### 3.3.1 セッション管理の仕組み

- **セッショントークン**: CookieまたはローカルストレージでクライアントがTTL付きトークンを保持
- **自動保存**: クライアント側で定期的（30秒ごと）およびユーザー操作時に保存
- **セッション復帰**: セッショントークンによる状態復元
- **タイムアウト**: 非アクティブ時間に基づく警告と自動終了

```go
// internal/domain/model/survey/response_state.go
package survey

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    "crypto/rand"
    "encoding/base64"
    "github.com/neko-dream/server/internal/domain/model/shared"
    "github.com/neko-dream/server/internal/domain/model/user"
)

// ResponseState 回答セッションの永続化可能な状態
type ResponseState struct {
    ResponseID        shared.UUID[SurveyResponse]
    SurveyID          shared.UUID[Survey]
    CurrentPage       int
    CurrentQuestionID *shared.UUID[Question]
    AnsweredQuestions []string
    Progress          SurveyProgress
    LastSavedAt       time.Time
    LastActivityAt    time.Time  // 最後のユーザー操作時刻
    ExpiresAt         time.Time
    Metadata          StateMetadata
}

type StateMetadata struct {
    SaveCount              int
    TotalTimeSpent         time.Duration
    DeviceChanges          int
    NetworkInterruptions   int
    AutoSaveCount          int           // 自動保存回数
    LastUserAgent          string        // デバイス変更検出用
    ShowInactivityWarning  bool          // 非アクティブ警告フラグ
}

// ResponseSessionManager セッション管理
type ResponseSessionManager struct {
    stateStore     ResponseStateStore
    responseRepo   SurveyResponseRepository
    answerRepo     AnswerRepository
    surveyRepo     SurveyRepository
    sessionTimeout time.Duration
}

type ResponseStateStore interface {
    Save(ctx context.Context, sessionToken string, state *ResponseState) error
    Load(ctx context.Context, sessionToken string) (*ResponseState, error)
    Delete(ctx context.Context, sessionToken string) error
    ExtendExpiration(ctx context.Context, sessionToken string, duration time.Duration) error
}

func NewResponseSessionManager(
    stateStore ResponseStateStore,
    responseRepo SurveyResponseRepository,
    answerRepo AnswerRepository,
    surveyRepo SurveyRepository,
    sessionTimeout time.Duration,
) *ResponseSessionManager {
    return &ResponseSessionManager{
        stateStore:     stateStore,
        responseRepo:   responseRepo,
        answerRepo:     answerRepo,
        surveyRepo:     surveyRepo,
        sessionTimeout: sessionTimeout,
    }
}

// セッションの作成
func (m *ResponseSessionManager) CreateSession(
    ctx context.Context,
    survey *Survey,
    responderID *shared.UUID[user.User],
    ipAddress string,
    userAgent string,
) (*SurveyResponse, string, error) {
    sessionToken := m.generateSecureToken()

    // 質問数を取得
    questions, err := m.surveyRepo.GetQuestions(ctx, survey.SurveyID())
    if err != nil {
        return nil, "", fmt.Errorf("failed to get questions: %w", err)
    }

    requiredCount := 0
    for _, q := range questions {
        if q.Required() {
            requiredCount++
        }
    }

    response := NewSurveyResponse(
        survey.SurveyID(),
        responderID,
        sessionToken,
        ipAddress,
        userAgent,
        len(questions),
        requiredCount,
    )

    // レスポンスの保存
    if err := m.responseRepo.Create(ctx, response); err != nil {
        return nil, "", fmt.Errorf("failed to create response: %w", err)
    }

    // セッション状態の保存
    state := m.createStateFromResponse(response)
    state.ExpiresAt = time.Now().Add(m.sessionTimeout)

    if err := m.stateStore.Save(ctx, sessionToken, state); err != nil {
        return nil, "", fmt.Errorf("failed to save session state: %w", err)
    }

    return response, sessionToken, nil
}

// セッションの復帰
func (m *ResponseSessionManager) ResumeSession(
    ctx context.Context,
    sessionToken string,
) (*SurveyResponse, error) {
    state, err := m.stateStore.Load(ctx, sessionToken)
    if err != nil {
        return nil, fmt.Errorf("failed to load session: %w", err)
    }

    // タイムアウトチェック
    if time.Now().After(state.ExpiresAt) {
        m.stateStore.Delete(ctx, sessionToken)
        return nil, fmt.Errorf("session expired")
    }

    // 非アクティブ警告チェック（30分）
    inactivityDuration := time.Since(state.LastActivityAt)
    if inactivityDuration > 30*time.Minute {
        state.Metadata.ShowInactivityWarning = true
    }

    // 回答セッションの復元
    response, err := m.responseRepo.FindByID(ctx, state.ResponseID)
    if err != nil {
        return nil, err
    }

    // 回答データの復元
    answers, err := m.answerRepo.FindByResponseID(ctx, state.ResponseID)
    if err != nil {
        return nil, err
    }

    response.RestoreAnswers(answers)
    response.RestoreState(state)

    // セッションの延長
    m.stateStore.ExtendExpiration(ctx, sessionToken, m.sessionTimeout)

    return response, nil
}

// 自動保存
func (m *ResponseSessionManager) AutoSave(
    ctx context.Context,
    sessionToken string,
    response *SurveyResponse,
    userAgent string,
) error {
    // 回答の永続化
    for _, answer := range response.GetNewAnswers() {
        if err := m.answerRepo.Save(ctx, answer); err != nil { // 新規作成のみ
            return fmt.Errorf("failed to save answer: %w", err)
        }
    }

    // セッション状態の更新
    state := m.createStateFromResponse(response)
    state.Metadata.AutoSaveCount++
    state.LastSavedAt = time.Now()
    state.LastActivityAt = time.Now()

    // デバイス変更検出
    if state.Metadata.LastUserAgent != "" && state.Metadata.LastUserAgent != userAgent {
        state.Metadata.DeviceChanges++
    }
    state.Metadata.LastUserAgent = userAgent

    return m.stateStore.Save(ctx, sessionToken, state)
}

// セキュアなトークン生成
func (m *ResponseSessionManager) generateSecureToken() string {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
    return base64.URLEncoding.EncodeToString(b)
}

// レスポンスから状態を作成
func (m *ResponseSessionManager) createStateFromResponse(response *SurveyResponse) *ResponseState {
    answeredQuestions := make([]string, 0, len(response.answerOrder))
    for _, qID := range response.answerOrder {
        answeredQuestions = append(answeredQuestions, qID.String())
    }

    return &ResponseState{
        ResponseID:        response.surveyResponseID,
        SurveyID:          response.surveyID,
        CurrentPage:       response.currentPage,
        AnsweredQuestions: answeredQuestions,
        Progress:          response.progress,
        LastSavedAt:       time.Now(),
        Metadata: StateMetadata{
            TotalTimeSpent: time.Since(response.startedAt),
        },
    }
}

// RestoreAnswers SurveyResponseに回答を復元
func (sr *SurveyResponse) RestoreAnswers(answers []*Answer) {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()

    sr.answers = make(map[shared.UUID[Question]]*Answer)
    for _, answer := range answers {
        sr.answers[answer.QuestionID()] = answer
    }
}

// RestoreState SurveyResponseに状態を復元
func (sr *SurveyResponse) RestoreState(state *ResponseState) {
    sr.mutex.Lock()
    defer sr.mutex.Unlock()

    sr.currentPage = state.CurrentPage
    sr.progress = state.Progress
}
```

#### 3.3.2 インフラ層でのセッション管理実装

```go
// internal/infrastructure/persistence/repository/response_state_store.go
package repository

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/neko-dream/server/internal/domain/model/survey"
)

type redisResponseStateStore struct {
    client         *redis.Client
    defaultTTL     time.Duration
}

func NewRedisResponseStateStore(client *redis.Client) survey.ResponseStateStore {
    return &redisResponseStateStore{
        client:     client,
        defaultTTL: 24 * time.Hour, // デフォルト24時間
    }
}

// Save セッション状態の保存
func (s *redisResponseStateStore) Save(
    ctx context.Context,
    sessionToken string,
    state *survey.ResponseState,
) error {
    key := fmt.Sprintf("survey:session:%s", sessionToken)

    data, err := json.Marshal(state)
    if err != nil {
        return fmt.Errorf("failed to marshal state: %w", err)
    }

    // TTL付きで保存
    ttl := time.Until(state.ExpiresAt)
    if ttl <= 0 {
        ttl = s.defaultTTL
    }

    return s.client.Set(ctx, key, data, ttl).Err()
}

// Load セッション状態の読み込み
func (s *redisResponseStateStore) Load(
    ctx context.Context,
    sessionToken string,
) (*survey.ResponseState, error) {
    key := fmt.Sprintf("survey:session:%s", sessionToken)

    data, err := s.client.Get(ctx, key).Bytes()
    if err != nil {
        if err == redis.Nil {
            return nil, fmt.Errorf("session not found")
        }
        return nil, fmt.Errorf("failed to get session: %w", err)
    }

    var state survey.ResponseState
    if err := json.Unmarshal(data, &state); err != nil {
        return nil, fmt.Errorf("failed to unmarshal state: %w", err)
    }

    return &state, nil
}

// Delete セッション状態の削除
func (s *redisResponseStateStore) Delete(
    ctx context.Context,
    sessionToken string,
) error {
    key := fmt.Sprintf("survey:session:%s", sessionToken)
    return s.client.Del(ctx, key).Err()
}

// ExtendExpiration 有効期限の延長
func (s *redisResponseStateStore) ExtendExpiration(
    ctx context.Context,
    sessionToken string,
    duration time.Duration,
) error {
    key := fmt.Sprintf("survey:session:%s", sessionToken)
    return s.client.Expire(ctx, key, duration).Err()
}

// クリーンアップタスク（定期実行）
func (s *redisResponseStateStore) CleanupExpiredSessions(ctx context.Context) error {
    // Redisの自動期限切れ機能を使用するため、明示的なクリーンアップは不要
    // ただし、異常終了したセッションの統計情報収集などが必要な場合はここに実装
    return nil
}
```

#### 3.3.3 クライアント側での自動保存実装例

```javascript
// フロントエンド実装例（TypeScript）
class SurveySessionManager {
    private sessionToken: string;
    private autosaveInterval: number = 30000; // 30秒
    private autosaveTimer?: NodeJS.Timer;
    private lastActivity: number = Date.now();
    private pendingAnswers: Map<string, any> = new Map();

    constructor(sessionToken: string) {
        this.sessionToken = sessionToken;
        this.startAutosave();
        this.setupActivityTracking();
    }

    // 定期的な自動保存
    private startAutosave(): void {
        this.autosaveTimer = setInterval(() => {
            this.saveCurrentState();
        }, this.autosaveInterval);
    }

    // ユーザー操作の追跡
    private setupActivityTracking(): void {
        ['click', 'keypress', 'scroll'].forEach(event => {
            document.addEventListener(event, () => {
                this.lastActivity = Date.now();
            });
        });
    }

    // ユーザー入力時の処理（デバウンス付き）
    public onUserInput(questionId: string, value: any): void {
        this.lastActivity = Date.now();
        this.pendingAnswers.set(questionId, value);

        // 2秒後に保存（デバウンス）
        clearTimeout(this.saveTimeout);
        this.saveTimeout = setTimeout(() => {
            this.saveAnswer(questionId, value);
        }, 2000);
    }

    // 現在の状態を保存
    private async saveCurrentState(): Promise<void> {
        if (this.pendingAnswers.size === 0) {
            return;
        }

        try {
            const response = await fetch(`/api/responses/${this.responseId}/autosave`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
                body: JSON.stringify({
                    answers: Array.from(this.pendingAnswers.entries()).map(
                        ([questionId, value]) => ({ questionId, value })
                    ),
                    lastActivity: this.lastActivity,
                }),
            });

            if (response.ok) {
                this.pendingAnswers.clear();
            }
        } catch (error) {
            console.error('Autosave failed:', error);
            // ローカルストレージにバックアップ
            this.saveToLocalStorage();
        }
    }

    // ローカルストレージへのバックアップ
    private saveToLocalStorage(): void {
        const backup = {
            sessionToken: this.sessionToken,
            answers: Array.from(this.pendingAnswers.entries()),
            timestamp: Date.now(),
        };
        localStorage.setItem('survey_backup', JSON.stringify(backup));
    }

    // セッション復帰
    public async resumeSession(): Promise<SurveyState | null> {
        try {
            const response = await fetch(`/api/responses/session/${this.sessionToken}`, {
                credentials: 'include',
            });

            if (response.ok) {
                const state = await response.json();

                // 非アクティブ警告の表示
                if (state.showInactivityWarning) {
                    this.showInactivityWarning();
                }

                return state;
            }
        } catch (error) {
            console.error('Failed to resume session:', error);
        }

        // ローカルストレージから復元を試みる
        return this.restoreFromLocalStorage();
    }

    // 非アクティブ警告の表示
    private showInactivityWarning(): void {
        // UIに警告を表示する実装
        alert('30分以上操作がありませんでした。セッションを継続しますか？');
    }

    // クリーンアップ
    public destroy(): void {
        if (this.autosaveTimer) {
            clearInterval(this.autosaveTimer);
        }
        clearTimeout(this.saveTimeout);
    }
}
```

## 4. 回答の集計と分析

### 4.1 集計サービス

```go
// internal/domain/service/survey_analytics_service.go
package service

import (
    "context"
    "math"
    "sort"
    "strings"
    "time"
    "github.com/neko-dream/server/internal/domain/model/survey"
    "github.com/neko-dream/server/internal/domain/model/shared"
)

type SurveyAnalyticsService struct {
    answerRepo   survey.AnswerRepository
    responseRepo survey.SurveyResponseRepository
    surveyRepo   survey.SurveyRepository
}

func NewSurveyAnalyticsService(
    answerRepo survey.AnswerRepository,
    responseRepo survey.SurveyResponseRepository,
    surveyRepo survey.SurveyRepository,
) *SurveyAnalyticsService {
    return &SurveyAnalyticsService{
        answerRepo:   answerRepo,
        responseRepo: responseRepo,
        surveyRepo:   surveyRepo,
    }
}

// 質問ごとの集計結果
type QuestionStatistics struct {
    QuestionID     shared.UUID[survey.Question]
    QuestionTitle  string
    QuestionType   survey.QuestionType
    TotalResponses int
    ValidResponses int
    SkippedCount   int
    Statistics     interface{}
}

// テキスト回答の統計
type TextStatistics struct {
    Responses      []TextResponse
    WordFrequency  map[string]int
    AverageLength  float64
    CommonPhrases  []Phrase
}

type TextResponse struct {
    ResponseID shared.UUID[survey.SurveyResponse]
    Text       string
    Timestamp  time.Time
}

type Phrase struct {
    Text  string
    Count int
}

// 選択式回答の統計
type SelectionStatistics struct {
    Options      []OptionStatistic
    OtherValues  []string
    Mode         []string          // 最頻値
    ResponseRate float64
}

type OptionStatistic struct {
    OptionID    string
    OptionLabel string
    Count       int
    Percentage  float64
}

// 数値回答の統計
type NumericStatistics struct {
    Min               float64
    Max               float64
    Mean              float64
    Median            float64
    Mode              []float64
    StandardDeviation float64
    Variance          float64
    Percentiles       map[int]float64 // 25, 50, 75, 90, 95
    Distribution      []BinCount
    Outliers          []float64
}

type BinCount struct {
    Min   float64
    Max   float64
    Count int
}

// スケール回答の統計
type ScaleStatistics struct {
    NumericStatistics
    ResponseDistribution map[int]int
    NetPromoterScore     *float64    // NPS（該当する場合）
}

// 日時回答の統計
type DateTimeStatistics struct {
    Earliest     time.Time
    Latest       time.Time
    Mode         []time.Time
    Distribution map[string]int // 日付ごとの分布
    Patterns     []DatePattern
}

type DatePattern struct {
    Pattern     string  // "weekday", "weekend", "month_start", etc.
    Count       int
    Percentage  float64
}

// アンケート全体の集計
func (s *SurveyAnalyticsService) AggregateSurveyResults(
    ctx context.Context,
    surveyID shared.UUID[survey.Survey],
) (*survey.SurveyResults, error) {
    // 全回答セッションを取得
    responses, err := s.responseRepo.FindBySurveyID(ctx, surveyID)
    if err != nil {
        return nil, err
    }

    // ステータス別の集計
    totalResponses := len(responses)
    completedResponses := 0
    inProgressResponses := 0
    abandonedResponses := 0

    for _, response := range responses {
        switch response.Status() {
        case survey.ResponseStatusSubmitted:
            completedResponses++
        case survey.ResponseStatusInProgress:
            inProgressResponses++
        case survey.ResponseStatusAbandoned:
            abandonedResponses++
        }
    }

    // 質問ごとの集計
    questionResults := []survey.QuestionResult{}
    questions, err := s.surveyRepo.GetQuestions(ctx, surveyID)
    if err != nil {
        return nil, err
    }

    for _, question := range questions {
        result, err := s.AggregateQuestionResults(ctx, question)
        if err != nil {
            return nil, err
        }
        questionResults = append(questionResults, *result)
    }

    return &survey.SurveyResults{
        SurveyID:            surveyID,
        TotalResponses:      totalResponses,
        CompletedResponses:  completedResponses,
        InProgressResponses: inProgressResponses,
        AbandonedResponses:  abandonedResponses,
        CompletionRate:      float64(completedResponses) / float64(totalResponses) * 100,
        QuestionResults:     questionResults,
        GeneratedAt:         time.Now(),
    }, nil
}

// 質問ごとの集計
func (s *SurveyAnalyticsService) AggregateQuestionResults(
    ctx context.Context,
    question survey.Question,
) (*survey.QuestionResult, error) {
    answers, err := s.answerRepo.FindByQuestionID(ctx, question.QuestionID())
    if err != nil {
        return nil, err
    }

    result := &survey.QuestionResult{
        QuestionID:    question.QuestionID(),
        QuestionTitle: question.Title(),
        QuestionType:  question.Type(),
        TotalAnswers:  len(answers),
    }

    switch question.Type() {
    case survey.QuestionTypeText, survey.QuestionTypeTextArea:
        result.Statistics = s.aggregateTextAnswers(answers)

    case survey.QuestionTypeRadio, survey.QuestionTypeCheckbox, survey.QuestionTypeDropdown:
        result.Statistics = s.aggregateSelectionAnswers(answers, question)

    case survey.QuestionTypeNumber:
        result.Statistics = s.aggregateNumericAnswers(answers)

    case survey.QuestionTypeScale:
        result.Statistics = s.aggregateScaleAnswers(answers, question)

    case survey.QuestionTypeDate, survey.QuestionTypeDateTime:
        result.Statistics = s.aggregateDateTimeAnswers(answers)
    }

    return result, nil
}

// テキスト回答の集計
func (s *SurveyAnalyticsService) aggregateTextAnswers(answers []*survey.Answer) TextStatistics {
    responses := []TextResponse{}
    totalLength := 0
    wordFreq := make(map[string]int)

    for _, answer := range answers {
        if textAnswer, ok := answer.Value().(survey.TextAnswer); ok {
            responses = append(responses, TextResponse{
                ResponseID: answer.SurveyResponseID(),
                Text:       textAnswer.Text,
                Timestamp:  answer.AnsweredAt(),
            })

            totalLength += len(textAnswer.Text)

            // 単語頻度の計算
            words := s.tokenizeText(textAnswer.Text)
            for _, word := range words {
                wordFreq[word]++
            }
        }
    }

    avgLength := 0.0
    if len(responses) > 0 {
        avgLength = float64(totalLength) / float64(len(responses))
    }

    // 共通フレーズの抽出
    commonPhrases := s.extractCommonPhrases(responses)

    return TextStatistics{
        Responses:     responses,
        WordFrequency: wordFreq,
        AverageLength: avgLength,
        CommonPhrases: commonPhrases,
    }
}

// 選択式回答の集計
func (s *SurveyAnalyticsService) aggregateSelectionAnswers(
    answers []*survey.Answer,
    question survey.Question,
) SelectionStatistics {
    optionCounts := make(map[string]int)
    otherValues := []string{}
    totalValid := 0

    // 選択肢の初期化
    for _, option := range question.Options() {
        optionCounts[option.ID] = 0
    }

    for _, answer := range answers {
        if selAnswer, ok := answer.Value().(survey.SelectionAnswer); ok {
            totalValid++
            for _, selectedID := range selAnswer.SelectedIDs {
                if selectedID == "other" && selAnswer.OtherValue != nil {
                    otherValues = append(otherValues, *selAnswer.OtherValue)
                } else {
                    optionCounts[selectedID]++
                }
            }
        }
    }

    // 統計の作成
    options := []OptionStatistic{}
    for _, option := range question.Options() {
        count := optionCounts[option.ID]
        percentage := 0.0
        if totalValid > 0 {
            percentage = float64(count) / float64(totalValid) * 100
        }

        options = append(options, OptionStatistic{
            OptionID:    option.ID,
            OptionLabel: option.Label,
            Count:       count,
            Percentage:  percentage,
        })
    }

    // 最頻値の計算
    mode := s.calculateMode(optionCounts)

    responseRate := 0.0
    if len(answers) > 0 {
        responseRate = float64(totalValid) / float64(len(answers)) * 100
    }

    return SelectionStatistics{
        Options:      options,
        OtherValues:  otherValues,
        Mode:         mode,
        ResponseRate: responseRate,
    }
}

// 数値回答の集計
func (s *SurveyAnalyticsService) aggregateNumericAnswers(answers []*survey.Answer) NumericStatistics {
    values := []float64{}

    for _, answer := range answers {
        switch v := answer.Value().(type) {
        case survey.NumericAnswer:
            values = append(values, v.Number)
        case survey.ScaleAnswer:
            values = append(values, float64(v.Scale))
        }
    }

    if len(values) == 0 {
        return NumericStatistics{}
    }

    // ソート（統計計算用）
    sort.Float64s(values)

    stats := NumericStatistics{
        Min:    values[0],
        Max:    values[len(values)-1],
        Mean:   s.calculateMean(values),
        Median: s.calculateMedian(values),
        Mode:   s.calculateNumericMode(values),
    }

    // 標準偏差と分散
    stats.StandardDeviation = s.calculateStandardDeviation(values, stats.Mean)
    stats.Variance = stats.StandardDeviation * stats.StandardDeviation

    // パーセンタイル
    stats.Percentiles = map[int]float64{
        25: s.calculatePercentile(values, 25),
        50: s.calculatePercentile(values, 50),
        75: s.calculatePercentile(values, 75),
        90: s.calculatePercentile(values, 90),
        95: s.calculatePercentile(values, 95),
    }

    // 分布（ヒストグラム）
    stats.Distribution = s.calculateDistribution(values, 10) // 10ビン

    // 外れ値の検出
    stats.Outliers = s.detectOutliers(values)

    return stats
}

// 平均値の計算
func (s *SurveyAnalyticsService) calculateMean(values []float64) float64 {
    sum := 0.0
    for _, v := range values {
        sum += v
    }
    return sum / float64(len(values))
}

// 中央値の計算
func (s *SurveyAnalyticsService) calculateMedian(sortedValues []float64) float64 {
    n := len(sortedValues)
    if n%2 == 0 {
        return (sortedValues[n/2-1] + sortedValues[n/2]) / 2
    }
    return sortedValues[n/2]
}

// 標準偏差の計算
func (s *SurveyAnalyticsService) calculateStandardDeviation(values []float64, mean float64) float64 {
    sumSquaredDiff := 0.0
    for _, v := range values {
        diff := v - mean
        sumSquaredDiff += diff * diff
    }
    variance := sumSquaredDiff / float64(len(values))
    return math.Sqrt(variance)
}

// パーセンタイルの計算
func (s *SurveyAnalyticsService) calculatePercentile(sortedValues []float64, percentile int) float64 {
    n := len(sortedValues)
    index := float64(percentile) / 100 * float64(n-1)
    lower := int(math.Floor(index))
    upper := int(math.Ceil(index))

    if lower == upper {
        return sortedValues[lower]
    }

    // 線形補間
    weight := index - float64(lower)
    return sortedValues[lower]*(1-weight) + sortedValues[upper]*weight
}

// 外れ値の検出（IQR法）
func (s *SurveyAnalyticsService) detectOutliers(sortedValues []float64) []float64 {
    q1 := s.calculatePercentile(sortedValues, 25)
    q3 := s.calculatePercentile(sortedValues, 75)
    iqr := q3 - q1

    lowerBound := q1 - 1.5*iqr
    upperBound := q3 + 1.5*iqr

    outliers := []float64{}
    for _, v := range sortedValues {
        if v < lowerBound || v > upperBound {
            outliers = append(outliers, v)
        }
    }

    return outliers
}

// 文字列のトークン化（簡易版）
func (s *SurveyAnalyticsService) tokenizeText(text string) []string {
    // 実際の実装では形態素解析を使用
    words := strings.Fields(strings.ToLower(text))
    return words
}

// 共通フレーズの抽出（簡易版）
func (s *SurveyAnalyticsService) extractCommonPhrases(responses []TextResponse) []Phrase {
    // 実際の実装ではn-gramや自然言語処理を使用
    return []Phrase{}
}

// 最頻値の計算
func (s *SurveyAnalyticsService) calculateMode(counts map[string]int) []string {
    maxCount := 0
    for _, count := range counts {
        if count > maxCount {
            maxCount = count
        }
    }

    mode := []string{}
    for value, count := range counts {
        if count == maxCount {
            mode = append(mode, value)
        }
    }

    return mode
}

// 数値の最頻値
func (s *SurveyAnalyticsService) calculateNumericMode(values []float64) []float64 {
    counts := make(map[float64]int)
    for _, v := range values {
        counts[v]++
    }

    maxCount := 0
    for _, count := range counts {
        if count > maxCount {
            maxCount = count
        }
    }

    mode := []float64{}
    for value, count := range counts {
        if count == maxCount {
            mode = append(mode, value)
        }
    }

    return mode
}

// 分布の計算
func (s *SurveyAnalyticsService) calculateDistribution(values []float64, binCount int) []BinCount {
    if len(values) == 0 {
        return []BinCount{}
    }

    min := values[0]
    max := values[len(values)-1]
    binWidth := (max - min) / float64(binCount)

    bins := make([]BinCount, binCount)
    for i := 0; i < binCount; i++ {
        bins[i] = BinCount{
            Min:   min + float64(i)*binWidth,
            Max:   min + float64(i+1)*binWidth,
            Count: 0,
        }
    }

    for _, v := range values {
        binIndex := int((v - min) / binWidth)
        if binIndex >= binCount {
            binIndex = binCount - 1
        }
        bins[binIndex].Count++
    }

    return bins
}

// スケール回答の集計
func (s *SurveyAnalyticsService) aggregateScaleAnswers(
    answers []*survey.Answer,
    question survey.Question,
) ScaleStatistics {
    numericStats := s.aggregateNumericAnswers(answers)

    // スケール値の分布を計算
    distribution := make(map[int]int)
    for _, answer := range answers {
        if scaleAnswer, ok := answer.Value().(survey.ScaleAnswer); ok {
            distribution[scaleAnswer.Scale]++
        }
    }

    stats := ScaleStatistics{
        NumericStatistics:    numericStats,
        ResponseDistribution: distribution,
    }

    // NPS計算（0-10スケールの場合）
    if question.Validation() != nil {
        min, hasMin := question.Validation().Parameters["min"].(float64)
        max, hasMax := question.Validation().Parameters["max"].(float64)

        if hasMin && hasMax && min == 0 && max == 10 {
            promoters := 0
            detractors := 0
            total := 0

            for scale, count := range distribution {
                total += count
                if scale >= 9 {
                    promoters += count
                } else if scale <= 6 {
                    detractors += count
                }
            }

            if total > 0 {
                nps := float64(promoters-detractors) / float64(total) * 100
                stats.NetPromoterScore = &nps
            }
        }
    }

    return stats
}

// 日時回答の集計
func (s *SurveyAnalyticsService) aggregateDateTimeAnswers(answers []*survey.Answer) DateTimeStatistics {
    dates := []time.Time{}

    for _, answer := range answers {
        if dateAnswer, ok := answer.Value().(survey.DateTimeAnswer); ok {
            dates = append(dates, dateAnswer.DateTime)
        }
    }

    if len(dates) == 0 {
        return DateTimeStatistics{}
    }

    // ソート
    sort.Slice(dates, func(i, j int) bool {
        return dates[i].Before(dates[j])
    })

    // 分布の計算
    distribution := make(map[string]int)
    for _, date := range dates {
        key := date.Format("2006-01-02")
        distribution[key]++
    }

    // パターンの分析
    patterns := s.analyzeDatePatterns(dates)

    return DateTimeStatistics{
        Earliest:     dates[0],
        Latest:       dates[len(dates)-1],
        Distribution: distribution,
        Patterns:     patterns,
    }
}

// 日付パターンの分析
func (s *SurveyAnalyticsService) analyzeDatePatterns(dates []time.Time) []DatePattern {
    weekdayCount := 0
    weekendCount := 0

    for _, date := range dates {
        switch date.Weekday() {
        case time.Saturday, time.Sunday:
            weekendCount++
        default:
            weekdayCount++
        }
    }

    total := len(dates)
    patterns := []DatePattern{}

    if weekdayCount > 0 {
        patterns = append(patterns, DatePattern{
            Pattern:    "weekday",
            Count:      weekdayCount,
            Percentage: float64(weekdayCount) / float64(total) * 100,
        })
    }

    if weekendCount > 0 {
        patterns = append(patterns, DatePattern{
            Pattern:    "weekend",
            Count:      weekendCount,
            Percentage: float64(weekendCount) / float64(total) * 100,
        })
    }

    return patterns
}
```

## 5. 回答データの永続化

### 5.1 リポジトリインターフェース

```go
// internal/domain/model/survey/answer_repository.go
package survey

import (
    "context"
    "time"
)

type AnswerRepository interface {
    // 基本的なCRUD
    Save(ctx context.Context, answer *Answer) error // 新規作成のみ（Createの動作）
    SaveBatch(ctx context.Context, answers []*Answer) error
    Update(ctx context.Context, answer *Answer) error
    Delete(ctx context.Context, answerID shared.UUID[Answer]) error

    // 取得
    FindByID(ctx context.Context, answerID shared.UUID[Answer]) (*Answer, error)
    FindByResponseID(ctx context.Context, responseID shared.UUID[SurveyResponse]) ([]*Answer, error)
    FindByQuestionID(ctx context.Context, questionID shared.UUID[Question]) ([]*Answer, error)
    FindByResponseAndQuestion(ctx context.Context, responseID shared.UUID[SurveyResponse], questionID shared.UUID[Question]) (*Answer, error)

    // 一時保存
    SaveDraft(ctx context.Context, responseID shared.UUID[SurveyResponse], answers []*Answer) error
    LoadDraft(ctx context.Context, responseID shared.UUID[SurveyResponse]) ([]*Answer, error)
    DeleteDraft(ctx context.Context, responseID shared.UUID[SurveyResponse]) error

    // 統計用クエリ
    CountByQuestion(ctx context.Context, questionID shared.UUID[Question]) (int, error)
    GetAnswerDistribution(ctx context.Context, questionID shared.UUID[Question]) (map[string]int, error)
    GetAverageByQuestion(ctx context.Context, questionID shared.UUID[Question]) (float64, error)

    // 高度なクエリ
    FindByQuestionWithPagination(ctx context.Context, questionID shared.UUID[Question], offset, limit int) ([]*Answer, int, error)
    FindByTimeRange(ctx context.Context, surveyID shared.UUID[Survey], start, end time.Time) ([]*Answer, error)

    // バルク操作
    DeleteByResponseID(ctx context.Context, responseID shared.UUID[SurveyResponse]) error
    DeleteBySurveyID(ctx context.Context, surveyID shared.UUID[Survey]) error
}

type SurveyResponseRepository interface {
    // 基本的なCRUD
    Create(ctx context.Context, response *SurveyResponse) error
    Update(ctx context.Context, response *SurveyResponse) error
    Delete(ctx context.Context, responseID shared.UUID[SurveyResponse]) error

    // 取得
    FindByID(ctx context.Context, responseID shared.UUID[SurveyResponse]) (*SurveyResponse, error)
    FindBySurveyID(ctx context.Context, surveyID shared.UUID[Survey]) ([]*SurveyResponse, error)
    FindByResponderID(ctx context.Context, responderID shared.UUID[user.User]) ([]*SurveyResponse, error)
    FindBySessionToken(ctx context.Context, sessionToken string) (*SurveyResponse, error)

    // 条件付き取得
    FindByStatus(ctx context.Context, surveyID shared.UUID[Survey], status ResponseStatus) ([]*SurveyResponse, error)
    FindByTimeRange(ctx context.Context, surveyID shared.UUID[Survey], start, end time.Time) ([]*SurveyResponse, error)

    // 統計
    CountByStatus(ctx context.Context, surveyID shared.UUID[Survey]) (map[ResponseStatus]int, error)
    GetAverageCompletionTime(ctx context.Context, surveyID shared.UUID[Survey]) (time.Duration, error)
}
```

### 5.2 PostgreSQL実装

```go
// internal/infrastructure/persistence/repository/answer_repository.go
package repository

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "time"

    "braces.dev/errtrace"
    "github.com/google/uuid"
    "github.com/neko-dream/server/internal/domain/model/survey"
    "github.com/neko-dream/server/internal/domain/model/shared"
    "github.com/neko-dream/server/internal/infrastructure/persistence/db"
    "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
)

type answerRepository struct {
    *db.DBManager
}

func NewAnswerRepository(dbManager *db.DBManager) survey.AnswerRepository {
    return &answerRepository{
        DBManager: dbManager,
    }
}

// Save 回答の保存
func (r *answerRepository) Save(ctx context.Context, answer *survey.Answer) error {
    queries := r.GetQueries(ctx)

    // AnswerValueをJSONに変換
    // 注: ドメイン層のAnswerValueは純粋な値オブジェクトのため、
    // インフラ層でJSONマーシャリングを行う
    valueJSON, err := r.marshalAnswerValue(answer.Value())
    if err != nil {
        return errtrace.Wrap(fmt.Errorf("failed to marshal answer value: %w", err))
    }

    // メタデータをJSONに変換
    metadataJSON, err := json.Marshal(answer.Metadata())
    if err != nil {
        return errtrace.Wrap(fmt.Errorf("failed to marshal metadata: %w", err))
    }

    params := model.CreateAnswerParams{
        AnswerID:         answer.AnswerID().UUID(),
        SurveyResponseID: answer.SurveyResponseID().UUID(),
        QuestionID:       answer.QuestionID().UUID(),
        Value:            valueJSON,
        Metadata:         metadataJSON,
        AnsweredAt:       answer.AnsweredAt(),
    }

    if answer.Duration() != nil {
        duration := int32(answer.Duration().Seconds())
        params.Duration = sql.NullInt32{Int32: duration, Valid: true}
    }

    return errtrace.Wrap(queries.CreateAnswer(ctx, params))
}

// SaveBatch バッチ保存
func (r *answerRepository) SaveBatch(ctx context.Context, answers []*survey.Answer) error {
    return r.ExecTx(ctx, func(ctx context.Context) error {
        for _, answer := range answers {
            if err := r.Save(ctx, answer); err != nil {
                return err
            }
        }
        return nil
    })
}

// FindByResponseID レスポンスIDによる取得
func (r *answerRepository) FindByResponseID(
    ctx context.Context,
    responseID shared.UUID[survey.SurveyResponse],
) ([]*survey.Answer, error) {
    rows, err := r.GetQueries(ctx).FindAnswersByResponseID(ctx, responseID.UUID())
    if err != nil {
        return nil, errtrace.Wrap(err)
    }

    answers := make([]*survey.Answer, 0, len(rows))
    for _, row := range rows {
        answer, err := r.mapRowToAnswer(ctx, row)
        if err != nil {
            return nil, err
        }
        answers = append(answers, answer)
    }

    return answers, nil
}

// mapRowToAnswer 行データをAnswerモデルに変換
func (r *answerRepository) mapRowToAnswer(ctx context.Context, row model.Answer) (*survey.Answer, error) {
    // 質問タイプを取得
    questionType, err := r.getQuestionType(ctx, row.QuestionID)
    if err != nil {
        return nil, errtrace.Wrap(err)
    }

    // JSONから回答値を復元
    answerValue, err := r.unmarshalAnswerValue(row.Value, questionType)
    if err != nil {
        return nil, err
    }

    answer := survey.ReconstructAnswer(
        shared.MustParseUUID[survey.Answer](row.AnswerID.String()),
        shared.MustParseUUID[survey.SurveyResponse](row.SurveyResponseID.String()),
        shared.MustParseUUID[survey.Question](row.QuestionID.String()),
        answerValue,
        row.AnsweredAt,
    )

    if row.Duration.Valid {
        duration := time.Duration(row.Duration.Int32) * time.Second
        answer.SetDuration(duration)
    }

    return answer, nil
}

// marshalAnswerValue AnswerValueをJSONに変換（インフラ層の責務）
func (r *answerRepository) marshalAnswerValue(value survey.AnswerValue) (json.RawMessage, error) {
    // AnswerValueの型に応じて適切なJSON構造を作成
    switch v := value.(type) {
    case survey.TextAnswer:
        return json.Marshal(map[string]interface{}{
            "text": v.Text,
        })
    case survey.SelectionAnswer:
        return json.Marshal(map[string]interface{}{
            "selected": v.SelectedIDs,
            "otherValue": v.OtherValue,
        })
    case survey.NumericAnswer:
        return json.Marshal(map[string]interface{}{
            "number": v.Number,
        })
    case survey.ScaleAnswer:
        return json.Marshal(map[string]interface{}{
            "scale": v.Scale,
        })
    case survey.DateTimeAnswer:
        return json.Marshal(map[string]interface{}{
            "datetime": v.DateTime,
        })
    default:
        return nil, fmt.Errorf("unknown answer value type: %T", value)
    }
}

// unmarshalAnswerValue JSONから適切な回答値型に変換
func (r *answerRepository) unmarshalAnswerValue(
    data json.RawMessage,
    questionType survey.QuestionType,
) (survey.AnswerValue, error) {
    switch questionType {
    case survey.QuestionTypeText, survey.QuestionTypeTextArea:
        var dto struct {
            Text string `json:"text"`
        }
        if err := json.Unmarshal(data, &dto); err != nil {
            return nil, errtrace.Wrap(err)
        }
        return survey.TextAnswer{Text: dto.Text}, nil

    case survey.QuestionTypeRadio, survey.QuestionTypeCheckbox, survey.QuestionTypeDropdown:
        var dto struct {
            SelectedIDs []string `json:"selected"`
            OtherValue  *string  `json:"otherValue,omitempty"`
        }
        if err := json.Unmarshal(data, &dto); err != nil {
            return nil, errtrace.Wrap(err)
        }
        return survey.SelectionAnswer{
            SelectedIDs: dto.SelectedIDs,
            OtherValue:  dto.OtherValue,
        }, nil

    case survey.QuestionTypeNumber:
        var dto struct {
            Number float64 `json:"number"`
        }
        if err := json.Unmarshal(data, &dto); err != nil {
            return nil, errtrace.Wrap(err)
        }
        return survey.NumericAnswer{Number: dto.Number}, nil

    case survey.QuestionTypeScale:
        var dto struct {
            Scale int `json:"scale"`
        }
        if err := json.Unmarshal(data, &dto); err != nil {
            return nil, errtrace.Wrap(err)
        }
        return survey.ScaleAnswer{Scale: dto.Scale}, nil

    case survey.QuestionTypeDate, survey.QuestionTypeDateTime:
        var dto struct {
            DateTime time.Time `json:"datetime"`
        }
        if err := json.Unmarshal(data, &dto); err != nil {
            return nil, errtrace.Wrap(err)
        }
        return survey.DateTimeAnswer{DateTime: dto.DateTime}, nil

    default:
        return nil, fmt.Errorf("unknown question type: %s", questionType)
    }
}

// getQuestionType 質問タイプを取得（キャッシュまたはDB）
func (r *answerRepository) getQuestionType(ctx context.Context, questionID uuid.UUID) (survey.QuestionType, error) {
    question, err := querier.GetQueries(ctx).GetQuestionByID(ctx, questionID)
    if err != nil {
        return "", errtrace.Wrap(err)
    }
    return survey.QuestionType(question.Type), nil
}

// GetAnswerDistribution 回答分布の取得
func (r *answerRepository) GetAnswerDistribution(
    ctx context.Context,
    questionID shared.UUID[survey.Question],
) (map[string]int, error) {
    rows, err := r.GetQueries(ctx).GetAnswerDistribution(ctx, questionID.UUID())
    if err != nil {
        return nil, errtrace.Wrap(err)
    }

    distribution := make(map[string]int)
    for _, row := range rows {
        distribution[row.Value] = int(row.Count)
    }

    return distribution, nil
}

// DeleteByResponseID レスポンスIDで削除
func (r *answerRepository) DeleteByResponseID(
    ctx context.Context,
    responseID shared.UUID[survey.SurveyResponse],
) error {
    return errtrace.Wrap(r.GetQueries(ctx).DeleteAnswersByResponseID(ctx, responseID.UUID()))
}
```

### 5.3 SQLクエリ定義

```sql
-- internal/infrastructure/persistence/sqlc/queries/survey/answer.sql

-- name: CreateAnswer :exec
INSERT INTO answers (
    answer_id,
    survey_response_id,
    question_id,
    value,
    metadata,
    duration,
    answered_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: UpdateAnswer :exec
UPDATE answers SET
    value = $2,
    metadata = $3,
    answered_at = $4,
    updated_at = now()
WHERE answer_id = $1;

-- name: FindAnswerByID :one
SELECT * FROM answers
WHERE answer_id = $1;

-- name: FindAnswersByResponseID :many
SELECT * FROM answers
WHERE survey_response_id = $1
ORDER BY answered_at ASC;

-- name: FindAnswersByQuestionID :many
SELECT * FROM answers
WHERE question_id = $1
ORDER BY answered_at DESC;

-- name: FindAnswerByResponseAndQuestion :one
SELECT * FROM answers
WHERE survey_response_id = $1 AND question_id = $2;

-- name: GetAnswerDistribution :many
SELECT
    value->>'selected' as value,
    COUNT(*) as count
FROM answers
WHERE question_id = $1
GROUP BY value->>'selected'
ORDER BY count DESC;

-- name: GetNumericAverage :one
SELECT
    AVG((value->>'number')::float) as average
FROM answers
WHERE question_id = $1
    AND value->>'number' IS NOT NULL;

-- name: CountAnswersByQuestion :one
SELECT COUNT(*) FROM answers
WHERE question_id = $1;

-- name: DeleteAnswer :exec
DELETE FROM answers
WHERE answer_id = $1;

-- name: DeleteAnswersByResponseID :exec
DELETE FROM answers
WHERE survey_response_id = $1;

-- name: DeleteAnswersBySurveyID :exec
DELETE FROM answers
WHERE survey_response_id IN (
    SELECT survey_response_id
    FROM survey_responses
    WHERE survey_id = $1
);

-- name: GetAnswerTimeSeries :many
SELECT
    DATE_TRUNC('day', answered_at) as date,
    COUNT(*) as count
FROM answers
WHERE question_id = $1
    AND answered_at BETWEEN $2 AND $3
GROUP BY date
ORDER BY date;

-- name: SaveDraftAnswer :exec
INSERT INTO answer_drafts (
    answer_id,
    survey_response_id,
    question_id,
    value,
    saved_at
) VALUES (
    $1, $2, $3, $4, now()
) ON CONFLICT (survey_response_id, question_id)
DO UPDATE SET
    value = EXCLUDED.value,
    saved_at = EXCLUDED.saved_at;

-- name: LoadDraftAnswers :many
SELECT * FROM answer_drafts
WHERE survey_response_id = $1;

-- name: DeleteDraftAnswers :exec
DELETE FROM answer_drafts
WHERE survey_response_id = $1;

-- name: GetQuestionByID :one
SELECT question_id, type FROM questions
WHERE question_id = $1;
```

## 6. エクスポート機能

### 6.1 エクスポートサービス

```go
// internal/application/query/survey_query/export_responses_query.go
package survey_query

import (
    "context"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "strings"
    "time"

    "golang.org/x/text/encoding/japanese"
    "golang.org/x/text/transform"

    "github.com/neko-dream/server/internal/domain/model/survey"
    "github.com/neko-dream/server/internal/domain/model/shared"
)

type ExportResponsesQuery struct {
    surveyRepo   survey.SurveyRepository
    responseRepo survey.SurveyResponseRepository
    answerRepo   survey.AnswerRepository
    questionRepo survey.QuestionRepository
}

func NewExportResponsesQuery(
    surveyRepo survey.SurveyRepository,
    responseRepo survey.SurveyResponseRepository,
    answerRepo survey.AnswerRepository,
    questionRepo survey.QuestionRepository,
) *ExportResponsesQuery {
    return &ExportResponsesQuery{
        surveyRepo:   surveyRepo,
        responseRepo: responseRepo,
        answerRepo:   answerRepo,
        questionRepo: questionRepo,
    }
}

type ExportFormat string

const (
    ExportFormatCSV  ExportFormat = "csv"
    ExportFormatJSON ExportFormat = "json"
    ExportFormatSPSS ExportFormat = "spss"
)

type ExportOptions struct {
    Format            ExportFormat
    IncludeMetadata   bool
    IncludeTimestamps bool
    IncludeIPAddress  bool
    AnonymizeData     bool
    DateFormat        string
    Encoding          string // UTF-8, Shift-JIS, etc.
}

// Export エクスポートのエントリーポイント
func (q *ExportResponsesQuery) Export(
    ctx context.Context,
    surveyID shared.UUID[survey.Survey],
    writer io.Writer,
    options ExportOptions,
) error {
    switch options.Format {
    case ExportFormatCSV:
        return q.ExportAsCSV(ctx, surveyID, writer, options)
    case ExportFormatJSON:
        return q.ExportAsJSON(ctx, surveyID, writer, options)
    case ExportFormatSPSS:
        return q.ExportAsSPSS(ctx, surveyID, writer, options)
    default:
        return fmt.Errorf("unsupported export format: %s", options.Format)
    }
}

// ExportAsCSV CSV形式でのエクスポート
func (q *ExportResponsesQuery) ExportAsCSV(
    ctx context.Context,
    surveyID shared.UUID[survey.Survey],
    writer io.Writer,
    options ExportOptions,
) error {
    survey, err := q.surveyRepo.FindByID(ctx, surveyID)
    if err != nil {
        return err
    }

    responses, err := q.responseRepo.FindBySurveyID(ctx, surveyID)
    if err != nil {
        return err
    }

    // エンコーディング対応
    var csvWriter *csv.Writer
    if options.Encoding == "Shift-JIS" {
        // Shift-JISエンコーダーでラップ
        encoder := japanese.ShiftJIS.NewEncoder()
        writer = transform.NewWriter(writer, encoder)
    }
    csvWriter = csv.NewWriter(writer)
    defer csvWriter.Flush()

    // BOM追加（UTF-8 with BOM）
    if options.Encoding == "UTF-8-BOM" {
        writer.Write([]byte{0xEF, 0xBB, 0xBF})
    }

    // 質問を取得
    questions, err := q.surveyRepo.GetQuestions(ctx, surveyID)
    if err != nil {
        return err
    }

    // ヘッダー行の作成
    headers := q.createCSVHeaders(survey, questions, options)
    if err := csvWriter.Write(headers); err != nil {
        return err
    }

    // 各回答の書き込み
    for _, response := range responses {
        row, err := q.createCSVRow(response, survey, questions, options)
        if err != nil {
            return err
        }
        if err := csvWriter.Write(row); err != nil {
            return err
        }
    }

    return nil
}

// createCSVHeaders CSVヘッダーの作成
func (q *ExportResponsesQuery) createCSVHeaders(
    survey *survey.Survey,
    questions []survey.Question,
    options ExportOptions,
) []string {
    headers := []string{"Response ID"}

    if !options.AnonymizeData {
        headers = append(headers, "Respondent ID")
    }

    if options.IncludeTimestamps {
        headers = append(headers, "Started At", "Submitted At")
    }

    if options.IncludeIPAddress && !options.AnonymizeData {
        headers = append(headers, "IP Address")
    }

    headers = append(headers, "Status", "Completion Rate (%)")

    // 各質問のヘッダー
    for _, question := range questions {
        headers = append(headers, question.Title())

        // メタデータ列
        if options.IncludeMetadata {
            headers = append(headers,
                fmt.Sprintf("%s (Duration)", question.Title()),
                fmt.Sprintf("%s (Changes)", question.Title()))
        }
    }

    return headers
}

// createCSVRow CSV行の作成
func (q *ExportResponsesQuery) createCSVRow(
    response *survey.SurveyResponse,
    survey *survey.Survey,
    questions []survey.Question,
    options ExportOptions,
) ([]string, error) {
    row := []string{response.SurveyResponseID().String()}

    if !options.AnonymizeData {
        if response.ResponderID() != nil {
            row = append(row, response.ResponderID().String())
        } else {
            row = append(row, "")
        }
    }

    if options.IncludeTimestamps {
        row = append(row, q.formatDateTime(response.StartedAt(), options.DateFormat))
        if response.SubmittedAt() != nil {
            row = append(row, q.formatDateTime(*response.SubmittedAt(), options.DateFormat))
        } else {
            row = append(row, "")
        }
    }

    if options.IncludeIPAddress && !options.AnonymizeData {
        row = append(row, response.IPAddress())
    }

    row = append(row, string(response.Status()))
    row = append(row, fmt.Sprintf("%.1f", response.Progress().Percentage))

    // 各質問の回答
    for _, question := range questions {
        answer := response.GetAnswer(question.QuestionID())

        if answer == nil {
            row = append(row, "")
        } else {
            row = append(row, q.formatAnswerValue(answer.Value(), question))
        }

        // メタデータ
        if options.IncludeMetadata {
            if answer != nil && answer.Duration() != nil {
                row = append(row, fmt.Sprintf("%.1f", answer.Duration().Seconds()))
            } else {
                row = append(row, "")
            }

            if answer != nil {
                row = append(row, fmt.Sprintf("%d", len(answer.PreviousValues())))
            } else {
                row = append(row, "0")
            }
        }
    }

    return row, nil
}

// formatAnswerValue 回答値のフォーマット
func (q *ExportResponsesQuery) formatAnswerValue(value survey.AnswerValue, question survey.Question) string {
    switch v := value.(type) {
    case survey.TextAnswer:
        return v.Text
    case survey.SelectionAnswer:
        // 選択肢のラベルを取得
        labels := []string{}
        for _, selectedID := range v.SelectedIDs {
            if selectedID == "other" && v.OtherValue != nil {
                labels = append(labels, fmt.Sprintf("その他: %s", *v.OtherValue))
            } else {
                for _, option := range question.Options() {
                    if option.ID == selectedID {
                        labels = append(labels, option.Label)
                        break
                    }
                }
            }
        }
        return strings.Join(labels, ", ")
    case survey.NumericAnswer:
        return fmt.Sprintf("%g", v.Number)
    case survey.ScaleAnswer:
        return fmt.Sprintf("%d", v.Scale)
    case survey.DateTimeAnswer:
        return v.DateTime.Format("2006-01-02 15:04:05")
    default:
        return value.ToString()
    }
}

// formatDateTime 日時のフォーマット
func (q *ExportResponsesQuery) formatDateTime(t time.Time, format string) string {
    if format == "" {
        format = "2006-01-02 15:04:05"
    }
    return t.Format(format)
}

// ExportAsJSON JSON形式でのエクスポート
func (q *ExportResponsesQuery) ExportAsJSON(
    ctx context.Context,
    surveyID shared.UUID[survey.Survey],
    writer io.Writer,
    options ExportOptions,
) error {
    survey, err := q.surveyRepo.FindByID(ctx, surveyID)
    if err != nil {
        return err
    }

    responses, err := q.responseRepo.FindBySurveyID(ctx, surveyID)
    if err != nil {
        return err
    }

    // エクスポート用の構造体
    type ExportData struct {
        Survey    SurveyInfo     `json:"survey"`
        Responses []ResponseData `json:"responses"`
        Metadata  ExportMetadata `json:"metadata"`
    }

    type SurveyInfo struct {
        ID          string    `json:"id"`
        Title       string    `json:"title"`
        Description string    `json:"description,omitempty"`
        CreatedAt   time.Time `json:"created_at"`
    }

    type ResponseData struct {
        ID           string                 `json:"id"`
        ResponderID  *string                `json:"responder_id,omitempty"`
        StartedAt    time.Time              `json:"started_at"`
        SubmittedAt  *time.Time             `json:"submitted_at,omitempty"`
        Status       string                 `json:"status"`
        Progress     float32                `json:"progress"`
        Answers      map[string]interface{} `json:"answers"`
        Metadata     *ResponseMetadata      `json:"metadata,omitempty"`
    }

    type ResponseMetadata struct {
        IPAddress  string        `json:"ip_address,omitempty"`
        UserAgent  string        `json:"user_agent,omitempty"`
        TotalTime  time.Duration `json:"total_time,omitempty"`
        DeviceType string        `json:"device_type,omitempty"`
    }

    type ExportMetadata struct {
        ExportedAt time.Time     `json:"exported_at"`
        TotalCount int           `json:"total_count"`
        ExportedBy string        `json:"exported_by,omitempty"`
        Options    ExportOptions `json:"options"`
    }

    // データの構築
    description := ""
    if survey.Description() != nil {
        description = *survey.Description()
    }

    exportData := ExportData{
        Survey: SurveyInfo{
            ID:          survey.SurveyID().String(),
            Title:       survey.Title(),
            Description: description,
            CreatedAt:   survey.CreatedAt(),
        },
        Responses: []ResponseData{},
        Metadata: ExportMetadata{
            ExportedAt: time.Now(),
            TotalCount: len(responses),
            Options:    options,
        },
    }

    // 質問を取得
    questions, err := q.surveyRepo.GetQuestions(ctx, surveyID)
    if err != nil {
        return err
    }

    for _, response := range responses {
        responseData := ResponseData{
            ID:          response.SurveyResponseID().String(),
            StartedAt:   response.StartedAt(),
            SubmittedAt: response.SubmittedAt(),
            Status:      string(response.Status()),
            Progress:    response.Progress().Percentage,
            Answers:     make(map[string]interface{}),
        }

        if !options.AnonymizeData && response.ResponderID() != nil {
            responderID := response.ResponderID().String()
            responseData.ResponderID = &responderID
        }

        // 回答データの追加
        for _, question := range questions {
            answer := response.GetAnswer(question.QuestionID())
            if answer != nil {
                responseData.Answers[question.QuestionID().String()] = q.formatAnswerForJSON(answer, question)
            }
        }

        // メタデータの追加
        if options.IncludeMetadata {
            responseData.Metadata = &ResponseMetadata{}
            if response.SubmittedAt() != nil {
                responseData.Metadata.TotalTime = response.SubmittedAt().Sub(response.StartedAt())
            }
            if !options.AnonymizeData {
                responseData.Metadata.IPAddress = response.IPAddress()
                responseData.Metadata.UserAgent = response.UserAgent()
            }
        }

        exportData.Responses = append(exportData.Responses, responseData)
    }

    // JSON出力
    encoder := json.NewEncoder(writer)
    encoder.SetIndent("", "  ")
    return encoder.Encode(exportData)
}

// formatAnswerForJSON JSON用の回答フォーマット
func (q *ExportResponsesQuery) formatAnswerForJSON(answer *survey.Answer, question survey.Question) interface{} {
    type AnswerData struct {
        QuestionTitle string      `json:"question_title"`
        Value         interface{} `json:"value"`
        AnsweredAt    time.Time   `json:"answered_at"`
        Duration      *float64    `json:"duration_seconds,omitempty"`
    }

    data := AnswerData{
        QuestionTitle: question.Title(),
        Value:         answer.Value(),
        AnsweredAt:    answer.AnsweredAt(),
    }

    if answer.Duration() != nil {
        duration := answer.Duration().Seconds()
        data.Duration = &duration
    }

    return data
}

// ExportAsSPSS SPSS形式でのエクスポート（簡易実装）
func (q *ExportResponsesQuery) ExportAsSPSS(
    ctx context.Context,
    surveyID shared.UUID[survey.Survey],
    writer io.Writer,
    options ExportOptions,
) error {
    // SPSS形式への変換は複雑なため、CSV形式で代替
    options.Format = ExportFormatCSV
    return q.ExportAsCSV(ctx, surveyID, writer, options)
}
```

## 7. まとめ

このアンケート回答システムの詳細設計では、以下の要素を網羅しました：

1. **型安全な回答値システム**
   - 各質問タイプに対応した専用の回答型（テキスト、選択式、数値、スケール、日時）
   - 厳密なバリデーション
   - 拡張可能な設計

2. **高度なセッション管理**
   - 自動保存機能
   - セッション復帰
   - タイムアウト管理
   - 並行アクセス制御

3. **包括的な集計機能**
   - リアルタイム統計
   - クロス集計
   - 高度な分析（外れ値検出など）

4. **多様なエクスポート形式**
   - CSV、JSON対応
   - エンコーディング対応
   - メタデータ含有オプション

この設計により、シンプルで直感的なアンケート機能とエンタープライズグレードの信頼性を両立したアンケートシステムを実現できます。プロジェクトの実装パターンに準拠し、Clean Architecture原則に従った設計となっています。
