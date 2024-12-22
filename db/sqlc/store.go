package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (store *Store) execTx(ctx context.Context, fn func(query *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return err
	}

	q := store.WithTx(tx)

	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("fn error %v;tx rollback failed: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, transferParams TransferTxParams) (error, TransferTxResult) {
	var transferResult TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		transferResult.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: transferParams.FromAccountID,
			ToAccountID:   transferParams.ToAccountID,
			Amount:        transferParams.Amount,
		})
		if err != nil {
			return err
		}
		transferResult.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transferParams.FromAccountID,
			Amount:    -transferParams.Amount,
		})
		if err != nil {
			return err
		}
		transferResult.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: transferParams.ToAccountID,
			Amount:    transferParams.Amount,
		})
		if err != nil {
			return err
		}

		//TODO: update account balance

		transferResult.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: transferParams.FromAccountID,
			Amount: -transferParams.Amount,
		})
		if err != nil {
			return err
		}

		transferResult.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID: transferParams.ToAccountID,
			Amount: transferParams.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return err, transferResult

}
