package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"webCrawler3/config"
	"webCrawler3/infra/crawler"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/log"
	"webCrawler3/infra/mgr/cluster"
	"webCrawler3/infra/mgr/router"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// This is the main function
	log.InitLogger()
	if err := config.LoadConfig(); err != nil {
		log.Logger.Fatal("Error loading config", zap.Error(err))
	}

	_, err := repoDb.Initialize(ctx, *config.Conf)
	if err != nil {
		log.Logger.Fatal("Error initializing database", zap.Error(err))
	}

	// Captura sinais de interrupção para cancelar o contexto
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	// Inicia Ouvinte as Rotas HTTP
	go router.InitRouters(ctx, config.Conf)

	// Inicia Ouvinte UDP
	go cluster.InitCluster(ctx, config.Conf)

	// Agendamento de tarefas
	//go tasks.Handle(ctx)

	// Inicia o processo de rastreamento
	go crawler.HandleQueue()

	// Aguarda o cancelamento do contexto
	<-ctx.Done()
	//defer dbConn.Close(ctx)
	fmt.Println("Crawler finalizado")
}
