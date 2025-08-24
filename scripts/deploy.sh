#!/bin/bash

set -eu

SCRIPT_DIR=$(cd $(dirname $0); pwd)
source "${SCRIPT_DIR}/utils.sh"

if [ $# -eq 0 ]; then
    print_error "エラー: 環境を指定してください"
    print_info "使い方: $0 <dev|prod>"
    print_info "例:"
    echo "  $0 dev   # 開発環境にデプロイ"
    echo "  $0 prod  # 本番環境にデプロイ"
    exit 1
fi

ENV=$1
if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
    print_error "エラー: 無効な環境 '$ENV'"
    print_info "有効な環境: dev, prod"
    exit 1
fi

export ACCOUNT_ID=$(login_aws $AWS_PROFILE)
if [ -z "$ACCOUNT_ID" ]; then
    print_error "AWSアカウントIDの取得に失敗しました"
    exit 1
fi

aws --profile $AWS_PROFILE ecr get-login-password --region ap-northeast-1 |
docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com

cd admin-ui
if [ "$ENV" = "prod" ]; then
    bun build:prod
else
    bun build:dev
fi
cd ../

go mod tidy
go mod download
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -tags timetzdata -ldflags="-s -w" -trimpath -o server ./cmd/server

COMMIT_HASH=$(git show --format='%h' --no-patch)
TIMESTAMP=$(date '+%s%3')
export IMAGE_TAG=$COMMIT_HASH$TIMESTAMP

if [ "$ENV" = "prod" ]; then
    ECR_REPO="kotohiro-prd-api"
else
    ECR_REPO="kotohiro-dev-api"
fi
TAG=$ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com/$ECR_REPO:$IMAGE_TAG

docker build -f ./container/Dockerfile.gh . -t $TAG --no-cache
docker push $TAG

rm -rf ./server

if [ "$ENV" = "prod" ]; then
    ecspresso deploy --config ./.ecspresso/api/prd/ecspresso.yml --force-new-deployment
else
    ecspresso deploy --config ./.ecspresso/api/dev/ecspresso.yml --force-new-deployment
fi

