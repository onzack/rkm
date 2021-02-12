package influxdb

import (
	"time"

	"github.com/onzack/rkm/internal/config"
	"github.com/onzack/rkm/internal/logger"
	"github.com/onzack/rkm/internal/metrics"

	client "github.com/influxdata/influxdb1-client/v2"
)

// InfluxSender is an interface that defindes an send methond
type InfluxSender interface {
	Send(metricsCollector *metrics.Collector) error
}

// InfluxDbClient is a struct that contains the following fields: config, logger, points
type InfluxDbClient struct {
	config *config.InfluxConfig
	logger *logger.Logger
	points []*client.Point
}

// NewInfluxDbClient is a function that takes some config arguments and returns an new InfluxDB Client
func NewInfluxDbClient(config *config.InfluxConfig, logger *logger.Logger) *InfluxDbClient {
	return &InfluxDbClient{
		config: config,
		logger: logger,
	}
}

func (i *InfluxDbClient) Send(metricsCollector *metrics.Collector) error {
	httpConfig := client.HTTPConfig{
		Addr: i.config.InfluxDbUrl,
	}

	if i.config.AuthEnabled {
		i.logger.Debug().Msg("influx auth enabled, using username password")
		httpConfig.Username = i.config.InfluxDbUser
		httpConfig.Password = i.config.InfluxDbPW
	}

	c, err := client.NewHTTPClient(httpConfig)
	if err != nil {
		return err
	}
	defer c.Close()

	bp, err := i.addMetricPoints(metricsCollector)
	if err != nil {
		return err
	}

	i.logger.Info().Str("influxUrl", i.config.InfluxDbUrl).Int("totalPoints", len(bp.Points())).Msg("try sending data to influxdb")
	return c.Write(bp)
}

func (i *InfluxDbClient) addMetricPoints(m *metrics.Collector) (client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{Database: i.config.InfluxDbName, Precision: "s"})
	if err != nil {
		return nil, err
	}

	for _, m := range m.MetricPoints {
		p, err := client.NewPoint(m.Measurement, m.Tags, m.Value, time.Now())
		if err != nil {
			i.logger.Error().Err(err).Msg("unable to add point to batchpoint")
			continue
		}
		i.logger.Debug().Msgf("%s", m)
		bp.AddPoint(p)
	}
	bp.AddPoints(i.points)
	return bp, nil
}
