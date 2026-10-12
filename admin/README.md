# Refine管理画面のローカル開発

管理画面は`admin/`、APIは`server/`にある。通常ユーザーの認証とは別のGoogle管理者認証を使う。組織所属は管理画面への入場条件にしない。

## 起動

Node.js 24、pnpm、Go 1.25.1、Dockerを使用する。worktreeのルートで依存パッケージをインストールする。

```sh
pnpm install
```

`admin/.env.local`に`VITE_GOOGLE_CLIENT_ID`を設定する。API側の`GOOGLE_CLIENT_ID`と同じ値を使い、Googleの承認済みJavaScript生成元に`http://localhost:3300`を登録する。ブラウザーへ渡すのはクライアントIDだけで、クライアントシークレットは置かない。

worktreeのルートで専用APIを起動する。

```sh
LOCAL_ADMIN_EMAIL=your-google-account@example.com bash scripts/dev-admin-api.sh
```

このコマンドは専用Composeプロジェクト`kotohiro-admin-local`のDBを5434番で起動し、そのDBにだけmigrationと初期管理者登録を行う。APIは8081番。既存の5432・5433番のDBとdev・prodは使わない。管理者登録をしない場合は`LOCAL_ADMIN_EMAIL`を省略する。

APIの設定読み込みは祖先ディレクトリの`.env`も探索する。起動スクリプトはDB・環境・署名キーをローカル値に固定する。GoogleクライアントIDなどは明示的な環境変数か、その設定ファイルで指定する。

別ターミナルで管理画面を起動する。

```sh
cd admin
pnpm dev
```

ブラウザーで`http://localhost:3300/admin/`を開く。Viteが`/v1`と`/auth`を8081番へプロキシするので、APIを別オリジンとして直接呼ばない。開発サーバーはループバックにだけ公開する。

## ビルドと確認

```sh
cd admin
pnpm test
pnpm typecheck
pnpm build
```

ビルド出力は`server/static/admin-ui/`。このworktreeの既存API配信先をRefine版で更新する。APIの再ビルド・再起動後に`http://localhost:8081/admin/`から配信される。

サーバー側の認証テストは`server/`で実行する。

```sh
ENV=local DATABASE_URL='postgres://kotohiro:kotohiro@127.0.0.1:5434/kotohiro?sslmode=disable' \
  go test -race ./internal/presentation/handler ./internal/application/usecase/admin_usecase ./internal/infrastructure/auth/admin
```

## 確認する動作

- 未ログインで保護画面を開くとログイン画面へ戻る。
- 許可リストに登録済みのGoogleアカウントだけがログインできる。
- ログイン後にトークセッション一覧・管理者一覧を取得できる。
- 一覧取得失敗は空一覧と区別してエラーを表示し、再試行できる。
- 期限切れのaccess tokenは401を返し、並行リクエストでもrefreshは1回にまとめる。再試行は最大1回。
- ログアウトはrefreshの失効とCookie削除を行う。失敗した場合は成功扱いにしない。

access JWT自体の即時失効は現行実装にない。ログアウトはブラウザーのCookie削除とrefresh失効であり、既に取得されたaccess JWTは最長15分有効という制限がある。

APIの障害を「管理者未登録」と表示しない。一般ユーザーの`SessionId`を管理者認証に流用しない。ユーザー管理・配布アカウント作成・作成権限の付与は次の作業であり、ここでは追加しない。
