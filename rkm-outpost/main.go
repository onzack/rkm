package main

import (
	"rkm-outpost/cmd"
	"rkm-outpost/internal/config"
	"rkm-outpost/internal/influxdb"
	"rkm-outpost/internal/k8sclient"
	"rkm-outpost/internal/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("rkm-outpost starting")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error while creating config")
	}
	logger := logger.New(cfg.Debug)

	influx := influxdb.NewInfluxDbClient(cfg.InfluxConfig, logger)

	k8sClient, err := k8sclient.NewK8sClient(cfg.K8sConfig, influx, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while creating k8s client")
	}

	app := cmd.NewApp(influx, k8sClient, logger)
	app.Run()
	log.Info().Msg("rkm-outpost stopping")
}
