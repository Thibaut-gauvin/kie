package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var clusterImagesCount = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "kie",
		Subsystem: "cluster",
		Name:      "images_total",
		Help:      "k8s cluster running images",
	},
	[]string{"name", "container", "pod", "namespace"},
)

func ObserveClusterImagesCount(name, container, pod, namespace string) {
	clusterImagesCount.With(
		prometheus.Labels{
			"name":      name,
			"container": container,
			"pod":       pod,
			"namespace": namespace,
		},
	).Inc()
}
