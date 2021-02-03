package k8sclient

import (
	"context"
	"log"
	"rkm-outpost/internal/config"
	"rkm-outpost/internal/logger"
	"rkm-outpost/internal/metrics"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	rkmNodeConditionTypeHealth string = "rkm_node_conditiontype_health"
)

type K8sClient struct {
	config           *config.K8sConfig
	clientset        *kubernetes.Clientset
	logger           *logger.Logger
	metricsCollector *metrics.MetricsCollector
}

func NewK8sClient(config *config.K8sConfig, logger *logger.Logger) (*K8sClient, error) {
	var err error
	k8sClient := K8sClient{config: config, logger: logger}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", k8sClient.config.ConfigPath)
		if err != nil {
			return nil, err
		}
	}

	k8sClient.clientset, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	k8sClient.metricsCollector = &metrics.MetricsCollector{ClusterName: config.ClusterName}

	return &k8sClient, nil
}

func (k *K8sClient) GetMetrics() *metrics.MetricsCollector {
	return k.metricsCollector
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
				"condition": string(c.Type),
			}

			k.metricsCollector.AddMetricsEntry(rkmNodeConditionTypeHealth, n.Name, tags, status)
		}
	}
}
