#!/bin/bash

# 色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 色なし

# メッセージ表示関数
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 指定されたポートを使用しているプロセスを終了
kill_port() {
    local port=$1
    local pids=$(lsof -t -i:$port)
    if [ ! -z "$pids" ]; then
        print_message $YELLOW "ポート$port のプロセスを終了します (PID: $pids)"
        kill -9 $pids 2>/dev/null
        sleep 1
    fi
}

print_error() {
    print_message $RED "$1"
}

print_success() {
    print_message $GREEN "$1"
}

print_warning() {
    print_message $YELLOW "$1"
}

print_info() {
    print_message $BLUE "$1"
}

login_aws() {
  aws sts get-caller-identity --profile $AWS_PROFILE &>/dev/null || aws sso login --profile $AWS_PROFILE &>/dev/null
  if [ -z "$1" ]; then
      echo "AWS profile name is required" >&2
      return 1
  fi

  ACCOUNT_ID=$(aws sts get-caller-identity --query "Account" --output text --profile $1)
  if [ -z "$ACCOUNT_ID" ]; then
      echo "Failed to get AWS account ID" >&2
      return 1
  fi
  echo $ACCOUNT_ID
}

urlencode() {
    local string="${1}"
    local strlen=${#string}
    local encoded=""
    local pos c o

    for (( pos=0 ; pos<strlen ; pos++ )); do
        c=${string:$pos:1}
        case "$c" in
            [-_.~a-zA-Z0-9] ) o="${c}" ;;
            * ) printf -v o '%%%02x' "'$c" ;;
        esac
        encoded+="${o}"
    done
    echo "${encoded}"
}
