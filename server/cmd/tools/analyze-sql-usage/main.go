package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	// SQLクエリ関数名を読み込む
	sqlFunctions := make(map[string]bool)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "-- name: ") {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				funcName := parts[2]
				sqlFunctions[funcName] = true
			}
		}
	}

	// 使用状況を格納するマップ
	usages := make(map[string]int)

	// Goファイルを再帰的に探索
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 除外するディレクトリ
		if strings.Contains(path, "/vendor/") ||
			strings.Contains(path, "/testdata/") ||
			strings.Contains(path, "/oas/") ||
			strings.Contains(path, "/generated/") {
			return nil
		}

		// .goファイルのみを処理
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// ファイルを解析
		findUsagesInFile(path, sqlFunctions, usages)
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking files: %v\n", err)
		os.Exit(1)
	}

	// 結果を関数名でソートして出力
	var funcNames []string
	for name := range sqlFunctions {
		funcNames = append(funcNames, name)
	}
	sort.Strings(funcNames)

	// 未使用の関数のみを出力
	for _, funcName := range funcNames {
		count := usages[funcName]
		if count == 0 {
			fmt.Println(funcName)
		}
	}
}

func findUsagesInFile(filename string, sqlFunctions map[string]bool, usages map[string]int) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	// ファイル内容を文字列として扱い、関数名を検索
	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	for funcName := range sqlFunctions {
		for _, line := range lines {
			// 関数名が単語境界で出現するかチェック
			if containsWord(line, funcName) {
				usages[funcName]++
			}
		}
	}
}

func containsWord(line, word string) bool {
	// 単語境界をチェックするための簡易的な実装
	index := strings.Index(line, word)
	if index == -1 {
		return false
	}

	// 前後の文字が単語の一部でないことを確認
	before := index == 0 || !isWordChar(line[index-1])
	after := index+len(word) >= len(line) || !isWordChar(line[index+len(word)])

	return before && after
}

func isWordChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_'
}
