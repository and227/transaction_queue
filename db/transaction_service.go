package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TransactionService struct {
	dbConn *sql.DB
}

func NewTransactionService(dbConn *sql.DB) *TransactionService {
	return &TransactionService{dbConn: dbConn}
}

func (s TransactionService) Create(ctx context.Context, newTx TransactionCreateDTO) (Transaction, error) {
	var createdTx Transaction
	var balanceId int
	var balanceAmount int

	fmt.Printf("Handle tx: %v", newTx)

	query := `
		INSERT INTO transaction(user_id, amount, tx_type) 
		VALUES ($1,$2,$3)
		RETURNING id, user_id, amount, tx_type
	`
	err := s.dbConn.QueryRow(query, newTx.UserId, newTx.Amount, newTx.Type).
		Scan(&createdTx.Id, &createdTx.UserId, &createdTx.Amount, &createdTx.Type)
	if err != nil {
		return Transaction{}, err
	}
	fmt.Printf("tx created: %v\n", createdTx)

	dbTx, err := s.dbConn.BeginTx(ctx, nil)
	if err != nil {
		return Transaction{}, err
	}
	defer dbTx.Rollback()

	if createdTx.Type == "withdraw" {
		query = `
		UPDATE balance
		SET amount = amount - $1
		WHERE user_id = $2
		RETURNING id, amount
	`
	} else if createdTx.Type == "deposit" {
		query = `
		UPDATE balance
		SET amount = amount + $1
		WHERE user_id = $2
		RETURNING id, amount
	`
	} else {
		return Transaction{}, fmt.Errorf("Not valid transaction type")
	}

	err = dbTx.QueryRow(query, createdTx.Amount, createdTx.UserId).
		Scan(&balanceId, &balanceAmount)
	if err != nil {
		return Transaction{}, err
	}
	if balanceAmount < 0 {
		return Transaction{}, fmt.Errorf("Balance cannot be less than 0")
	}
	fmt.Println("balance ok")

	query = `
		INSERT INTO hold(transaction_id, balance_id, amount, tx_type)
		VALUES ($1,$2,$3,$4)
	`
	_, err = dbTx.Exec(query, createdTx.Id, balanceId, createdTx.Amount, createdTx.Type)
	if err != nil {
		return Transaction{}, err
	}
	fmt.Println("hold ok")

	if err = dbTx.Commit(); err != nil {
		return Transaction{}, err
	}
	fmt.Println("commit ok")

	return createdTx, nil
}

func (s TransactionService) GetAll() ([]Transaction, error) {
	var (
		txs []Transaction
		tx  Transaction
	)
	query := "SELECT id, user_id, amount, tx_type from transaction"
	rows, err := s.dbConn.Query(query)
	if err != nil {
		return []Transaction{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&tx.Id, &tx.UserId, &tx.Amount, &tx.Type)
		if err != nil {
			return []Transaction{}, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}
