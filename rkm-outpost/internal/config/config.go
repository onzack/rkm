package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

var ErrAuthDetailsMissing = errors.New("auth credentials missing")

type Config struct {
	Debug        bool `envconfig:"DEBUG" default:"true"`
	InfluxConfig *InfluxConfig
	K8sConfig    *K8sConfig
}

type K8sConfig struct {
	ClusterName string `envconfig:"CLUSTER_NAME" default:"k8sdev"`
	ConfigPath  string `envconfig:"CONFIG_PATH" default:"/Users/andiexer/.kube/config"`
}

type InfluxConfig struct {
	InfluxDbUrl  string `envconfig:"INFLUXDB_URL" default:"http://localhost:8087"`
	InfluxDbName string `envconfig:"INFLUXDB_NAME" default:"rkm-outpost"`
	InfluxDbUser string `envconfig:"INFLUXDB_USER" default:""`
	InfluxDbPass string `envconfig:"INFLUXDB_PASS" default:""`
	AuthEnabled  bool   `envconfig:"AUTH_ENABLED" default:"false"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.InfluxConfig.AuthEnabled {
		if cfg.InfluxConfig.InfluxDbUser == "" || cfg.InfluxConfig.InfluxDbPass == "" {
			return nil, ErrAuthDetailsMissing
		}
	}

	return &cfg, nil
}
