# ドメイン層（Domain Layer）

## 概要

ビジネスロジックの中核を担う層です。技術的な詳細から独立し、純粋なビジネスルールを表現します。

## ディレクトリ構造

```
domain/
├── messages/      # エラーメッセージ定義
├── model/         # ドメインモデル
│   ├── analysis/
│   ├── auth/
│   ├── opinion/
│   ├── organization/
│   ├── session/
│   ├── shared/    # 共通の値オブジェクト
│   ├── survey/
│   ├── talksession/
│   ├── user/
│   └── vote/
└── service/       # ドメインサービス
    ├── organization/
    └── timeline/
```

## 主要コンポーネント

### 1. エンティティ（model/）

IDによる同一性を持つオブジェクト。

**実装パターン:**
- プライベートフィールドによる内部状態の隠蔽
- `NewXxx()`コンストラクタによる生成
- ゲッターメソッドによる読み取り専用アクセス
- ビジネスロジックを含むメソッド

**リポジトリインターフェース:**
- 各エンティティファイル内に定義
- 例：`user.go`内に`UserRepository`インターフェース

### 2. 値オブジェクト

**型安全なID:**
- `shared.UUID[T]` - ジェネリクスを使用した型安全なUUID
- 例：`UUID[User]`、`UUID[TalkSession]`

**ドメイン固有の値:**
- `UserName`、`UserSubject` - ユーザー関連
- `DateOfBirth` - バリデーション付き生年月日
- `Location` - 位置情報
- `RestrictionAttribute` - 参加制限

### 3. ドメインサービス（service/）

エンティティに属さないビジネスロジック。

**主要サービス:**
- `UserService` - ユーザーID重複チェック
- `OrganizationService` - 組織の作成・管理
- `OpinionService` - 意見への投票状態確認
- `AuthService`、`SessionService` - 認証・セッション管理

### 4. エラー定義（messages/）

統一的なエラー表現。

- `APIError` - HTTPステータスコード、エラーコード、メッセージを含む
- 事前定義されたエラー（Unauthorized、Forbidden、NotFoundなど）

## 設計原則

- **技術からの独立**: DB、HTTPなどの技術的詳細を含まない
- **ビジネスロジックの集約**: エンティティ内にロジックを実装
- **型安全性**: ジェネリクスを活用した型安全な実装
- **バリデーション**: コンストラクタやセッターでの検証

## アンチパターン

- 貧血ドメインモデル（ロジックのないデータ構造）
- 技術的関心事の混入
- 過度に細かい値オブジェクトの作成