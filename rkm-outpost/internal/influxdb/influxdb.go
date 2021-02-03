package influxdb

import (
	"rkm-outpost/internal/config"
	"rkm-outpost/internal/logger"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxDbClient struct {
	config *config.InfluxConfig
	logger *logger.Logger
	points []*client.Point
}

func NewInfluxDbClient(config *config.InfluxConfig, logger *logger.Logger) *InfluxDbClient {
	return &InfluxDbClient{
		config: config,
		logger: logger,
	}
}

func (i *InfluxDbClient) AddNewPoint(measurement string, tags map[string]string, value interface{}) error {
	fields := map[string]interface{}{
		"value": value,
	}
	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		return err
	}
	i.points = append(i.points, pt)
	i.logger.Debug().Msgf("adding point to measurement=%s with tags=%+v and value=%v", measurement, tags, value)
	return nil
}

func (i *InfluxDbClient) Send() error {
	httpConfig := client.HTTPConfig{
		Addr: i.config.InfluxDbUrl,
	}

	if i.config.AuthEnabled {
		i.logger.Debug().Msg("influx auth enabled, using username password")
		httpConfig.Username = i.config.InfluxDbUser
		httpConfig.Password = i.config.InfluxDbPass
	}

	c, err := client.NewHTTPClient(httpConfig)
	if err != nil {
		return err
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{Database: i.config.InfluxDbName, Precision: "s"})
	if err != nil {
		return err
	}
	bp.AddPoints(i.points)

	i.logger.Info().Str("influxUrl", i.config.InfluxDbUrl).Int("totalPoints", len(i.points)).Msg("try sending data to influxdb")
	return c.Write(bp)
}

type InfluxSender interface {
	AddNewPoint(measurement string, tags map[string]string, value interface{}) error
	Send() error
}
