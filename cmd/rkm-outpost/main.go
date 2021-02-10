package main

import (
	"github.com/onzack/rkm/internal/config"
	"github.com/onzack/rkm/internal/influxdb"
	"github.com/onzack/rkm/internal/k8sclient"
	"github.com/onzack/rkm/internal/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("rkm-outpost starting")
	rkmOutpostConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error while creating config")
	}
	logger := logger.New(rkmOutpostConfig.Debug)

	influx := influxdb.NewInfluxDbClient(rkmOutpostConfig.InfluxConfig, logger)

	k8sClient, err := k8sclient.NewK8sClient(rkmOutpostConfig.K8sConfig, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error while creating k8s client")
	}

	log.Info().Msg("read metrics data")

	// k8s client fetch data
	k8sClient.GetNodeStatus()
	k8sClient.GetEndpointStatus()

	// send to influx
	if err := influx.Send(k8sClient.GetMetrics()); err != nil {
		logger.Error().Err(err).Msg("error while sending metrics to influx")
	}

	log.Info().Msg("data sent to influx")
	log.Info().Msg("rkm-outpost stopping")
}
