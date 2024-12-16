package app

import (
	"database/sql"
	"fmt"

	"github.com/dexguitar/p2p_service/config"
	"github.com/dexguitar/p2p_service/internal/gateways/rest"
	"github.com/dexguitar/p2p_service/internal/rabbitmq"
	"github.com/dexguitar/p2p_service/internal/repo"
	"github.com/dexguitar/p2p_service/internal/usecase"
)

type App struct {
	server *rest.Server
	mq     *rabbitmq.RabbitMQ
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB))
	if err != nil {
		return nil, err
	}

	mq, err := rabbitmq.NewRabbitMQ(cfg)
	if err != nil {
		return nil, err
	}

	producer := rabbitmq.NewProducer(mq)
	repoTrans := repo.NewTransactionRepo(db)
	usecaseP2P := usecase.NewP2PTransferUsecase(repoTrans, producer)
	handler := rest.NewHandler(usecaseP2P)
	router := rest.NewRouter(handler)
	server := rest.NewServer(router, ":8080")

	mq.StartConsumer()

	rabbitmq.StartOutboxDispatcher(db, producer)

	return &App{server: server, mq: mq}, nil
}

func (a *App) Run() error {
	return a.server.Start()
}

func (a *App) Shutdown() {
	a.mq.Close()
}
