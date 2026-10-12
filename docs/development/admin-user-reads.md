# 管理画面のユーザー閲覧

運営の管理画面アカウントで、ことひろの一般ユーザーを閲覧する。管理画面アカウント自体の一覧（管理者メニュー）とは別の画面。

## 画面

- `/admin/users`: 表示名・表示ID・ユーザーIDによる部分一致検索、20件ずつのページ送り、メール・状態・登録日時の表示。
- `/admin/users/show/:id`: 基本情報、確認済みメールの状態、ログイン方法、登録・更新・退会日時の表示。
- 検索とページはURLに保持する。詳細へのリンクは内部ユーザーIDを使う。
- 一覧取得失敗と空一覧、詳細の404とサーバー障害を区別し、再試行できる。

## API

- GET `/v1/manage/users/list?offset=0&limit=20&search=...`
- GET `/v1/manage/users/{userID}`

既存の未実装一覧APIを実装し、レスポンスを `{users, totalCount}` に変更した。管理画面用OpenAPIとGoコードを再生成する。現在は登録日時降順・ユーザーID降順の固定順序で、未実装だった並び替えパラメーターは削除した。

既存のAdminCookieAuthを必須とする。一般ユーザーのSessionIdや組織への所属では認証されない。メールは既存Encryptorで復号する。復号失敗は既存ユーザー詳細と同じ方法で記録し、そのメールを非表示・未確認として他の項目の閲覧を続ける。暗号化メールに対する文字列検索は行わない。パスワードハッシュ・ソルト・OAuth subject・トークン・人口統計情報は返さない。

読み取り専用で、ユーザーの作成・編集・削除・権限付与は含まない。独立したセッション作成権限の機構は別途実装するため、この画面には推測した権限状態を表示しない。マイグレーションは不要。

管理画面のビルド成果物はコミットしない。`server/static/admin-ui/index.html`はgo:embed用のプレースホルダーとして残し、デプロイ時の`pnpm build`で環境ごとの成果物を生成する。

## 検証

`admin/`で `pnpm test` / `pnpm typecheck` / `pnpm build`。
`server/`で `go test -race ./internal/presentation/handler ./internal/infrastructure/persistence/query/user`。
DBクエリテストには、migration 57まで適用した使い捨てDBの `ADMIN_TEST_DATABASE_URL` を指定する。作成したテストユーザーはトランザクションをrollbackする。
