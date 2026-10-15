#!/bin/bash
# 手元のチェックアウトを ECS にデプロイする。GitHub Actions（deploy-dev.yml / deploy-prod.yaml）と同じ手順で行う。
#
# 使い方: server/scripts/deploy.sh <dev|prod> [--yes]
#   - デプロイするのは「実行したディレクトリの git チェックアウト」の内容
#   - AWS_PROFILE が未設定なら peco で選ぶ
#   - ADMIN_ALLOWED_ORIGINS / INITIAL_ADMIN_EMAIL / INITIAL_ADMIN_NAME / ADMIN_GOOGLE_CLIENT_ID は、
#     未設定なら GitHub の Environment（development / production）の変数を gh で読む
#   - prod は main ブランチで、未コミットの変更がないときだけデプロイできる

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)
source "${SCRIPT_DIR}/utils.sh"

usage() {
    print_info "使い方: $0 <dev|prod> [--yes]"
    echo "  $0 dev    # 開発環境にデプロイ"
    echo "  $0 prod   # 本番環境にデプロイ（main ブランチのみ）"
    echo "  --yes     # 確認を省略する"
}

ENV=${1:-}
ASSUME_YES=false
[ "${2:-}" = "--yes" ] && ASSUME_YES=true

case "$ENV" in
    dev)
        GH_ENVIRONMENT=development
        APP_ENV=dev
        ECR_REPO=kotohiro-dev-api
        ;;
    prod)
        GH_ENVIRONMENT=production
        APP_ENV=prd
        ECR_REPO=kotohiro-prd-api
        ;;
    *)
        print_error "エラー: 環境を dev か prod で指定してください"
        usage
        exit 1
        ;;
esac

for cmd in aws docker ecspresso go pnpm gh git; do
    if ! command -v "$cmd" &>/dev/null; then
        print_error "$cmd がインストールされていません"
        exit 1
    fi
done
if ! docker info &>/dev/null; then
    print_error "Docker が起動していません"
    exit 1
fi

REPO_ROOT=$(git rev-parse --show-toplevel)
SERVER_DIR="${REPO_ROOT}/server"
BRANCH=$(git -C "$REPO_ROOT" branch --show-current)
COMMIT_HASH=$(git -C "$REPO_ROOT" rev-parse --short HEAD)
DIRTY=$(git -C "$REPO_ROOT" status --porcelain --untracked-files=no)

if [ "$ENV" = "prod" ]; then
    if [ "$BRANCH" != "main" ] || [ -n "$DIRTY" ]; then
        print_error "本番は main ブランチで、未コミットの変更がないときだけデプロイできます（現在: ${BRANCH:-detached}）"
        exit 1
    fi
elif [ -n "$DIRTY" ]; then
    print_warning "未コミットの変更があります。変更を含めてデプロイします"
fi

# GitHub Actions と同じ値を使う。環境変数で上書きできる
github_var() {
    local name=$1
    local value=${!name:-}
    if [ -z "$value" ]; then
        value=$(gh variable get "$name" --env "$GH_ENVIRONMENT" 2>/dev/null || true)
    fi
    echo "$value"
}
ADMIN_ALLOWED_ORIGINS=$(github_var ADMIN_ALLOWED_ORIGINS)
INITIAL_ADMIN_EMAIL=$(github_var INITIAL_ADMIN_EMAIL)
INITIAL_ADMIN_NAME=$(github_var INITIAL_ADMIN_NAME)
INITIAL_ADMIN_NAME=${INITIAL_ADMIN_NAME:-管理者}
ADMIN_GOOGLE_CLIENT_ID=$(github_var ADMIN_GOOGLE_CLIENT_ID)
for name in ADMIN_ALLOWED_ORIGINS INITIAL_ADMIN_EMAIL ADMIN_GOOGLE_CLIENT_ID; do
    if [ -z "${!name}" ]; then
        print_error "$name を取得できませんでした。GitHub の ${GH_ENVIRONMENT} 環境の変数か、環境変数で指定してください"
        exit 1
    fi
done
export ADMIN_ALLOWED_ORIGINS INITIAL_ADMIN_EMAIL INITIAL_ADMIN_NAME

# AWS_PROFILEが設定されていない場合はpecoで選択
if [ -z "${AWS_PROFILE:-}" ]; then
    if ! command -v peco &>/dev/null; then
        print_error "AWS_PROFILE 環境変数を設定するか、peco をインストールしてください"
        exit 1
    fi
    print_info "AWSプロファイルを選択してください:"
    AWS_PROFILE=$(aws configure list-profiles | peco --prompt "AWS Profile > ")
    if [ -z "$AWS_PROFILE" ]; then
        print_error "プロファイルが選択されませんでした"
        exit 1
    fi
fi
export AWS_PROFILE

print_info "デプロイ先: ${ENV} / ブランチ: ${BRANCH:-detached} (${COMMIT_HASH}) / AWSプロファイル: ${AWS_PROFILE}"
if [ "$ASSUME_YES" != true ]; then
    read -r -p "デプロイしますか？ [y/N] " answer
    [ "$answer" = "y" ] || [ "$answer" = "Y" ] || { print_warning "中止しました"; exit 1; }
fi

ACCOUNT_ID=$(login_aws "$AWS_PROFILE")
if [ -z "$ACCOUNT_ID" ]; then
    print_error "AWSアカウントIDの取得に失敗しました"
    exit 1
fi
REGISTRY="${ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com"
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin "$REGISTRY"

# ビルドはコミット済みのファイルを書き換える（管理画面の index.html、go mod download による go.work.sum）。
# 元々変更がなかったものは、終了時に元に戻して作業ツリーを汚さない
RESTORE_FILES=()
for file in server/static/admin-ui/index.html go.work.sum; do
    if git -C "$REPO_ROOT" ls-files --error-unmatch -- "$file" &>/dev/null && git -C "$REPO_ROOT" diff --quiet -- "$file"; then
        RESTORE_FILES+=("$file")
    fi
done
cleanup() {
    rm -f "${SERVER_DIR}/server"
    if [ ${#RESTORE_FILES[@]} -gt 0 ]; then
        git -C "$REPO_ROOT" checkout -- "${RESTORE_FILES[@]}"
    fi
}
trap cleanup EXIT

# 管理画面をビルドする（出力先は server/static/admin-ui。サーバーのバイナリに埋め込まれる）
print_info "管理画面をビルドしています"
(cd "$REPO_ROOT" && pnpm install --frozen-lockfile)
(cd "${REPO_ROOT}/admin" && VITE_GOOGLE_CLIENT_ID="$ADMIN_GOOGLE_CLIENT_ID" VITE_API_URL='' pnpm build)

print_info "サーバーをビルドしています"
cd "$SERVER_DIR"
go mod download
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags timetzdata -ldflags="-s -w" -trimpath -o server ./cmd/server

export IMAGE_TAG="${COMMIT_HASH}-$(date '+%s')"
TAG="${REGISTRY}/${ECR_REPO}:${IMAGE_TAG}"
docker build --platform linux/arm64 -f ./container/Dockerfile.gh . -t "$TAG" --build-arg ENV="$APP_ENV"
docker push "$TAG"

ecspresso deploy --config "./.ecspresso/api/${APP_ENV}/ecspresso.yml" --force-new-deployment
print_success "デプロイしました: ${TAG}"
