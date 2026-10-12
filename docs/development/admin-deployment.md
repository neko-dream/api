# 管理画面のデプロイ

現段階ではRefineをAPIの `/admin/` から配信する。CIは `admin/` をビルドしてから `server/` のGoバイナリへ埋め込む。旧 `admin-ui/` はビルドしない。S3 + CloudFrontへの配信分離は未実装。

## 環境設定

GitHub Environmentsの `development` / `production` にそれぞれ以下のVariablesを設定する。

| Variable | 内容 |
| --- | --- |
| INITIAL_ADMIN_EMAIL | 最初の管理者のGoogleアカウントのメールアドレス |
| INITIAL_ADMIN_NAME | 表示名。省略時は「管理者」 |
| ADMIN_GOOGLE_CLIENT_ID | APIのGOOGLE_CLIENT_IDと同じ公開クライアントID |

初期管理者メールとクライアントIDが未設定の場合、CIはデプロイ前に失敗する。ローカルのメールやGoogle IDを本番用の既定値にはしない。

Google OAuth設定で `https://api.dev.kotohiro.com` / `https://api.kotohiro.com` を各環境の承認済みJavaScriptオリジンへ登録する。この管理画面はGoogle IDトークンをAPIへ渡すため、一般ユーザーのOAuthコールバック設定とは別にオリジンの登録が必要。

## 初期管理者

ECSへINITIAL_ADMIN_EMAIL / INITIAL_ADMIN_NAMEを渡し、マイグレーション完了後、HTTPサーバー起動前に登録する。指定メールが登録済みなら変更しない。未設定なら誰も登録しない。「DBが空なら最初のログインユーザーを管理者にする」という経路はない。Googleトークンの検証は通常どおり必須。

envが残っている間は指定メールがDBから削除されても次回起動時に再登録される。初回登録後に再登録を止める場合は、タスク定義からINITIAL_ADMIN_EMAILを外し、CIの必須チェックも変更する。メールを変更した場合、旧管理者を自動削除する機能はない。

## リリース前の確認

- Aurora移行完了と、適用済みマイグレーション番号・内容の一致を確認する。別ブランチで53以降の番号が衝突しているため、番号だけで互換性を判断しない。
- 上記GitHub VariablesとGoogleの承認済みオリジンを設定する。
- devでログイン・一覧取得・refresh・ログアウトを確認した後、本番へ進める。

この変更自体ではAWSリソース、既存DB、GitHub Variables、Google OAuth設定を変更しない。
