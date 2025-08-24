#!/bin/bash

SCRIPT_DIR=$(cd $(dirname $0); pwd)
source "${SCRIPT_DIR}/utils.sh"

# 使い方を表示
show_usage() {
    print_info "使い方:"
    echo "  $0 [テストファイルパス]"
    echo ""
    print_info "例:"
    echo "  $0                                    # すべてのE2Eテストを実行 (./e2e/**/*)"
    echo "  $0 ./e2e/auth/*                       # authディレクトリのテストのみ実行"
    echo "  $0 ./e2e/auth/test_withdraw_user.yaml # 特定のテストファイルのみ実行"
    echo ""
    exit 0
}

if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_usage
fi

wait_for_server() {
    local port=$1
    local max_attempts=30
    local attempt=0

    print_warning "ポート$port でサーバーの起動を待っています..."

    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:$port/health >/dev/null 2>&1; then
            print_success "サーバーが起動しました！"
            return 0
        fi
        attempt=$((attempt + 1))
        sleep 1
        echo -n "."
    done

    echo ""
    print_error "サーバーが30秒以内に起動しませんでした"
    return 1
}

cleanup() {
    print_warning "\nクリーンアップ中..."
    if [ ! -z "$SERVER_PID" ] && kill -0 $SERVER_PID 2>/dev/null; then
        print_warning "サーバーを停止します (PID: $SERVER_PID)"
        kill -TERM $SERVER_PID 2>/dev/null
        wait $SERVER_PID 2>/dev/null
    fi
    kill_port 3000
    exit $EXIT_CODE
}

trap cleanup EXIT INT TERM

print_success "=== テストスイート開始 ==="

TEST_PATH=${1:-"./e2e/**/*"}
print_warning "テスト対象: $TEST_PATH"

kill_port 3000

print_warning "サーバーを起動しています..."
go run ./cmd/server/main.go &
SERVER_PID=$!

if [ -z "$SERVER_PID" ] || ! kill -0 $SERVER_PID 2>/dev/null; then
    print_error "サーバーの起動に失敗しました"
    EXIT_CODE=1
    exit 1
fi

print_success "サーバーを起動しました (PID: $SERVER_PID)"

# サーバーの準備完了を待つ
if ! wait_for_server 3000; then
    EXIT_CODE=1
    exit 1
fi

print_success "\n=== E2Eテスト実行中 ==="
runn run $TEST_PATH --verbose --debug

EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
    print_success "\n✅ すべてのテストが成功しました！"
else
    print_error "\n❌ いくつかのテストが失敗しました (終了コード: $EXIT_CODE)"
fi

