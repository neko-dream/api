package db

import (
	"context"
	"database/sql"
	"fmt"

	"braces.dev/errtrace"
	model "github.com/neko-dream/server/internal/infrastructure/persistence/sqlc/generated"
	"github.com/neko-dream/server/pkg/utils"
	"go.opentelemetry.io/otel"
)

type txKey struct{}

func getTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// DBManager DBコネクションおよびトランザクションを管理
type DBManager struct {
	db *sql.DB
}

func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{
		db: db,
	}
}

func (s *DBManager) TestTx(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, span := otel.Tracer("db").Start(ctx, "DBManager.TestTx")
	defer span.End()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		utils.HandleError(ctx, err, "トランザクションの開始に失敗")
		return errtrace.Wrap(err)
	}

	if err := fn(context.WithValue(ctx, txKey{}, tx)); err != nil {
		utils.HandleError(ctx, err, "トランザクション内の処理で失敗")
		return err
	}

	if err := tx.Rollback(); err != nil {
		utils.HandleError(ctx, err, "トランザクションのロールバックに失敗")
		return errtrace.Wrap(err)
	}

	return nil
}

func (s *DBManager) ExecTx(ctx context.Context, fn func(context.Context) error) error {
	ctx, span := otel.Tracer("db").Start(ctx, "DBManager.ExecTx")
	defer span.End()

	var tx *sql.Tx

	if tmpTx := getTx(ctx); tmpTx != nil {
		// トランザクションが既に開始されているのならそれを使う
		return fn(ctx)
	}

	// トランザクションが開始されていない場合は新しくトランザクションを開始する
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		utils.HandleError(ctx, err, "トランザクションの開始に失敗")
		return errtrace.Wrap(err)
	}

	defer func() {
		if r := recover(); r != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				utils.HandleError(ctx, rbErr, "トランザクションのロールバックに失敗")
			}
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(txCtx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			utils.HandleError(ctx, rbErr, fmt.Sprintf("ロールバックに失敗: %v (元エラー: %v)", rbErr, err))
			return fmt.Errorf("ロールバックに失敗: %v (元エラー: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(ctx, err, "トランザクションのコミットに失敗")
		return errtrace.Wrap(err)
	}

	return nil
}

// GetQueries トランザクションが開始されている場合はトラトランザクションを返す。そうでない場合はDBコネクションを返す。
func (s *DBManager) GetQueries(ctx context.Context) *model.Queries {
	ctx, span := otel.Tracer("db").Start(ctx, "DBManager.GetQueries")
	defer span.End()

	// トランザクションが開始されている場合はトランザクションを返す
	if tx := getTx(ctx); tx != nil {
		return model.New(tx)
	}

	// トランザクションが開始されていない場合はDBコネクションを返す
	return model.New(s.db)
}
