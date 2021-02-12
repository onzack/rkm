package main

import (
	"time"

	"github.com/onzack/rkm/internal/config"
	"github.com/onzack/rkm/internal/influxdb"
	"github.com/onzack/rkm/internal/k8sclient"
	"github.com/onzack/rkm/internal/logger"

	"github.com/rs/zerolog/log"
)

func main() {
	StartTime := time.Now()
	OverallHealth := 0
	KubeAPIHealth := 0
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

	// k8s client collect metrics
	k8sClient.GetNodeStatus(&OverallHealth, &KubeAPIHealth)
	k8sClient.GetEndpointStatus(&KubeAPIHealth)
	k8sClient.GetComponentStatus(&OverallHealth, &KubeAPIHealth)
	k8sClient.SetKubeAPIAndOverallHealth(&OverallHealth, &KubeAPIHealth)
	k8sClient.StopTimer(StartTime)

	// send to influx
	if err := influx.Send(k8sClient.GetMetrics()); err != nil {
		logger.Error().Err(err).Msg("error while sending metrics to influx")
	}

	log.Info().Msg("data sent to influx")
	log.Info().Msg("rkm-outpost stopping")
}
