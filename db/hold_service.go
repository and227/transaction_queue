package db

import (
	"context"
	"database/sql"
	"fmt"
)

type HoldService struct {
	dbConn *sql.DB
}

func NewHoldService(dbConn *sql.DB) *HoldService {
	return &HoldService{dbConn: dbConn}
}

func (s HoldService) Delete(txId int) error {
	query := "DELETE from hold where transaction_id = $1"
	_, err := s.dbConn.Exec(query, txId)
	return err
}

func (s HoldService) Reverse(ctx context.Context, txId int) error {
	var (
		holdAmount int
		balanceId  int
		txType     string
	)
	dbTx, err := s.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer dbTx.Rollback()

	query := `
		SELECT transaction_id, balance_id, amount
		FROM hold
		WHERE transaction_id = $1
		RETURNING balance_id, amount, tx_type
	`
	err = dbTx.QueryRow(query, txId).Scan(&balanceId, &holdAmount, &txType)

	if txType == "withdraw" {
		query = `
			UPDATE balance
			SET amount = amount + $1
			WHERE id = $2
		`
	} else if txType == "deposit" {
		query = `
			UPDATE balance
			SET amount = amount - $1
			WHERE id = $2
		`
	} else {
		return fmt.Errorf("Wrong transaction type")
	}

	_, err = dbTx.Exec(query, holdAmount, balanceId)

	query = "DELETE from hold where transaction_id = $1"
	_, err = dbTx.Exec(query, txId)

	if err = dbTx.Commit(); err != nil {
		return err
	}

	return err
}
