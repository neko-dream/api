# プレゼンテーション層（Presentation Layer）

## 概要

HTTPリクエスト/レスポンスの処理を担当する最外層です。OpenAPIから自動生成されたコードを基盤とし、型安全なAPI実装を提供します。

## ディレクトリ構造

```
presentation/
├── handler/         # HTTPハンドラー実装
│   ├── auth_handler.go
│   ├── opinion_handler.go
│   ├── talk_session_handler.go
│   ├── error_handler.go      # エラーハンドリング
│   └── session_helper.go     # 認証ヘルパー
└── oas/            # OpenAPI自動生成コード
```

## 実装フロー

1. TypeSpec定義 (`tsp/`) → OpenAPIスペック生成
2. OpenAPIスペック → ogenによるGoコード生成 (`presentation/oas/`)
3. 自動生成されたインターフェースを`handler/`で実装
4. ビジネスロジックはApplication層に委譲

## 認証パターン

TypeSpecでの定義により3つのパターンがあります：

### 1. 認証必須（デフォルト）
```typescript
@route("/opinions/{opinionId}/report")
@post
op reportOpinion(...): void;
```
- `@useAuth`を付けない
- 認証がない場合は401エラー

### 2. 認証不要
```typescript
@route("/test")
@useAuth([])
op test(): TestResponse;
```
- 公開API
- 誰でもアクセス可能

### 3. 認証オプショナル
```typescript
@route("/talk-sessions/{id}/opinions")
@useAuth(OptionalCookieAuth)
op opinionComments2(...): OpinionCommentsResponse;
```
- 認証があればユーザー情報を取得
- 認証がなくてもアクセス可能

## 設計原則

- **薄い層**: HTTP固有の処理のみ実装
- **OpenAPI駆動**: 仕様が単一の真実の源
- **認証のデフォルト**: 明示的に指定しない限り認証必須
- **型安全**: 自動生成された型を活用

## アンチパターン

- ハンドラー内でのビジネスロジック実装
- 手動でのリクエスト/レスポンス型定義
- 直接的なDB操作