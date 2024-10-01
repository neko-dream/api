package db

import (
	"context"
	"database/sql"

	"braces.dev/errtrace"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
	"github.com/neko-dream/server/pkg/utils"
)

// DBManager DBコネクションおよびトランザクションを管理
type DBManager struct {
	db *sql.DB
}

func (s *DBManager) ExecTx(ctx context.Context, fn func(context.Context) error) error {
	panicked := true
	var tx *sql.Tx

	if tmpTx := getTransaction(ctx); tmpTx != nil {
		// トランザクションが既に開始されているのならそれを使う
		tx = tmpTx
	} else {
		// トランザクションが開始されていない場合は新しくトランザクションを開始する
		tmpTx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			utils.HandleError(ctx, err, "トランザクションの開始に失敗")
			return errtrace.Wrap(err)
		}
		tx = tmpTx
	}

	defer func() {
		if panicked {
			if rbErr := tx.Rollback(); rbErr != nil {
				utils.HandleError(ctx, rbErr, "トランザクションのロールバックに失敗")
			}
		}
	}()

	if err := fn(context.WithValue(ctx, key, tx)); err != nil {
		utils.HandleError(ctx, err, "トランザクション内の処理で失敗")
		return err
	}

	if err := tx.Commit(); err != nil {
		utils.HandleError(ctx, err, "トランザクションのコミットに失敗")
		return errtrace.Wrap(err)
	}

	panicked = false

	return nil
}

// GetQueries トランザクションが開始されている場合はトラトランザクションを返す。そうでない場合はDBコネクションを返す。
func (s *DBManager) GetQueries(ctx context.Context) *model.Queries {
	// トランザクションが開始されている場合はトランザクションを返す
	if tx, ok := ctx.Value(key).(*sql.Tx); ok {
		return model.New(tx)
	}
	// トランザクションが開始されていない場合はDBコネクションを返す
	return model.New(s.db)
}

type transactionCtxKey string

const (
	key transactionCtxKey = "TransactionKey"
)

func getTransaction(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(key).(*sql.Tx)
	if ok && tx != nil {
		return tx
	}
	return nil
}

func NewDBManager(db *sql.DB) *DBManager {
	return &DBManager{
		db: db,
	}
}
