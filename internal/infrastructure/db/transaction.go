package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"braces.dev/errtrace"
	model "github.com/neko-dream/server/internal/infrastructure/db/sqlc"
)

// DBManager DBコネクションおよびトランザクションを管理
type DBManager struct {
	db *sql.DB
}

// ExecTx トランザクションを実行
func (s *DBManager) ExecTx(ctx context.Context, fn func(context.Context) error) error {
	var tx *sql.Tx

	// トランザククションが終了するまで待つ
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Done()

	if tmpTx := getTransaction(ctx); tmpTx != nil {
		// トランザクションが既に開始されているのならそれを使う
		tx = tmpTx
	} else {
		// トランザクションが開始されていない場合は新しくトランザクションを開始する
		tmpTx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return errtrace.Wrap(err)
		}
		tx = tmpTx
	}

	// トランザクションをコンテキストにセットし処理を実行
	if err := fn(context.WithValue(ctx, key, tx)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return errtrace.Wrap(fmt.Errorf("tx err: %v, rb err: %v", err, rbErr))
		}
		return errtrace.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		return errtrace.Wrap(err)
	}

	return nil
}

// GetQueries トランザクションが開始されている場合はトラトランザクションを返す。そうでない場合はDBコネクションを返す。
func (s *DBManager) GetQueries(ctx context.Context) *model.Queries {
	// トランザクションが開始されている場合はトランザクションを返す
	if tx := getTransaction(ctx); tx != nil {
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
