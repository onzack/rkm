package cmd

import (
	"rkm-outpost/internal/influxdb"
	"rkm-outpost/internal/k8sclient"
	"rkm-outpost/internal/logger"
)

type App struct {
	influxClient influxdb.InfluxSender
	k8sClient    *k8sclient.K8sClient
	logger       *logger.Logger
}

func NewApp(influxDbClient influxdb.InfluxSender, k8sClient *k8sclient.K8sClient, logger *logger.Logger) *App {
	return &App{influxClient: influxDbClient, k8sClient: k8sClient, logger: logger}
}

func (a *App) Run() {
	a.logger.Info().Msg("read metrics data")

	// k8s client fetch data
	a.k8sClient.GetNodeStatus()

	// send to influx
	if err := a.influxClient.Send(a.k8sClient.GetMetrics()); err != nil {
		a.logger.Error().Err(err).Msg("error while sending metrics to influx")
	}
	a.logger.Info().Msg("data sent to influx")

}
