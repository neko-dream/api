# プロジェクトドキュメントガイド

このファイルは、開発時にどのドキュメントを参照すべきかを示すインデックスです。

## アーキテクチャを理解したい

[アーキテクチャ概要](docs/development/architecture/architecture.md)を参照してください。

### 各層の詳細を知りたい

- **ドメイン層**: [domain-layer.md](docs/development/architecture/layers/domain-layer.md)
- **アプリケーション層**: [application-layer.md](docs/development/architecture/layers/application-layer.md)
- **インフラストラクチャ層**: [infrastructure-layer.md](docs/development/architecture/layers/infrastructure-layer.md)
- **プレゼンテーション層**: [presentation-layer.md](docs/development/architecture/layers/presentation-layer.md)

## 新機能を実装したい

[開発ガイドライン](docs/development/guidelines/development-guideline.md)を参照してください。
実装手順、命名規則、コード生成の使い方が記載されています。

## TDDで開発したい

[TDD開発ガイドライン](docs/development/guidelines/tdd-guideline.md)を参照してください。
Kent BeckのTDDとTidy Firstアプローチに基づいた開発方法が記載されています。

## APIを追加・変更したい

1. `tsp/`ディレクトリでTypeSpec定義を編集
2. [TypeSpec命名規則](docs/development/guidelines/typespec-naming-conventions.md)を参照
3. [開発ガイドライン](docs/development/guidelines/development-guideline.md)の「API定義」セクションを参照

## 認証パターンを理解したい

[プレゼンテーション層](docs/development/architecture/layers/presentation-layer.md)の「認証パターン」セクションを参照してください。

## コマンド一覧

```bash
# コード生成
go generate ./...

# テスト実行
go test ./...

# リンター
golangci-lint run

# フォーマット
go fmt ./...
```

## 重要な注意事項

- **コード生成を忘れずに**: TypeSpecやSQLを変更したら必ず`go generate ./...`を実行
- **テストファースト**: 新機能は必ずテストから書き始める
- **層の依存関係を守る**: presentation → application → domain ← infrastructure
