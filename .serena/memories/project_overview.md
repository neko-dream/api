# プロジェクト概要

## プロジェクト名: kotohiro API Server
意見や言葉を重ねて、よりよい意思決定を目指すサービスのAPIサーバー

## テックスタック
- 言語: Go 1.24.3
- フレームワーク:
  - API定義: TypeSpec → OpenAPI → ogen
  - SQL: SQLC
  - DI: dig
  - テレメトリ: OpenTelemetry
- データベース: PostgreSQL 16 + PostGIS
- 開発ツール: mise（ツール管理）、air（ホットリロード）

## アーキテクチャ
- DDDベースのレイヤードアーキテクチャ（4層構造）
  - presentation層: APIエンドポイント
  - application層: ユースケース（CQRS実装）
  - domain層: ビジネスロジック
  - infrastructure層: 技術的実装
- 依存関係: presentation → application → domain ← infrastructure

## ライセンス
GNU Affero General Public License v3.0 (AGPL-3.0)