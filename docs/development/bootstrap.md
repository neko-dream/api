# Bootstrapパッケージ

## 概要

`cmd/server/bootstrap`パッケージは、サーバー起動時の初期化処理を管理する。

## いつ触るか

- 新しいミドルウェアを追加するとき → `routes.go`
- HTTPサーバーの設定を変更するとき → `server.go`
- 起動時の初期化処理を追加するとき → `bootstrap.go`
- Swagger UIの挙動を変更するとき → `swagger.go`

## 構成

### bootstrap.go

アプリケーション起動の制御。DIコンテナから依存を取得し、以下の順で処理を実行:

1. マイグレーション実行
2. HTTPサーバー起動

### server.go

HTTPサーバーの設定。以下の環境変数でタイムアウトを制御可能:

- `HTTP_READ_TIMEOUT` (デフォルト: 15秒)
- `HTTP_WRITE_TIMEOUT` (デフォルト: 15秒)
- `HTTP_IDLE_TIMEOUT` (デフォルト: 60秒)

### routes.go

ルーティング定義。主なエンドポイント:

- `/api/*` - API (認証、CORS、エラーハンドリング付き)
- `/static/*` - 静的ファイル
- `/admin/*` - 管理画面
- `/docs/*` - Swagger UI

新しいエンドポイントやミドルウェアを追加する場合はここを編集。

### swagger.go

Swagger UIの設定。環境ごとにOpenAPI定義のURLを切り替え:

- DEV: 特定ドメインからのアクセス時はそのドメインを使用
- PROD: 固定ドメイン
- LOCAL: localhost

## 注意点

- 起動順序に依存関係がある場合は`bootstrap.go`の`Run()`メソッドで制御
- グローバルなミドルウェアは`routes.go`の`setupRoutes()`で設定
- 環境固有の設定は環境変数で管理し、ハードコードは避ける
