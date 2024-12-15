package repo

import (
	"database/sql"

	"github.com/dexguitar/p2p_service/internal/domains"
)

type TransactionRepo interface {
	CreateTransaction(t domains.Transaction) (int, error)
}

type transactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) TransactionRepo {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) CreateTransaction(t domains.Transaction) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO transactions (sender, receiver, amount) VALUES ($1, $2, $3) RETURNING id",
		t.Sender, t.Receiver, t.Amount,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
