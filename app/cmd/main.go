package main

// @title Warehouse API Documentation
// @description This is a sample API for a warehouse application
// @version 1
// @host localhost:8080
// @BasePath /api/v1
import (
	"context"
	"log"

	"github.com/DmitriiKumancev/lamoda-test/internal/app"
	"github.com/DmitriiKumancev/lamoda-test/internal/config"
	"github.com/DmitriiKumancev/lamoda-test/pkg/logging"
	
	_ "github.com/lib/pq"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logger := logging.GetLogger(ctx)

	logger.Info("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	ctx = logging.ContextWithLogger(ctx, logger)

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Running Application")
	a.Run(ctx)
}
