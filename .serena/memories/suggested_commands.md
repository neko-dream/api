# 推奨コマンド一覧

## 開発コマンド
- `go generate ./...` - コード生成（TypeSpec/SQL変更後は必須）
- `air` - ホットリロードでサーバー起動
- `go test ./...` - テスト実行
- `golangci-lint run` - リンター実行
- `go fmt ./...` - コードフォーマット

## データベース関連
- `docker compose up -d db` - ローカルDB起動
- `./scripts/migrate-create.sh` - 新規マイグレーション作成

## システムコマンド（Darwin）
- `git` - バージョン管理
- `ls` - ファイルリスト表示
- `cd` - ディレクトリ移動
- `grep` - ファイル内検索（ripgrepも利用可能）
- `find` - ファイル検索

## 環境セットアップ
- `mise install` - 必要なツールのインストール
- `cp .env.example .env` - 環境変数設定
