# Web + API のローカル開発

## 構成

- Web: `https://localhost:3000`
- ブラウザ用 API: `https://localhost:3000/api`（Vite が `/api` を除いて API に中継）
- API / SSR 用 API: `http://127.0.0.1:8080`
- PostgreSQL: `127.0.0.1:5433`（専用 Compose project / volume）
- メール受信: `http://127.0.0.1:8025`（Mailpit。SMTP は `127.0.0.1:1025`）

ブラウザからは Web と API が同一オリジンになる。ローカル環境の
`SessionId` は Domain 指定なしの host-only Cookie で、`Secure` / `HttpOnly` /
`SameSite=Lax` を維持する。SSR は受信した Cookie を HTTP のローカル API に
転送するので、自己署名証明書の検証を無効化する設定は不要。
ブラウザの変更系API操作はこのHTTPS Webプロキシ経由で行う。
HTTPの直接Swaggerや別ポートの管理画面開発サーバーは許可Originに含めない。
本番では既存のCookie Domain / SameSite=Laxを維持する。
developmentでは既存のクロスサイト認証修正（SameSite=None）を含む。

## 起動

前提: Docker（Rancher Desktop 等）起動済み、Go（`.mise.toml` 参照）、
Web 用の Node 24+ と pnpm。Web は `dev:local` 対応版を使う。

API リポジトリで:

```sh
bash scripts/dev-local.sh
```

専用の DB を起動し、API 起動時にマイグレーションを適用する。
`.env` の DB 接続先は使わない。その他のアプリ設定は通常の `.env` を読み込む。
DB データは再起動後も保持される。

別のターミナルで Web リポジトリに移動して:

```sh
pnpm install --frozen-lockfile
pnpm dev:local
```

`CLOUDFLARE_ENV` は指定しない。`https://localhost:3000` を開く。
ローカル証明書は `.cert/localhost` に自動生成される。自己署名のため、初回は
ブラウザで証明書の警告を確認して続行する。`localhost` と `127.0.0.1` を
ブラウザのアドレスとして混在させない。hosts ファイルの変更は不要。

## 認証の確認

トップ画面からゲストログインし、再読み込み後もログイン状態が維持されるか確認する。
開発ログインも Web 経由で呼べる:

```text
https://localhost:3000/api/auth/dev/login?id=local-user&redirect_url=https%3A%2F%2Flocalhost%3A3000%2Fauth%2Fsignup
```

Google / LINE を使う場合は、プロバイダーの管理設定にローカルの callback URL
登録が必要。Cookie 接続の確認だけならゲスト / 開発ログインで外部認証は不要。
組織専用アカウントのログインはこの仕組みとは別途実装する。

疎通と Cookie 発行・認証・失効を確認するには:

```sh
bash scripts/check-local-auth.sh
```

## メール送信・招待のテスト

`ENV=local` のメールはローカルSMTPに送る。通常の開発起動で使う場合は
`.env` に `SMTP_ADDR=127.0.0.1:1025` を指定する。省略時もこのアドレスを使う。
外部SMTPはローカル環境では拒否し、development / production は従来どおり SES を使う。

メール基盤と招待処理をまとめて検証する:

```sh
bash scripts/test-mail.sh
```

外部送信なしのSES・本文検証、ローカルSMTPの正常／失敗／キャンセル／期限切れ、
組織招待の単体テストに加えて、PostgreSQLの一時スキーマで送信失敗時の
ロールバック・再試行・再招待を検証する。CIでもDB統合テストを実行する。
一時スキーマは終了時に削除し、通常のデータは変更しない。

Mailpitの画面でローカル送信メールを確認できる。
ログインごとのメールリンク／PIN本人確認は [メールログイン](password-email-login.md) を参照。

送信とDBコミットは分散トランザクションではない。送信が成功してから
DBコミットが失敗した場合や、送信結果がネットワーク障害で不明な場合の厳密な
整合性には、永続的なoutbox等が必要。現段階で exactly-once 配信は保証しない。

## 停止

Web と API のターミナルで Ctrl+C。DB は次で停止する（データは保持）:

```sh
docker compose -p kotohiro-local -f compose.local.yaml stop
```
