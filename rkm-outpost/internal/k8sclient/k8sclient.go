package k8sclient

import (
	"context"
	"log"
	"rkm-outpost/internal/config"
	"rkm-outpost/internal/influxdb"
	"rkm-outpost/internal/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	config    *config.K8sConfig
	clientset *kubernetes.Clientset
	influxDb  influxdb.InfluxSender
	logger    *logger.Logger
}

func NewK8sClient(config *config.K8sConfig, influxDb influxdb.InfluxSender, logger *logger.Logger) (*K8sClient, error) {
	var err error
	k8sClient := K8sClient{config: config, influxDb: influxDb, logger: logger}

	cfg, err := clientcmd.BuildConfigFromFlags("", k8sClient.config.ConfigPath)
	if err != nil {
		return nil, err
	}

	k8sClient.clientset, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &k8sClient, nil
}

func (k *K8sClient) GetNodeStatus() {
	nodes, err := k.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range nodes.Items {
		status := 0
		for _, c := range n.Status.Conditions {
			switch c.Type {
			case "Ready":
				if c.Status == v1.ConditionTrue {
					status = 1
				}
			default:
				if c.Status == v1.ConditionFalse {
					status = 1
				}
			}

			tags := map[string]string{
				"clustername": k.config.ClusterName,
				"node":        n.Name,
				"condition":   string(c.Type),
			}

			k.influxDb.AddNewPoint("rkm_node_conditiontype_health", tags, status)
		}

	}
}
