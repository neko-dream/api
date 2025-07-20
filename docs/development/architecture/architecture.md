# プロジェクトアーキテクチャ

本プロジェクトはDDD（ドメイン駆動設計）とクリーンアーキテクチャの原則に基づいた4層構造を採用しています。

## ディレクトリ構造

```
internal/
├── application/     # アプリケーション層
├── domain/         # ドメイン層
├── infrastructure/ # インフラストラクチャ層
├── presentation/   # プレゼンテーション層
└── test/          # テストユーティリティ
```

## 依存関係

```
presentation → application → domain ← infrastructure
```

## 開発ガイドライン

[開発ガイドライン](../guidelines/development-guideline.md)を参照してください。

## 各層の詳細

各層についてはそれぞれのドキュメントを参照してください。

- [ドメイン層](./layers/domain-layer.md) - ビジネスロジックの実装
- [アプリケーション層](./layers/application-layer.md) - ユースケースの実装
- [インフラストラクチャ層](./layers/infrastructure-layer.md) - 技術的実装
- [プレゼンテーション層](./layers/presentation-layer.md) - APIエンドポイント
