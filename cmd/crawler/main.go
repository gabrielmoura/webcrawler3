package main

import (
	"context"
	"go.uber.org/zap"
	"webCrawler3/config"
	"webCrawler3/infra/crawler"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/log"
	"webCrawler3/infra/mgr/cluster"
	"webCrawler3/infra/mgr/router"
)

func main() {
	ctx := context.Background()
	// This is the main function
	log.InitLogger()
	if err := config.LoadConfig(); err != nil {
		log.Logger.Fatal("Error loading config", zap.Error(err))
	}

	_, err := repoDb.Initialize(ctx, *config.Conf)
	if err != nil {
		log.Logger.Fatal("Error initializing database", zap.Error(err))
	}

	// Inicia Ouvinte as Rotas HTTP
	go router.InitRouters(ctx, config.Conf)

	// Inicia Ouvinte UDP
	go cluster.InitCluster(ctx, config.Conf)

	// Agendamento de tarefas
	//go tasks.Handle(ctx)

	// Inicia o processo de rastreamento
	crawler.HandleQueue()

	// Aguarda o cancelamento do contexto
	<-ctx.Done()
	//defer dbConn.Close(ctx)
	println("Crawler finalizado")
}
