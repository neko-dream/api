package postgresql

import (
	"database/sql"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/neko-dream/api/internal/infrastructure/config"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

func Connect(config *config.Config, tp trace.TracerProvider) *sql.DB {
	db, err := otelsql.Open("pgx",
		config.DatabaseURL,
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithTracerProvider(tp),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			Ping:                 false,
			OmitConnResetSession: true,
			OmitRows:             true,
			OmitConnectorConnect: true,
			RecordError: func(err error) bool {
				return err != sql.ErrNoRows
			},
			RowsNext: false,
		}),
	)
	if err != nil {
		panic(err)
	}

	// DBStatsメトリクスを登録
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		panic(err)
	}

	// コネクションプールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	// 使っていない接続を閉じて、Aurora Serverless v2 がアイドル時に自動一時停止できるようにする
	db.SetConnMaxIdleTime(5 * time.Minute)

	return db
}
