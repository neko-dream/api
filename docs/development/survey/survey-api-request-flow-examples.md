# アンケートAPI リクエストフロー例

このAPIは次のフローをサポートします：
1. **アンケート作成フロー**: アンケートの作成、質問の追加、公開
2. **回答フロー**: 回答セッションの開始、回答の送信、完了
3. **集計・分析フロー**: 統計データの取得、個別回答の閲覧、エクスポート

## 1. アンケート作成フロー

### 1.1 基本的なアンケート作成

```bash
# 1. ログイン（JWT取得）
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "山田太郎"
  }
}

# 2. アンケート作成
POST /surveys
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "title": "顧客満足度調査",
  "description": "サービスの改善のため、ご意見をお聞かせください",
  "settings": {
    "requireLogin": false,
    "allowMultipleSubmit": false,
    "showProgressBar": true,
    "confirmationMessage": "ご回答ありがとうございました"
  },
  "startAt": "2024-01-01T00:00:00Z",
  "endAt": "2024-12-31T23:59:59Z"
}

# Response
{
  "id": "7f3e4d5c-1234-5678-9abc-def012345678",
  "title": "顧客満足度調査",
  "description": "サービスの改善のため、ご意見をお聞かせください",
  "creatorId": "550e8400-e29b-41d4-a716-446655440000",
  "status": "draft",
  "settings": {
    "requireLogin": false,
    "allowMultipleSubmit": false,
    "showProgressBar": true,
    "confirmationMessage": "ご回答ありがとうございました"
  },
  "startAt": "2024-01-01T00:00:00Z",
  "endAt": "2024-12-31T23:59:59Z",
  "createdAt": "2023-12-15T10:30:00Z",
  "updatedAt": "2023-12-15T10:30:00Z"
}
```

### 1.2 質問の追加

```bash
# 3. テキスト質問の追加
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/questions
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "type": "text",
  "title": "お名前を教えてください",
  "description": "フルネームでご記入ください",
  "required": true,
  "order": 1,
  "validation": {
    "type": "text",
    "parameters": {
      "minLength": 2,
      "maxLength": 50
    }
  }
}

# Response
{
  "id": "q1-550e8400-e29b-41d4-a716-446655440001",
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "type": "text",
  "title": "お名前を教えてください",
  "description": "フルネームでご記入ください",
  "required": true,
  "order": 1,
  "validation": {
    "type": "text",
    "parameters": {
      "minLength": 2,
      "maxLength": 50
    }
  }
}

# 4. 単一選択質問の追加
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/questions
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "type": "radio",
  "title": "サービスの満足度を教えてください",
  "required": true,
  "order": 2,
  "options": [
    {"id": "opt1", "label": "とても満足", "value": "5", "order": 1},
    {"id": "opt2", "label": "満足", "value": "4", "order": 2},
    {"id": "opt3", "label": "普通", "value": "3", "order": 3},
    {"id": "opt4", "label": "不満", "value": "2", "order": 4},
    {"id": "opt5", "label": "とても不満", "value": "1", "order": 5}
  ]
}

# Response
{
  "id": "q2-550e8400-e29b-41d4-a716-446655440002",
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "type": "radio",
  "title": "サービスの満足度を教えてください",
  "required": true,
  "order": 2,
  "options": [
    {"id": "opt1", "label": "とても満足", "value": "5", "order": 1},
    {"id": "opt2", "label": "満足", "value": "4", "order": 2},
    {"id": "opt3", "label": "普通", "value": "3", "order": 3},
    {"id": "opt4", "label": "不満", "value": "2", "order": 4},
    {"id": "opt5", "label": "とても不満", "value": "1", "order": 5}
  ]
}

# 5. テキストエリア質問の追加
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/questions
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "type": "textarea",
  "title": "ご意見・ご感想をお聞かせください",
  "description": "改善のための貴重なご意見として活用させていただきます",
  "required": false,
  "order": 3,
  "validation": {
    "type": "text",
    "parameters": {
      "minLength": 10,
      "maxLength": 1000
    }
  }
}

# 6. 複数選択質問の追加
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/questions
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "type": "checkbox",
  "title": "どの機能を利用していますか？（複数選択可）",
  "required": true,
  "order": 4,
  "options": [
    {"id": "feat1", "label": "基本機能", "value": "basic", "order": 1},
    {"id": "feat2", "label": "レポート機能", "value": "report", "order": 2},
    {"id": "feat3", "label": "API連携", "value": "api", "order": 3},
    {"id": "feat4", "label": "データ分析", "value": "analytics", "order": 4},
    {"id": "other", "label": "その他", "value": "other", "order": 5}
  ],
  "validation": {
    "type": "selection",
    "parameters": {
      "minSelections": 1,
      "maxSelections": 5,
      "allowOther": true
    }
  }
}

# 7. スケール質問の追加
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/questions
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "type": "scale",
  "title": "友人に推薦する可能性はどのくらいですか？",
  "description": "0（全く推薦しない）から10（強く推薦する）でお答えください",
  "required": true,
  "order": 5,
  "options": {
    "min": 0,
    "max": 10,
    "minLabel": "全く推薦しない",
    "maxLabel": "強く推薦する"
  }
}

# 8. アンケートの公開
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/publish
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Response
{
  "id": "7f3e4d5c-1234-5678-9abc-def012345678",
  "status": "published",
  "publishedAt": "2023-12-15T11:00:00Z"
}
```

## 2. 回答フロー

### 2.1 匿名回答フロー

```bash
# 1. アンケート詳細の取得（公開情報）
GET /surveys/7f3e4d5c-1234-5678-9abc-def012345678

# Response
{
  "id": "7f3e4d5c-1234-5678-9abc-def012345678",
  "title": "顧客満足度調査",
  "description": "サービスの改善のため、ご意見をお聞かせください",
  "status": "published",
  "settings": {
    "requireLogin": false,
    "allowMultipleSubmit": false,
    "showProgressBar": true
  },
  "questions": [
    {
      "id": "q1-550e8400-e29b-41d4-a716-446655440001",
      "type": "text",
      "title": "お名前を教えてください",
      "description": "フルネームでご記入ください",
      "required": true,
      "order": 1
    },
    {
      "id": "q2-550e8400-e29b-41d4-a716-446655440002",
      "type": "radio",
      "title": "サービスの満足度を教えてください",
      "required": true,
      "order": 2,
      "options": [
        {"id": "opt1", "label": "とても満足", "value": "5"},
        {"id": "opt2", "label": "満足", "value": "4"},
        {"id": "opt3", "label": "普通", "value": "3"},
        {"id": "opt4", "label": "不満", "value": "2"},
        {"id": "opt5", "label": "とても不満", "value": "1"}
      ]
    }
    // ... 他の質問
  ]
}

# 2. 回答セッションの開始
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/responses
Content-Type: application/json

{
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  "ipAddress": "192.168.1.100"
}

# Response
{
  "id": "resp-123e4567-e89b-12d3-a456-426614174000",
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "sessionToken": "sess_abcdef123456789",
  "status": "in_progress",
  "startedAt": "2023-12-20T14:00:00Z",
  "progress": {
    "totalQuestions": 5,
    "answeredQuestions": 0,
    "percentage": 0
  }
}

Set-Cookie: SurveySession=sess_abcdef123456789

# 3. 回答の送信（1問ずつ）
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ, SurveySession=sess_abcdef123456789

{
  "questionId": "q1-550e8400-e29b-41d4-a716-446655440001",
  "value": {
    "text": "山田太郎"
  }
}

# Response
{
  "id": "ans-111e4567-e89b-12d3-a456-426614174001",
  "responseId": "resp-123e4567-e89b-12d3-a456-426614174000",
  "questionId": "q1-550e8400-e29b-41d4-a716-446655440001",
  "answeredAt": "2023-12-20T14:01:00Z",
  "progress": {
    "totalQuestions": 5,
    "answeredQuestions": 1,
    "percentage": 20
  }
}

# 4. 次の質問の取得
GET /responses/resp-123e4567-e89b-12d3-a456-426614174000/next-question
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "question": {
    "id": "q2-550e8400-e29b-41d4-a716-446655440002",
    "type": "radio",
    "title": "サービスの満足度を教えてください",
    "required": true,
    "options": [
      {"id": "opt1", "label": "とても満足", "value": "5"},
      {"id": "opt2", "label": "満足", "value": "4"},
      {"id": "opt3", "label": "普通", "value": "3"},
      {"id": "opt4", "label": "不満", "value": "2"},
      {"id": "opt5", "label": "とても不満", "value": "1"}
    ]
  },
  "currentPage": 1,
  "totalPages": 3
}

# 5. 選択式回答の送信
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

{
  "questionId": "q2-550e8400-e29b-41d4-a716-446655440002",
  "value": {
    "selected": ["opt4"],
    "otherValue": null
  }
}

# 6. 次の質問の取得
GET /responses/resp-123e4567-e89b-12d3-a456-426614174000/next-question
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "question": {
    "id": "q3-550e8400-e29b-41d4-a716-446655440003",
    "type": "textarea",
    "title": "ご意見・ご感想をお聞かせください",
    "description": "改善のための貴重なご意見として活用させていただきます",
    "required": false
  }
}

# 7. 複数選択の回答（その他含む）
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

{
  "questionId": "q4-550e8400-e29b-41d4-a716-446655440004",
  "value": {
    "selected": ["feat1", "feat2", "other"],
    "otherValue": "カスタマーサポート"
  }
}

# 8. バッチ回答送信（ページ単位）
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers/batch
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

{
  "answers": [
    {
      "questionId": "q5-550e8400-e29b-41d4-a716-446655440005",
      "value": {
        "scale": 3
      }
    },
    {
      "questionId": "q6-550e8400-e29b-41d4-a716-446655440006",
      "value": {
        "datetime": "2023-12-20T15:00:00Z"
      }
    }
  ]
}

# 9. 回答の完了
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/submit
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "id": "resp-123e4567-e89b-12d3-a456-426614174000",
  "status": "submitted",
  "submittedAt": "2023-12-20T14:10:00Z",
  "completionMessage": "ご回答ありがとうございました"
}
```

### 2.2 ログインユーザーの回答フロー

```bash
# 1. ユーザーログイン
POST /auth/login
Content-Type: application/json

{
  "email": "customer@example.com",
  "password": "password456"
}

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-789e0123-e89b-12d3-a456-426614174000"
  }
}

# 2. 回答セッションの開始（認証済み）
POST /surveys/7f3e4d5c-1234-5678-9abc-def012345678/responses
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "id": "resp-456e7890-e89b-12d3-a456-426614174000",
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "responderId": "user-789e0123-e89b-12d3-a456-426614174000",
  "sessionToken": "sess_xyz987654321",
  "status": "in_progress",
  "startedAt": "2023-12-20T15:00:00Z"
}

# 3. 以降の回答フローは匿名回答と同様
```

### 2.3 自動保存とセッション復帰

**実装方針**: WebSocketを使用せず、HTTPベースのセッション管理を実装しています。

#### 2.3.1 単一質問の自動保存

```bash
# 1. 途中で回答を自動保存（単一質問）
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/autosave
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789
Content-Type: application/json

{
  "currentQuestionId": "q3-550e8400-e29b-41d4-a716-446655440003",
  "draftAnswer": {
    "text": "レスポンス速度が遅い時があり..."
  },
  "lastActivity": "2023-12-20T14:05:00Z"
}

# Response (200 OK)
{
  "savedAt": "2023-12-20T14:05:30Z",
  "status": "saved",
  "sessionExtended": true,
  "expiresAt": "2023-12-21T14:05:30Z"
}

# Response (セッション期限切れ: 401 Unauthorized)
{
  "error": {
    "code": "SESSION_EXPIRED",
    "message": "セッションの有効期限が切れました"
  }
}

```

#### 2.3.2 複数回答のバッチ自動保存

```bash
# 複数の回答をまとめて自動保存
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/autosave/batch
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789
Content-Type: application/json

{
  "answers": [
    {
      "questionId": "q3-550e8400-e29b-41d4-a716-446655440003",
      "value": {"text": "レスポンス速度が遅い時があり、特にピーク時間帯です"}
    },
    {
      "questionId": "q4-550e8400-e29b-41d4-a716-446655440004",
      "value": {"selected": ["feat1", "feat2"]}
    }
  ],
  "currentPage": 2,
  "lastActivity": "2023-12-20T14:06:00Z",
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
}

# Response (200 OK)
{
  "savedAt": "2023-12-20T14:06:30Z",
  "savedCount": 2,
  "totalAutoSaves": 5,
  "sessionExtended": true,
  "expiresAt": "2023-12-21T14:06:30Z"
}
```

#### 2.3.3 セッションの復帰

```bash
# セッション復帰（セッショントークンを使用）
GET /responses/session/sess_abcdef123456789
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response (200 OK)
{
  "responseId": "resp-123e4567-e89b-12d3-a456-426614174000",
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "status": "in_progress",
  "currentQuestionId": "q3-550e8400-e29b-41d4-a716-446655440003",
  "currentPage": 2,
  "progress": {
    "totalQuestions": 5,
    "answeredQuestions": 2,
    "percentage": 40
  },
  "answers": [
    {
      "questionId": "q1-550e8400-e29b-41d4-a716-446655440001",
      "value": {"text": "山田太郎"}
    },
    {
      "questionId": "q2-550e8400-e29b-41d4-a716-446655440002",
      "value": {"selected": ["opt4"]}
    }
  ],
  "draftAnswer": {
    "questionId": "q3-550e8400-e29b-41d4-a716-446655440003",
    "value": {"text": "レスポンス速度が遅い時があり..."}
  },
  "metadata": {
    "lastSavedAt": "2023-12-20T14:06:30Z",
    "totalAutoSaves": 5,
    "showInactivityWarning": false,
    "sessionExpiresAt": "2023-12-21T14:06:30Z"
  }
}

# Response (非アクティブ警告あり: 200 OK)
{
  "responseId": "resp-123e4567-e89b-12d3-a456-426614174000",
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "status": "in_progress",
  "metadata": {
    "showInactivityWarning": true,
    "inactivityDuration": "PT35M",  // ISO 8601 duration (35分)
    "lastActivityAt": "2023-12-20T13:30:00Z"
  }
  // ... 他のフィールド
}

# Response (セッションが見つからない: 404 Not Found)
{
  "error": {
    "code": "SESSION_NOT_FOUND",
    "message": "セッションが見つかりません"
  }
}
```

#### 2.3.4 セッション状態の確認

```bash
# セッション状態のチェック（軽量なエンドポイント）
GET /responses/session/sess_abcdef123456789/status
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response (200 OK)
{
  "sessionActive": true,
  "expiresAt": "2023-12-21T14:06:30Z",
  "lastActivityAt": "2023-12-20T14:06:30Z",
  "showInactivityWarning": false,
  "progress": {
    "percentage": 40
  }
}

# Response (セッション期限切れ: 401 Unauthorized)
{
  "sessionActive": false,
  "reason": "expired"
}
```

## 3. 集計・分析フロー

### 3.1 基本的な集計結果の取得

```bash
# 1. アンケート全体の集計
GET /surveys/7f3e4d5c-1234-5678-9abc-def012345678/results
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "totalResponses": 150,
  "completedResponses": 120,
  "inProgressResponses": 25,
  "abandonedResponses": 5,
  "completionRate": 80.0,
  "averageCompletionTime": "PT5M30S",
  "questionResults": [
    {
      "questionId": "q2-550e8400-e29b-41d4-a716-446655440002",
      "questionTitle": "サービスの満足度を教えてください",
      "questionType": "radio",
      "totalAnswers": 120,
      "statistics": {
        "options": [
          {"optionId": "opt1", "optionLabel": "とても満足", "count": 30, "percentage": 25.0},
          {"optionId": "opt2", "optionLabel": "満足", "count": 45, "percentage": 37.5},
          {"optionId": "opt3", "optionLabel": "普通", "count": 25, "percentage": 20.8},
          {"optionId": "opt4", "optionLabel": "不満", "count": 15, "percentage": 12.5},
          {"optionId": "opt5", "optionLabel": "とても不満", "count": 5, "percentage": 4.2}
        ],
        "mode": ["opt2"],
        "responseRate": 100.0
      }
    },
    {
      "questionId": "q5-550e8400-e29b-41d4-a716-446655440005",
      "questionTitle": "友人に推薦する可能性はどのくらいですか？",
      "questionType": "scale",
      "totalAnswers": 120,
      "statistics": {
        "min": 0,
        "max": 10,
        "mean": 7.2,
        "median": 8,
        "mode": [8],
        "standardDeviation": 2.1,
        "percentiles": {
          "25": 6,
          "50": 8,
          "75": 9,
          "90": 10,
          "95": 10
        },
        "netPromoterScore": 35.5
      }
    }
  ],
  "generatedAt": "2023-12-21T10:00:00Z"
}

# 2. 特定質問の詳細統計
GET /surveys/7f3e4d5c-1234-5678-9abc-def012345678/questions/q4-550e8400-e29b-41d4-a716-446655440004/statistics
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response（複数選択質問）
{
  "questionId": "q4-550e8400-e29b-41d4-a716-446655440004",
  "questionTitle": "どの機能を利用していますか？",
  "totalResponses": 120,
  "statistics": {
    "options": [
      {"optionId": "feat1", "label": "基本機能", "count": 110, "percentage": 91.7},
      {"optionId": "feat2", "label": "レポート機能", "count": 85, "percentage": 70.8},
      {"optionId": "feat3", "label": "API連携", "count": 40, "percentage": 33.3},
      {"optionId": "feat4", "label": "データ分析", "count": 65, "percentage": 54.2},
      {"optionId": "other", "label": "その他", "count": 15, "percentage": 12.5}
    ],
    "otherValues": [
      "カスタマーサポート",
      "モバイルアプリ",
      "Webhook連携"
    ],
    "averageSelectionsPerResponse": 2.6,
    "correlationMatrix": {
      "feat1-feat2": 0.75,
      "feat2-feat4": 0.82,
      "feat3-feat4": 0.65
    }
  }
}
```

### 3.2 個別回答の閲覧

```bash
# 1. 回答一覧の取得（ページネーション付き）
GET /surveys/7f3e4d5c-1234-5678-9abc-def012345678/responses?offset=0&limit=20&status=submitted
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "items": [
    {
      "response": {
        "id": "resp-123e4567-e89b-12d3-a456-426614174000",
        "responderId": null,
        "startedAt": "2023-12-20T14:00:00Z",
        "submittedAt": "2023-12-20T14:10:00Z",
        "status": "submitted"
      },
      "answers": [
        {
          "questionId": "q1-550e8400-e29b-41d4-a716-446655440001",
          "questionTitle": "お名前を教えてください",
          "answer": {
            "value": {"text": "山田太郎"},
            "answeredAt": "2023-12-20T14:01:00Z"
          }
        },
        {
          "questionId": "q2-550e8400-e29b-41d4-a716-446655440002",
          "questionTitle": "サービスの満足度を教えてください",
          "answer": {
            "value": {"selected": ["opt4"]},
            "answeredAt": "2023-12-20T14:02:00Z"
          }
        }
        // ... 他の回答
      ]
    }
    // ... 他の回答セッション
  ],
  "totalCount": 120,
  "offset": 0,
  "limit": 20,
  "hasMore": true
}

# 2. 特定の回答詳細
GET /responses/resp-123e4567-e89b-12d3-a456-426614174000
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "id": "resp-123e4567-e89b-12d3-a456-426614174000",
  "surveyId": "7f3e4d5c-1234-5678-9abc-def012345678",
  "surveyTitle": "顧客満足度調査",
  "responderId": null,
  "startedAt": "2023-12-20T14:00:00Z",
  "submittedAt": "2023-12-20T14:10:00Z",
  "completionTime": "PT10M",
  "status": "submitted",
  "ipAddress": "192.168.1.100",
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  "answers": [
    {
      "questionId": "q1-550e8400-e29b-41d4-a716-446655440001",
      "questionTitle": "お名前を教えてください",
      "questionType": "text",
      "answer": {
        "id": "ans-111e4567-e89b-12d3-a456-426614174001",
        "value": {"text": "山田太郎"},
        "answeredAt": "2023-12-20T14:01:00Z",
        "duration": 30.5,
        "changeCount": 0
      }
    }
    // ... 全ての回答
  ]
}
```

### 3.3 エクスポート

```bash
# 1. CSV形式でエクスポート
GET /surveys/7f3e4d5c-1234-5678-9abc-def012345678/export?format=csv&encoding=utf-8
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response Headers
Content-Type: text/csv; charset=utf-8
Content-Disposition: attachment; filename="survey_7f3e4d5c_responses_20231221.csv"

# Response Body
Response ID,Respondent ID,Started At,Submitted At,Status,Completion Rate (%),お名前を教えてください,サービスの満足度を教えてください,不満な点を具体的に教えてください,どの機能を利用していますか？,友人に推薦する可能性はどのくらいですか？
resp-123e4567-e89b-12d3-a456-426614174000,,2023-12-20T14:00:00Z,2023-12-20T14:10:00Z,submitted,100.0,山田太郎,不満,レスポンス速度が遅い時があり改善を希望します,"基本機能, レポート機能, その他: カスタマーサポート",3
...

# 2. JSON形式でエクスポート
GET /surveys/7f3e4d5c-1234-5678-9abc-def012345678/export?format=json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response
{
  "survey": {
    "id": "7f3e4d5c-1234-5678-9abc-def012345678",
    "title": "顧客満足度調査",
    "description": "サービスの改善のため、ご意見をお聞かせください",
    "createdAt": "2023-12-15T10:30:00Z"
  },
  "responses": [
    {
      "id": "resp-123e4567-e89b-12d3-a456-426614174000",
      "responderId": null,
      "startedAt": "2023-12-20T14:00:00Z",
      "submittedAt": "2023-12-20T14:10:00Z",
      "status": "submitted",
      "progress": 100.0,
      "answers": {
        "q1-550e8400-e29b-41d4-a716-446655440001": {
          "questionTitle": "お名前を教えてください",
          "value": {"text": "山田太郎"}
        },
        "q2-550e8400-e29b-41d4-a716-446655440002": {
          "questionTitle": "サービスの満足度を教えてください",
          "value": {"selected": ["opt4"], "label": "不満"}
        },
        "q3-550e8400-e29b-41d4-a716-446655440003": {
          "questionTitle": "不満な点を具体的に教えてください",
          "value": {"text": "レスポンス速度が遅い時があり改善を希望します"}
        },
        "q4-550e8400-e29b-41d4-a716-446655440004": {
          "questionTitle": "どの機能を利用していますか？",
          "value": {
            "selected": ["feat1", "feat2", "other"],
            "labels": ["基本機能", "レポート機能", "その他"],
            "otherValue": "カスタマーサポート"
          }
        },
        "q5-550e8400-e29b-41d4-a716-446655440005": {
          "questionTitle": "友人に推薦する可能性はどのくらいですか？",
          "value": {"scale": 3}
        }
      }
    }
    // ... 他の回答
  ],
  "metadata": {
    "exportedAt": "2023-12-21T10:30:00Z",
    "totalCount": 120,
    "options": {
      "format": "json",
      "includeMetadata": true,
      "includeTimestamps": true
    }
  }
}
```

## 4. エラーハンドリング

### 4.1 バリデーションエラー

```bash
# 必須項目が未入力
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

{
  "questionId": "q1-550e8400-e29b-41d4-a716-446655440001",
  "value": {
    "text": ""
  }
}

# Response (400 Bad Request)
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "回答は必須です",
    "field": "value.text",
    "questionId": "q1-550e8400-e29b-41d4-a716-446655440001"
  }
}

# 文字数制限違反
{
  "questionId": "q1-550e8400-e29b-41d4-a716-446655440001",
  "value": {
    "text": "あ"
  }
}

# Response (400 Bad Request)
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "最小文字数は2文字です",
    "field": "value.text",
    "minLength": 2,
    "currentLength": 1
  }
}
```

### 4.2 認証エラー

```bash
# トークン期限切れ
GET /surveys/7f3e4d5c-1234-5678-9abc-def012345678/results
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response (401 Unauthorized)
{
  "error": {
    "code": "TOKEN_EXPIRED",
    "message": "認証トークンの有効期限が切れています"
  }
}
```

### 4.3 権限エラー

```bash
# 他人のアンケート結果へのアクセス
GET /surveys/other-user-survey-id/results
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response (403 Forbidden)
{
  "error": {
    "code": "PERMISSION_DENIED",
    "message": "このアンケートの結果を閲覧する権限がありません"
  }
}
```

### 4.4 セッションエラー

```bash
# セッションタイムアウト
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response (401 Unauthorized)
{
  "error": {
    "code": "SESSION_EXPIRED",
    "message": "セッションの有効期限が切れました。最初からやり直してください。"
  }
}

# 既に提出済みの回答への追加
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response (400 Bad Request)
{
  "error": {
    "code": "RESPONSE_ALREADY_SUBMITTED",
    "message": "この回答は既に提出されています"
  }
}
```

## 5. 組織限定アンケート

```bash
# 1. 組織限定アンケートの作成
POST /surveys
Content-Type: application/json
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

{
  "title": "社内満足度調査",
  "description": "働きやすさに関するアンケート",
  "organizationId": "org-550e8400-e29b-41d4-a716-446655440000",
  "settings": {
    "requireLogin": true,
    "allowMultipleSubmit": false,
    "restrictToOrganization": true
  }
}

# 2. 組織外ユーザーのアクセス試行
GET /surveys/7f3e4d5c-organization-survey-id

# Response (403 Forbidden)
{
  "error": {
    "code": "ORGANIZATION_RESTRICTED",
    "message": "このアンケートは組織メンバーのみ回答可能です"
  }
}
```


## 7. レート制限とクォータ

```bash
# レート制限に達した場合
POST /responses/resp-123e4567-e89b-12d3-a456-426614174000/answers
Cookie: SessionId=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..., SurveySession=sess_abcdef123456789

# Response (429 Too Many Requests)
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "リクエストが多すぎます。しばらく待ってから再試行してください。",
    "retryAfter": 60
  }
}

# Response Headers
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1703145600
Retry-After: 60
```
