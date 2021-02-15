package k8sclient

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/onzack/rkm/internal/config"
	"github.com/onzack/rkm/internal/logger"
	"github.com/onzack/rkm/internal/metrics"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

const (
	rkmNodeConditionTypeHealth     string = "rkm_node_conditiontype_health"
	rkmAPIServerEndpointsTotal     string = "rkm_apiserver_endpoints_total"
	rkmComponentHealth             string = "rkm_component_health"
	rkmKubeAPIServerHealth         string = "rkm_kubeapiserver_health"
	rkmOverallHealth               string = "rkm_overall_health"
	rkmOutpostDurationMilliseconds string = "rkm_outpost_duration_milliseconds"
)

// K8sClient is a struct containing config, clientset, logger and metricsController
type K8sClient struct {
	config           *config.K8sConfig
	clientset        *kubernetes.Clientset
	logger           *logger.Logger
	metricsCollector *metrics.Collector
}

// NewK8sClient is a function which takes a config and a logger and returns a K8sClient
func NewK8sClient(config *config.K8sConfig, logger *logger.Logger) (*K8sClient, error) {
	var err error
	k8sClient := K8sClient{config: config, logger: logger}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		cfg, err = clientcmd.BuildConfigFromFlags("", defaultKubeconfig())
		if err != nil {
			return nil, err
		}
	}

	k8sClient.clientset, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	k8sClient.metricsCollector = &metrics.Collector{ClusterName: config.ClusterName}

	return &k8sClient, nil
}

// GetMetrics is am method for the K8sClient struct that returns the K8sClient metricsCollector
func (k *K8sClient) GetMetrics() *metrics.Collector {
	return k.metricsCollector
}

// GetNodeStatus is a method for the K8sClient struct that fetches the node status and calls
// the AddNodesMetricsEntry method to add the metrics to a MetricsPoint
func (k *K8sClient) GetNodeStatus(overallHealth *int, kubeAPIHealth *int) {
	nodes, err := k.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
		*kubeAPIHealth = 0
		*overallHealth = 0
	}
	for _, n := range nodes.Items {
		status := 1
		for _, c := range n.Status.Conditions {
			switch c.Type {
			case "Ready":
				if c.Status == v1.ConditionFalse || c.Status == v1.ConditionUnknown {
					status = 0
					*overallHealth = 0
				}
			default:
				if c.Status == v1.ConditionTrue || c.Status == v1.ConditionUnknown {
					status = 0
					*overallHealth = 0
				}
			}

			tags := map[string]string{
				"condition": string(c.Type),
			}

			k.metricsCollector.AddNodesMetricsEntry(rkmNodeConditionTypeHealth, n.Name, tags, status)
		}
	}
}

// GetEndpointStatus is a method for the K8sClient struct that fetches the endpoints status and calls
// the AddSimpleMetricsEntry method to add the metrics to a MetricsPoint
func (k *K8sClient) GetEndpointStatus(kubeAPIHealth *int) {
	endpoints, err := k.clientset.CoreV1().Endpoints("default").Get(context.Background(), "kubernetes", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
		*kubeAPIHealth = 0
	}
	endpointsCount := len(endpoints.Subsets[0].Addresses)
	k.metricsCollector.AddSimpleMetricsEntry(rkmAPIServerEndpointsTotal, endpointsCount)
}

// GetComponentStatus is a method for the K8sClient struct that fetches the components status and calls
// the AddComponentsMetricsEntry method to add the metrics to a MetricsPoint
func (k *K8sClient) GetComponentStatus(overallHealth *int, kubeAPIHealth *int) {
	components, err := k.clientset.CoreV1().ComponentStatuses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
		*kubeAPIHealth = 0
		*overallHealth = 0
	}
	for _, component := range components.Items {
		status := 1
		tags := map[string]string{}
		if component.Conditions[0].Status == "False" {
			status = 0
			*overallHealth = 0
			tags["component"] = string(component.Name)

		} else {
			tags["component"] = string(component.Name)
		}
		k.metricsCollector.AddComponentsMetricsEntry(rkmComponentHealth, tags, status)
	}
}

// SetKubeAPIAndOverallHealth is a method for the K8sClient struct that calls
// the AddSimpleMetricsEntry method to add the metrics to a MetricsPoint
func (k *K8sClient) SetKubeAPIAndOverallHealth(overallHealth *int, kubeAPIHealth *int) {
	k.metricsCollector.AddSimpleMetricsEntry(rkmKubeAPIServerHealth, *kubeAPIHealth)
	k.metricsCollector.AddSimpleMetricsEntry(rkmOverallHealth, *overallHealth)
}

// StopTimer is a method for the K8sClient struct that calculates the duration of rkm-oupost and calls the
// the AddSimpleMetricsEntry method to add the metrics to a MetricsPoint
func (k *K8sClient) StopTimer(startTime time.Time) {
	duration := time.Since(startTime)
	log.Printf("rkm-outpost took %s to finsh", duration)
	k.metricsCollector.AddSimpleMetricsEntry(rkmOutpostDurationMilliseconds, duration.Milliseconds())
}

// defaultKubeconfig is a function that reads the OS environments and returns the path to the kubeconfig
func defaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		klog.Warningf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
