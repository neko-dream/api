<div align="center">
    <img src="docs/public/assets/icon.png" width="50%" center />


# kotohiro API Server

ことひろは意見や言葉を重ねて、よりよい意思決定を目指すサービスです。

[![Go version][go_version_img]][go_dev_url]
[![Go report][go_report_img]][go_report_url]
[![License][repo_license_img]][repo_license_url]

</div>

## 🚀 セットアップ

### 1. 前提パッケージのインストール

miseを用いてツールの管理をしています。

```sh
# macosの場合
brew install mize

# その他の環境
curl https://mise.run | sh
```

### 2. プロジェクトのセットアップ

リポジトリのクローン

```sh
git clone https://github.com/neko-dream/api.git
cd api
```

必要なツールのインストール

```sh
mise install
```

環境変数の設定

```sh
cp .env.example .env
```

### 3. コード生成

自動生成スクリプトを実行

```sh
./scripts/gen.sh
```

### 3. DBのセットアップ

docker composeを使用します

```sh
docker compose up -d db
```

### 4. ホットリロード

airを用いてホットリロードが可能

```sh
air
```

## ⚖️ License

ことひろ は GNU Affero General Public License v3.0 (AGPL-3.0) の下で提供されています。詳細については[LICENSE](LICENSE)をご確認ください。


<!-- Go links -->

[go_report_url]: https://goreportcard.com/report/github.com/neko-dream/server
[go_report_img]: https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none
[go_dev_url]: https://pkg.go.dev/github.com/neko-dream/api
[go_version_img]: https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go

<!-- Repository links -->

[repo_license_url]: https://github.com/gowebly/gowebly/blob/main/LICENSE
[repo_license_img]: https://img.shields.io/badge/license-AGPL--v3-red?style=for-the-badge&logo=none

