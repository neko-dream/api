# インフラストラクチャ層（Infrastructure Layer）

## 概要

技術的な実装を提供し、ドメイン層のインターフェースを実装する層です。外部システムとの連携、データ永続化、認証・認可などを担当します。

## ディレクトリ構造

```
infrastructure/
├── auth/            # 認証・認可
│   ├── jwt/         # JWT処理
│   ├── oauth/       # OAuth認証（Google、LINE）
│   └── session/     # セッション管理
├── config/          # 設定管理
├── crypto/          # 暗号化処理
├── di/              # 依存性注入（dig）
├── email/           # メール送信
│   └── template/    # メールテンプレート
├── external/        # 外部API連携
│   ├── analysis/    # 分析API
│   └── aws/         # AWSサービス
│       └── ses/     # メール送信
├── http/            # HTTPサーバー
│   ├── cookie/      # Cookie管理
│   └── middleware/  # ミドルウェア
├── persistence/     # データ永続化
│   ├── db/          # DB管理・マイグレーション
│   ├── postgresql/  # PostgreSQL接続
│   ├── query/       # クエリ実装
│   ├── repository/  # リポジトリ実装
│   └── sqlc/        # SQLコード生成
│       ├── generated/
│       └── queries/
├── service/         # インフラサービス
└── telemetry/       # 監視・トレーシング
```

## 主要コンポーネント

### 1. データ永続化（persistence/）

- **SQLC**: SQLから型安全なGoコードを自動生成(./scripts/gen.sh)
- **リポジトリ**: ドメイン層のインターフェースを実装
- **クエリ**: 読み取り専用の最適化されたクエリ

### 2. 認証・認可（auth/）

- **JWT**: アクセストークンの生成・検証
- **OAuth**: Google、LINEによるソーシャルログイン
- **セッション**: Cookieベースのセッション管理

### 3. 外部サービス連携（external/）

- **分析API**: 外部の分析サービスとの連携
- **AWS**: S3（ファイルストレージ）、SES（メール送信）

### 4. 暗号化（crypto/）

- ユーザーの個人情報を暗号化
- AES-GCM/CBCによる暗号化実装
- バージョン管理された暗号化方式

### 5. 依存性注入（di/）

- digライブラリを使用したDIコンテナ
- 各層の依存関係を一元管理

## 技術スタック

- **データベース**: PostgreSQL
- **コード生成**: sqlc
- **認証**: JWT、OAuth2.0
- **クラウド**: AWS（S3、SES）
- **監視**: OpenTelemetry、Baselime、Sentry
- **DI**: dig

## 設計原則

- インターフェースの実装に専念
- 技術的詳細の隠蔽
- 設定の外部化（環境変数）
- エラーの適切な変換とログ出力

## アンチパターン

- ビジネスロジックの実装
- ドメイン知識の流出
- 設定のハードコーディング
