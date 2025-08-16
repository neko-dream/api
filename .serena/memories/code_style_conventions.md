# コードスタイルと規約

## 命名規則
- パッケージ名: 小文字、単数形（例: user_usecase, organization_query）
- ファイル名: snake_case（例: create_user.go, list_organizations_query.go）
- 構造体名: PascalCase（例: User, TalkSession）
- インターフェース名: PascalCase + 動詞/名詞（例: UserRepository, AuthService）

## ディレクトリ構造
```
internal/
├── application/     # アプリケーション層
│   ├── usecase/    # 書き込み処理（Command）
│   └── query/      # 読み取り処理（Query）
├── domain/         # ドメイン層
│   ├── messages/   # エラー定義
│   ├── model/      # エンティティ・値オブジェクト
│   └── service/    # ドメインサービス
├── infrastructure/ # インフラストラクチャ層
│   ├── persistence/# DB関連
│   └── di/         # 依存性注入
└── presentation/   # プレゼンテーション層
    └── handler/    # APIハンドラー
```

## 設計方針
- CQRS: 読み書きの分離
- リポジトリパターン: ドメイン層でインターフェース定義
- 依存性注入: digを使用
- コード生成の活用: 手動実装の最小化

## アンチパターン
- ドメイン層への技術的詳細の混入
- アプリケーション層でのドメインロジック実装
- プレゼンテーション層での直接的なDB操作
