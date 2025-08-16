#!/bin/bash

# ログファイル（デバッグ用）
LOG_FILE="/tmp/claude-go-hook.log"

echo "=== Go品質チェックフック開始 ===" >> "$LOG_FILE"
echo "日時: $(date)" >> "$LOG_FILE"
echo "ツール: $CLAUDE_TOOL" >> "$LOG_FILE"
echo "パラメータ: $CLAUDE_PARAMS" >> "$LOG_FILE"

# ツールパラメータからファイルパスを抽出
FILE_PATH=$(echo "$CLAUDE_PARAM_file_path" | jq -r '.')
echo "ファイルパス: $FILE_PATH" >> "$LOG_FILE"

# チェック失敗を追跡
CHECKS_FAILED=0

# Goファイルかチェック
if [[ "$FILE_PATH" == *.go ]]; then
    echo "🔍 Go品質チェックを実行中: $FILE_PATH"
    echo "Goファイルを処理中: $FILE_PATH" >> "$LOG_FILE"

    # go fmtを実行
    echo "go fmtを実行中..." >> "$LOG_FILE"
    FORMATTED=$(go fmt "$FILE_PATH" 2>&1)
    if [ -n "$FORMATTED" ]; then
        echo "⚠️  ファイルがフォーマットされていませんでした: $FORMATTED" | tee -a "$LOG_FILE"
        echo "フォーマットの問題を検出しました。go fmtによる自動修正を適用しました。" >> "$LOG_FILE"
        CHECKS_FAILED=1
    else
        echo "フォーマットチェック: OK" >> "$LOG_FILE"
    fi

    # golangci-lintを実行
    if command -v golangci-lint &> /dev/null; then
        echo "golangci-lintを実行中..." >> "$LOG_FILE"
        if ! golangci-lint run "$FILE_PATH" 2>&1 | tee -a "$LOG_FILE"; then
            echo "❌ golangci-lintで問題が見つかりました" | tee -a "$LOG_FILE"
            echo "静的解析で問題を検出しました。上記のエラーを修正してください。" >> "$LOG_FILE"
            CHECKS_FAILED=1
        else
            echo "静的解析チェック: OK" >> "$LOG_FILE"
        fi
    else
        echo "⚠️  golangci-lintが見つかりません。リンターチェックをスキップします" >> "$LOG_FILE"
        echo "golangci-lintのインストールを推奨します: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" >> "$LOG_FILE"
    fi

    # テストファイルの場合はテストを実行
    if [[ "$FILE_PATH" == *_test.go ]]; then
        DIR=$(dirname "$FILE_PATH")
        echo "テストを実行中: $DIR ..." >> "$LOG_FILE"
        echo "対象ディレクトリ: $DIR" >> "$LOG_FILE"
        if ! go test "$DIR" -v 2>&1 | tee -a "$LOG_FILE"; then
            echo "❌ テストが失敗しました" | tee -a "$LOG_FILE"
            echo "テストの失敗を検出しました。テスト結果を確認してください。" >> "$LOG_FILE"
            CHECKS_FAILED=1
        else
            echo "テスト実行: OK - 全てのテストが成功しました" >> "$LOG_FILE"
        fi
    else
        echo "テストファイルではないため、テスト実行をスキップ" >> "$LOG_FILE"
    fi

    if [ $CHECKS_FAILED -eq 0 ]; then
        echo "✅ 全てのGo品質チェックに合格しました"
        echo "結果: 全てのチェックに合格" >> "$LOG_FILE"
    else
        echo "❌ Go品質チェックで問題が見つかりました"
        echo "結果: 品質チェックで問題を検出" >> "$LOG_FILE"
    fi
else
    echo "Goファイルではないため、チェックをスキップします" >> "$LOG_FILE"
fi

echo "=== フック完了 ===" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# 適切な終了コードで終了
if [ $CHECKS_FAILED -eq 1 ]; then
    echo "チェックエラーのため終了コード2で終了します" >> "$LOG_FILE"
    exit 2
fi

echo "正常終了" >> "$LOG_FILE"
exit 0
