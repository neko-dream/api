package utils

import (
	"context"
	"log"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ErrorInfo はエラーに関する追加情報を保持する構造体です
type ErrorInfo struct {
	Err      error
	File     string
	Line     int
	Function string
}

// HandleError はエラーを処理し、追加情報を記録する関数です
func HandleError(ctx context.Context, err error, message string) {
	HandleErrorWithCaller(ctx, err, message, 2)
}

func HandleErrorWithCaller(ctx context.Context, err error, message string, caller int) {
	pc, file, line, _ := runtime.Caller(caller)
	function := runtime.FuncForPC(pc).Name()

	// エラーメッセージを出力
	log.Printf("%s: %s\n", message, err.Error())
	log.Printf("file: %s:%d, function: %s\n", file, line, function)

	// スパンに情報を追加
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err, trace.WithAttributes(
		attribute.String("error.file", file),
		attribute.Int("error.line", line),
		attribute.String("error.function", function),
		attribute.String("error.message", message),
	))
}
