# アプリケーション層（Application Layer）

## 概要

アプリケーション層は、ユースケースを実装し、ドメイン層のオブジェクトを協調させる層です。ビジネスロジック（ユースケース固有の処理フロー）を実装します。

## ディレクトリ構造

```
application/
├── usecase/          # 書き込み系処理（コマンド）
│   ├── analysis_usecase/
│   ├── auth_usecase/
│   ├── opinion_usecase/
│   ├── organization_usecase/
│   ├── policy_usecase/
│   ├── report_usecase/
│   ├── survey_usecase/
│   ├── talksession_usecase/
│   ├── timeline_usecase/
│   ├── user_usecase/
│   └── vote_usecase/
└── query/           # 読み取り専用処理
    ├── dto/         # データ転送オブジェクト
    ├── analysis_query/
    ├── opinion/
    ├── organization_query/
    ├── policy_query/
    ├── report_query/
    ├── survey_query/
    ├── talksession/
    ├── timeline_query/
    └── user/
```

## ビジネスロジックとは

### 定義
ユースケースを実現するための処理フロー、調整ロジック

### 特徴
- 複数のドメインオブジェクトを協調させる
- 外部サービスとの連携を含む
- トランザクション境界の管理
- ユースケース固有の処理

### 実装例

```go
// application/usecase/user_usecase/register.go
type RegisterUserUseCase struct {
    userRepo     domain.UserRepository
    emailService EmailService
    statsService StatsService
}

func (u *RegisterUserUseCase) Execute(ctx context.Context, input RegisterInput) error {
    // 1. ドメインロジックでユーザー作成
    user := domain.NewUser(input.Name, input.Email, input.Age)

    // 2. ドメインロジックで登録可能かチェック
    if !user.CanRegister() {
        return errors.New("登録条件を満たしていません")
    }

    // 3. 重複チェック（ビジネスロジック）
    exists, err := u.userRepo.ExistsByEmail(user.Email)
    if exists {
        return errors.New("既に登録されています")
    }

    // 4. 永続化
    if err := u.userRepo.Save(ctx, user); err != nil {
        return err
    }

    // 5. メール送信（ビジネスロジック）
    u.emailService.SendWelcomeMail(user.Email)

    // 6. 統計情報更新（ビジネスロジック）
    u.statsService.IncrementUserCount()

    return nil
}
```

## 主要コンポーネント

### 1. UseCase（ユースケース）

書き込み系の処理を実装。

**命名規則:**
- ディレクトリ: `[機能名]_usecase/`
- ファイル: `[アクション].go` (create.go, update.go, delete.go)

**実装方針:**
- 1ユースケース = 1ファイル
- 単一責任の原則に従う
- トランザクション境界を明確にする

### 2. Query（クエリ）

読み取り専用の処理を実装。

**命名規則:**
- ディレクトリ: `[機能名]_query/`
- ファイル: `[詳細処理名]_query.go`

**実装方針:**
- パフォーマンスを重視した実装
- 複数集約をまたぐ読み取りも可能
- ドメインオブジェクトではなくDTOを返す

### 3. DTO（Data Transfer Object）

層間のデータ転送用オブジェクト。

**実装例:**
```go
// application/query/dto/user_dto.go
type UserDTO struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    CreatedAt   time.Time `json:"created_at"`
}
```

## CQRS（Command Query Responsibility Segregation）

### コマンド側（UseCase）
- ドメインモデルを使用
- ビジネスルールの実行
- 状態の変更

### クエリ側（Query）
- 最適化されたクエリ
- DTOの使用
- 読み取り専用

## 設計原則

### 1. 薄いアプリケーション層
- ドメインロジックは含まない
- 処理の流れの制御に専念
- ドメインオブジェクトの協調

### 2. トランザクション管理
- ユースケースがトランザクション境界
- 適切なロールバック処理
- 一貫性の保証

### 3. 外部サービスとの連携
- インターフェース経由での依存
- エラーハンドリング
- リトライ処理

## アンチパターン

### 避けるべき実装
- ドメインロジックの実装
- 巨大なユースケース
- 複数の責任を持つユースケース
- 直接的な技術依存

### 推奨事項
- 単一責任の原則
- 明確なインターフェース定義
- 適切なエラーハンドリング
- テストしやすい設計
