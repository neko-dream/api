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
func HandleError(ctx context.Context, err error) ErrorInfo {
	// スタックトレース情報を取得
	pc, file, line, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc).Name()

	errorInfo := ErrorInfo{
		Err:      err,
		File:     file,
		Line:     line,
		Function: function,
	}

	log.Printf("Error occurred: %v\nfile: %s\nline: %d\nfunction: %s\n",
		err, file, line, function)

	// スパンに情報を追加
	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err, trace.WithAttributes(
		attribute.String("error.file", file),
		attribute.Int("error.line", line),
		attribute.String("error.function", function),
	))

	return errorInfo
}
