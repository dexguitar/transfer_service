package usecase

import (
	"errors"

	"github.com/dexguitar/p2p_service/internal/domains"
	"github.com/dexguitar/p2p_service/internal/rabbitmq"
	"github.com/dexguitar/p2p_service/internal/repo"
)

type P2PTransferUsecase interface {
	Transfer(sender, receiver string, amount float64) (int, error)
}

type p2pTransferUsecase struct {
	transactionRepo repo.TransactionRepo
	producer        rabbitmq.Producer
}

func NewP2PTransferUsecase(tr repo.TransactionRepo, p rabbitmq.Producer) P2PTransferUsecase {
	return &p2pTransferUsecase{transactionRepo: tr, producer: p}
}

func (u *p2pTransferUsecase) Transfer(sender, receiver string, amount float64) (int, error) {
	if sender == "" || receiver == "" || amount <= 0 {
		return 0, errors.New("invalid input")
	}

	id, err := u.transactionRepo.CreateTransaction(domains.Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	})
	if err != nil {
		return 0, err
	}

	err = u.producer.PublishTransactionID(id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
