package utils

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"go.opentelemetry.io/otel"
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
	ctx, span := otel.Tracer("utils").Start(ctx, "HandleError")
	defer span.End()

	HandleErrorWithCaller(ctx, err, message, 2)
}

func HandleErrorWithCaller(ctx context.Context, err error, message string, caller int) {
	_, span := otel.Tracer("utils").Start(ctx, "HandleErrorWithCaller")
	defer span.End()

	pc, file, line, _ := runtime.Caller(caller)
	function := runtime.FuncForPC(pc).Name()

	attrs := []attribute.KeyValue{
		attribute.String("error.file", file),
		attribute.Int("error.line", line),
		attribute.String("error.function", function),
		attribute.String("error.message", message),
	}

	// エラーメッセージを出力
	fmt.Println("+----------------------------------------")
	fmt.Printf("time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		fmt.Printf("mesg: %s: %s\n", message, err.Error())
		attrs = append(attrs, attribute.String("error", err.Error()))
	} else {
		fmt.Printf("mesg: %s\n", message)
	}
	fmt.Printf("file: %s:%d\n", file, line)
	fmt.Printf("func: %s\n", function)
	fmt.Println("+----------------------------------------")

	span.SetStatus(codes.Error, message)
	span.RecordError(err, trace.WithAttributes(
		attrs...,
	))
	span.SetAttributes(attrs...)
}
