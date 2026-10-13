# メール送信・組織招待のローカル検証

管理画面からの一般ユーザー作成と初回パスワード設定は [ユーザー招待](admin-user-invitations.md) を参照。

## メール送信先

`ENV=local` では `SMTP_ADDR` のローカルSMTPに送信する。
省略時は `127.0.0.1:1025` を使い、ループバック以外の送信先は拒否する。
development / production は従来どおり SES を使う。

Docker が起動済みの状態で、受信環境を起動する:

```sh
docker compose -p kotohiro-local -f compose.local.yaml up -d --wait db mailpit
```

- SMTP: `127.0.0.1:1025`
- Mailpit の受信画面: `http://127.0.0.1:8025`
- 検証用 PostgreSQL: `127.0.0.1:5433`

API をローカル起動する場合は、`EMAIL_FROM` に有効な送信元アドレス、
`WEBSITE_URL` に有効な HTTP(S) URL を設定する。

## テスト

リポジトリの `.mise.toml` に対応する Go と Docker を使う:

```sh
bash scripts/test-mail.sh
```

SES のリクエストは偽物のクライアントで検証し、外部にメールを送らない。
SMTP の正常・失敗・キャンセル・期限切れ、本文と入力の検証、
組織招待の単体テストを `-race` 付きで実行する。
Mailpit への実際の送信と、PostgreSQL 上の送信失敗時のロールバック、
再試行、再招待も検証する。

DB 統合テストは専用の一時スキーマに全マイグレーションを適用し、
終了時に削除する。通常のアプリケーションの行は変更しない。
CI では `MAIL_TEST_DATABASE_URL`、`MAIL_TEST_SMTP_ADDR`、
`MAIL_TEST_HTTP_URL` を設定し、同じ統合テストを実行する。

## 対象と制約

### デプロイ前の送信元設定

dev / prd の API タスク定義の `environment` に
`EMAIL_FROM=noreply@kotohiro.com` を設定する。送信元アドレスは秘密情報では
ないため、Secrets Manager に `EMAIL_FROM` キーを追加する必要はない。

2026-10-08 に admin プロファイルで東京リージョン（`ap-northeast-1`）の
SES を照会し、`noreply@kotohiro.com` が検証済み（`SUCCESS`）かつ
送信可能（`SendingEnabled=true`）であることを確認した。
送信元や AWS アカウント・リージョンを変更する場合は検証状態を再確認する。
このリポジトリの変更だけでは ECS の稼働中タスクや SES の設定は変更されない。

### 残課題

対象はメール基盤と既存の組織招待の改善（#317）。
ログインのたびのメール本人確認（#318）や組織専用アカウント（#319）は含まない。
再招待で既存アカウントのパスワードや既存所属の権限を変更しない。

メール送信は引き続き DB トランザクション内で行う。
送信成功後に DB コミットが失敗した場合や、ネットワーク障害で送信結果が
不明な場合の厳密な整合性には、永続的な outbox 等が必要。
現段階で exactly-once 配信や同時招待の重複防止は保証しない。

検証サービスを停止する（DB データは保持）:

```sh
docker compose -p kotohiro-local -f compose.local.yaml stop
```
