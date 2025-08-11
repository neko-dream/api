#!/bin/bash

set -eu

if [ $# -eq 0 ]; then
    echo "Usage: $0 <dev|prod>"
    exit 1
fi

ENV=$1
if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
    echo "Environment must be 'dev' or 'prod'"
    exit 1
fi

source scripts/utils.sh

export ACCOUNT_ID=$(login_aws $AWS_PROFILE)
if [ -z "$ACCOUNT_ID" ]; then
    echo "Failed to get AWS account ID" >&2
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

