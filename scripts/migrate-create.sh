#!/usr/bin/env bash
if [ -z $1 ]; then
  echo "マイグレーション名を指定してください"
  echo "使用方法: mise run migrate-create <マイグレーション名>"
  exit 1
fi
MIGRATION_DIR="migrations"
mkdir -p $MIGRATION_DIR
MIGRATION_NAME=$(echo $1 | tr '[:upper:]' '[:lower:]' | tr ' ' '_')
migrate create -ext sql -dir $MIGRATION_DIR -seq $MIGRATION_NAME
