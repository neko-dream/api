# 開発ガイドライン

## 新機能追加の手順

新機能を追加する際は、以下の順序で実装を進めてください。

### 1. ドメイン層の実装
1. `domain/model/`に新しいエンティティ・値オブジェクトを定義
2. 必要に応じて`domain/service/`にドメインサービスを追加
3. リポジトリインターフェースを同じファイル内に定義

### 2. API定義
1. `tsp/`にTypeSpecでAPIエンドポイントを定義
   - [TypeSpec命名規則](./typespec-naming-conventions.md)を参照
2. 認証パターンを明確に指定
   - デフォルト（認証必須）：`@useAuth`なし
   - 認証不要：`@useAuth([])`
   - 認証オプショナル：`@useAuth(OptionalCookieAuth)`

### 3. コード生成（1回目）
```bash
./scripts/gen.sh
```
- OpenAPIスペックが生成される
- `presentation/oas/`にインターフェースが生成される

### 4. アプリケーション層の実装
1. `application/usecase/[機能名]_usecase/`に書き込み処理を実装
2. `application/query/[機能名]_query/`に読み取り処理を実装
3. 必要に応じて`application/query/dto/`にDTOを定義

### 5. プレゼンテーション層の実装
1. `presentation/handler/`に自動生成されたインターフェースの実装を追加
2. ビジネスロジックはApplication層に委譲
3. エラーハンドリングは`error_handler.go`の仕組みを活用

### 6. データベーススキーマ定義
1. `migrations/`に新しいマイグレーションファイルを作成
2. テーブル定義、インデックス、制約を記述

### 7. SQLクエリ定義
1. `infrastructure/persistence/sqlc/queries/`にSQLクエリを記述
2. CRUDクエリとカスタムクエリを定義

### 8. コード生成（2回目）
```bash
./scripts/gen.sh
```
- SQLCによりクエリからGoコードが生成される

### 9. インフラストラクチャ層の実装
1. `infrastructure/persistence/repository/`にリポジトリ実装を追加
2. `infrastructure/persistence/query/`にクエリ実装を追加
3. 必要に応じて`infrastructure/di/`に依存関係を登録

## データフロー

### 読み取り処理（Query）
```
HTTP Request → presentation → application/query → infrastructure → Database
```

### 書き込み処理（Command）
```
HTTP Request → presentation → application/usecase → domain → infrastructure → Database
```

## 重要な設計方針

### CQRS（Command Query Responsibility Segregation）
- **Command（書き込み）**: `application/usecase/`で実装
- **Query（読み取り）**: `application/query/`で実装
- 読み取りは最適化されたクエリを直接実行可能

### リポジトリパターン
- ドメイン層でインターフェース定義
- インフラ層で具体的実装
- ドメインオブジェクトの永続化に特化

### 依存性注入
- `infrastructure/di/`でdigを使用
- インターフェースと実装の分離
- テスタビリティの向上

### コード生成の活用
- **TypeSpec → OpenAPI → ogen**: API定義からハンドラーインターフェース生成
- **SQLC**: SQLから型安全なGoコード生成
- 手動実装の最小化

## 命名規則

### パッケージ名
- 小文字、単数形
- 例：`user_usecase`、`organization_query`

### ファイル名
- snake_case
- 機能を明確に表す名前
- 例：`create_user.go`、`list_organizations_query.go`

### 構造体名
- PascalCase
- 例：`User`、`TalkSession`

### インターフェース名
- PascalCase + 動詞/名詞
- 例：`UserRepository`、`AuthService`

## テスト方針

- 単体テスト：各層で実装
- 統合テスト：`test/txtest/`を使用
- モック：インターフェースベースで作成

## 避けるべきアンチパターン

### ドメイン層
- 技術的な詳細の混入
- 貧血ドメインモデル
- 過度な値オブジェクトの使用

### アプリケーション層
- ドメインロジックの実装
- 巨大なユースケース
- 複数の責任を持つユースケース

### インフラストラクチャ層
- ビジネスロジックの実装
- ドメイン知識の流出
- 設定のハードコーディング

### プレゼンテーション層
- ビジネスロジックの実装
- 手動でのリクエスト/レスポンス型定義
- 直接的なDB操作