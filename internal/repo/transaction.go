package repo

import (
	"database/sql"
	"encoding/json"

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
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var transactionID int
	err = r.db.QueryRow(
		"INSERT INTO transactions (sender, receiver, amount) VALUES ($1, $2, $3) RETURNING id",
		t.Sender, t.Receiver, t.Amount,
	).Scan(&transactionID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	eventPayload := map[string]interface{}{
		"transaction_id": transactionID,
		"sender":         t.Sender,
		"receiver":       t.Receiver,
		"amount":         t.Amount,
	}
	payloadBytes, _ := json.Marshal(eventPayload)
	_, err = tx.Exec(
		"INSERT INTO outbox (aggregate_type, aggregate_id, event_type, payload) VALUES ($1, $2, $3, $4)",
		"transaction", transactionID, "transaction_created", payloadBytes,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return transactionID, nil
}
