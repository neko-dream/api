#!/bin/bash

# 出力用のカラーコード
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # 色なし

# 使い方を表示
show_usage() {
    echo -e "${BLUE}使い方:${NC}"
    echo "  $0 [テストファイルパス]"
    echo ""
    echo -e "${BLUE}例:${NC}"
    echo "  $0                                    # すべてのE2Eテストを実行 (./e2e/**/*)"
    echo "  $0 ./e2e/auth/*                       # authディレクトリのテストのみ実行"
    echo "  $0 ./e2e/auth/test_withdraw_user.yaml # 特定のテストファイルのみ実行"
    echo ""
    exit 0
}

# ヘルプオプションの確認
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_usage
fi

# 色付きメッセージを表示する関数
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# 指定ポートのプロセスを終了する関数
kill_port() {
    local port=$1
    local pids=$(lsof -t -i:$port)
    if [ ! -z "$pids" ]; then
        print_message $YELLOW "ポート$port のプロセスを終了します (PID: $pids)"
        kill -9 $pids 2>/dev/null
        sleep 1
    fi
}

# サーバーの起動を待つ関数
wait_for_server() {
    local port=$1
    local max_attempts=30
    local attempt=0

    print_message $YELLOW "ポート$port でサーバーの起動を待っています..."

    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:$port/health >/dev/null 2>&1; then
            print_message $GREEN "サーバーが起動しました！"
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 1
        echo -n "."
    done

    echo ""
    print_message $RED "サーバーが30秒以内に起動しませんでした"
    return 1
}

# クリーンアップ関数
cleanup() {
    print_message $YELLOW "\nクリーンアップ中..."
    if [ ! -z "$SERVER_PID" ] && kill -0 $SERVER_PID 2>/dev/null; then
        print_message $YELLOW "サーバーを停止します (PID: $SERVER_PID)"
        kill -TERM $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
    fi
    kill_port 3000
    exit $EXIT_CODE
}

# 終了時にクリーンアップを実行するようトラップを設定
trap cleanup EXIT INT TERM

# メイン処理
print_message $GREEN "=== テストスイート開始 ==="

# テストパスを設定（引数がなければデフォルト）
TEST_PATH=${1:-"./e2e/**/*"}
print_message $YELLOW "テスト対象: $TEST_PATH"

# ポート3000の既存プロセスを終了
kill_port 3000

# サーバーを起動
print_message $YELLOW "サーバーを起動しています..."
go run ./cmd/server/main.go &
SERVER_PID=$!

# サーバーが正常に起動したか確認
if [ -z "$SERVER_PID" ] || ! kill -0 $SERVER_PID 2>/dev/null; then
    print_message $RED "サーバーの起動に失敗しました"
    EXIT_CODE=1
    exit 1
fi

print_message $GREEN "サーバーを起動しました (PID: $SERVER_PID)"

# サーバーの準備完了を待つ
if ! wait_for_server 3000; then
    EXIT_CODE=1
    exit 1
fi

# テストを実行
print_message $GREEN "\n=== E2Eテスト実行中 ==="
runn run $TEST_PATH --verbose --debug

# テストの終了コードを取得
EXIT_CODE=$?

# テスト結果を表示
if [ $EXIT_CODE -eq 0 ]; then
    print_message $GREEN "\n✅ すべてのテストが成功しました！"
else
    print_message $RED "\n❌ いくつかのテストが失敗しました (終了コード: $EXIT_CODE)"
fi

# クリーンアップはトラップにより自動的に実行されます
