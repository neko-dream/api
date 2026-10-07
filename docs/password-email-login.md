# ログインごとのメール本人確認

パスワードの漏えいだけでログインできないよう、パスワードログインでは毎回
メールの確認リンクまたは6桁PINと、ログイン開始時のブラウザCookieの両方を必要とする。
過去の `emailVerified` は今回の本人確認の代わりにしない。

## ローカルで試す

Web + API のローカル開発環境を起動する（Web は `dev:local` 対応版）。
API は `ENV=local`、SMTP は `127.0.0.1:1025` の Mailpit を使う。
`EMAIL_VERIFICATION_BASE_URL` は省略すると local のみ
`https://localhost:3000/api` を使う。`WEBSITE_URL=https://localhost:3000` を設定する。

専用のローカルDB（127.0.0.1:5433）に検証用アカウントを作成する:

```sh
ENV=local go run ./cmd/local-email-account city-login@localhost.test local-email-login-only
```

このコマンドはリモートDBや `.env` を参照しない。既存アカウントのパスワードを
上書きせず、重複時はエラーにする。組織所属は作成しない。

1. `https://localhost:3000/` の「自治体・組織の方はこちら」から
   `https://localhost:3000/organization/login` を開く。
2. 利用規約への同意を選び、アカウントとパスワードを入力してログインする。
3. 同じブラウザで `http://localhost:8025` を開き、Mailpitのメールを確認する。
4. メールのリンクを開き、「ログインを確認する」を押す。
   リンクを開けない場合は、元の確認待ち画面にメールの6桁PINを入力し「PINで確認する」を押す。
5. 本人確認が完了すると初めて `SessionId` Cookieを発行し、Webの `/home` へ移動する。

外部メールは送らない。Webは既存カード・入力部品を流用し、一般向け入口と
自治体・組織向けのメールログインページを分離する。Figmaは入力・確認部品の参考に限定する。
専用ページへアクセスするだけで組織権限は付与されない。自動組織ログインは #319 の別作業。
メールリンクの遷移先は `EMAIL_VERIFICATION_PAGE_URL=https://localhost:3000/auth/email/verify`。
APIの診断ページ `/api/auth/password/verify` も残すが、そこにあるローカル検証フォームは
development / production には表示しない。

## API

- `POST /auth/password/login`: 既存のmultipart形式で `idOrEmail` / `password` を送る。
  成功時は **202** と `{verificationRequired:true, expiresAt}`、
  `__Host-PasswordLoginPending` Cookieを返す。認証済みセッションは作成しない。
- `GET /auth/password/verify`: 確認ページ。リンクを消費せず、Cookieを発行しない。
- `POST /auth/password/verify`: JSON `{token}` と開始時のCookieで本人確認する。
  成功時は **200** と `{redirectURL}`、`SessionId` Cookieを返し、開始時Cookieを削除する。
  無効・期限切れ・別ブラウザ・再利用は400。
- `POST /auth/password/verify-pin`: JSON `{pin:"000001"}` と開始時Cookieで本人確認する。
  先頭の0も文字列として保持する。成功時のCookie・移動先はリンクと同じ。
  間違いは400、5回目以降は429。失敗回数はトランザクションをコミットしてDBへ記録する。
- リンクとPINは同じ確認情報を共有し、どちらかの成功で両方を消費する。
  PINを5回間違えた場合も両方を無効化する。別ブラウザの試行では回数を増やさない。
- 再送は同じ `POST /auth/password/login` に再度パスワードを送る。60秒未満は429。
  再送成功時は旧リンク・PIN・開始Cookieを無効化し、失敗回数をリセットする。

メールURLの秘密値はフラグメント（`#...`）に置く。確認ページがボタン操作時に
JSON本文でAPIへ送るので、アクセスログやRefererには含まれない。
開始時Cookieはhost-onlyの `__Host-` プレフィックス付きで、他サブドメインからの
Cookie植え付けを防ぐ。CookieはSecure / HttpOnly。POSTのOriginは設定済みのWeb/APIだけを許可する。
本番Webのwwwとapexは同じアプリを配信するため両方を明示許可し、確認後の移動先はwwwへ統一する。
Credential付きCORSも同じ許可リストを使い、任意HTTPSサイトやworkers.devプレビューは許可しない。
ブラウザからの変更系リクエストも未許可Originを拒否し、単純フォームによるクロスサイト書き込みを防ぐ。
非ブラウザクライアントではOriginを省略できるが、秘密値とCookieの両方が必要。

## 検証

```sh
MAIL_TEST_DATABASE_URL='postgres://kotohiro:kotohiro@127.0.0.1:5433/kotohiro?sslmode=disable' \
go test -race ./internal/application/usecase/auth_usecase ./internal/domain/service ./internal/presentation/handler ./internal/infrastructure/email/...
```

DB統合テストは専用の一時スキーマを作成・削除し、通常の行は変更しない。
期限境界、再利用、同時使用（成功は1回だけ）、再送、送信失敗、セッション発行失敗、
パスワード・メール変更後の旧リンク無効化を含む。
PINの先頭0、別ブラウザ、期限、5回制限の永続化、並行失敗時の回数上限、
リンクとPINの競合（成功1回）、再送とセッション失敗時のロールバックも検証する。

## 制約・移行時の注意

- `ENV` は `production` / `development` / `local` のいずれか必須。
  未指定・誤記は起動時に拒否する。開発用ログイン・登録は `development` / `local` のみ許可する。
- 本番のWeb・メール確認リンク・認証後URLは `https://www.kotohiro.com` に統一する。
  Webのビルドには `CLOUDFLARE_ENV=production` が必要で、設定欠落時は停止する。

- 確認リンクは15分で失効し、DBにはランダム秘密値ではなくSHA-256ハッシュを保存する。
- PINも15分で失効する。DBにはPINそのものではなく、サーバーの `TOKEN_SECRET` と
  ブラウザ情報に紐づくHMAC-SHA256を保存する。キー変更で発行済みPINは使えなくなる。
- メールとDBコミットの原子性は未対応（#317のoutbox残課題）。
- パスワードログインのレスポンスが200から202へ変わるため、他の既存クライアントも
  「ログイン完了」ではなく「メール確認待ち」を扱う必要がある。Webの正式な画面は対応済み。
- この変更前に発行されたセッションは一括失効しない。次回の本人確認成功時に
  そのユーザーの既存セッションを失効する。
  このDB失効処理はパスワード本人確認に限定し、OAuthの共有SessionServiceは変更しない。
- Google / LINE / ゲストの既存ログインは変更しない。developmentの登録・開発用ログインは
  検証用のまま残す。本番のパスワード登録は禁止済み。
- 組織専用の1アカウント制・自動組織ログイン・同時セッションは #319 の範囲で、
  この実装だけでは完了しない。初期パスワードの紛失時の復旧も別途必要。
- メールアカウントの侵害、メールとパスワード双方のフィッシングは防げない。
  通常のパスワード推測対策（IPごとのレート制限等）は別途必要。
